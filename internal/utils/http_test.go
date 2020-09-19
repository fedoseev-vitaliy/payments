package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClient_Get(t *testing.T) {
	type message struct {
		Message string `json:"message"`
	}
	type fields struct {
		timeout     time.Duration
		handlerFunc http.HandlerFunc
	}
	type args struct {
		ctx        context.Context
		ctxTimeout time.Duration
		succResp   interface{}
		errResp    interface{}
	}

	tests := []struct {
		name               string
		fields             fields
		args               args
		wantSuccResp       interface{}
		wantErrResp        interface{}
		wantErr            bool
		withContextTimeout bool
		ctx                context.Context
		succResp           interface{}
		errResp            interface{}
	}{
		{
			name: "nil parser for successful response",
			fields: fields{
				handlerFunc: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			},
			args: args{
				ctx:      context.Background(),
				succResp: nil,
				errResp:  nil,
			},
			wantErr: false,
		},
		{
			name: "nil parser for error response",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) },
			},
			args: args{
				ctx:      context.Background(),
				succResp: nil,
				errResp:  nil,
			},
			wantErr: false,
		},
		{
			name: "not nil parser for successful response",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) { _ = json.NewEncoder(w).Encode(message{Message: "OK"}) },
			},
			args: args{
				ctx:      context.Background(),
				succResp: &message{},
				errResp:  nil,
			},
			wantSuccResp: &message{Message: "OK"},
			wantErrResp:  nil,
			wantErr:      false,
		},
		{
			name: "not nil parser for error response",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(message{Message: "ERROR"})
				},
			},
			args: args{
				ctx:      context.Background(),
				succResp: nil,
				errResp:  &message{},
			},
			wantSuccResp: nil,
			wantErrResp:  &message{Message: "ERROR"},
			wantErr:      false,
		},
		{
			name: "not nil parsers for successful response",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) { _ = json.NewEncoder(w).Encode(message{Message: "OK"}) },
			},
			args: args{
				ctx:      context.Background(),
				succResp: &message{},
				errResp:  &message{},
			},
			wantSuccResp: &message{Message: "OK"},
			wantErrResp:  &message{},
			wantErr:      false,
		},
		{
			name: "not nil parsers and for error response",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusInternalServerError)
					_ = json.NewEncoder(w).Encode(message{Message: "ERROR"})
				},
			},
			args: args{
				ctx:      context.Background(),
				succResp: &message{},
				errResp:  &message{},
			},
			wantSuccResp: &message{},
			wantErrResp:  &message{Message: "ERROR"},
			wantErr:      false,
		},
		{
			name: "marshal error for successful response",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {},
			},
			args: args{
				ctx:      context.Background(),
				succResp: &message{},
				errResp:  &message{},
			},
			wantSuccResp: &message{},
			wantErrResp:  &message{},
			wantErr:      true,
		},
		{
			name: "not nil parsers and for error response",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusInternalServerError) },
			},
			args: args{
				ctx:      context.Background(),
				succResp: &message{},
				errResp:  &message{},
			},
			wantSuccResp: &message{},
			wantErrResp:  &message{},
			wantErr:      true,
		},
		{
			name: "timeout error",
			fields: fields{
				timeout:     5 * time.Millisecond,
				handlerFunc: func(w http.ResponseWriter, r *http.Request) { time.Sleep(20 * time.Millisecond) },
			},
			args: args{
				ctx:      context.Background(),
				succResp: &message{},
				errResp:  &message{},
			},
			wantSuccResp: &message{},
			wantErrResp:  &message{},
			wantErr:      true,
		},
		{
			name: "context timeout error",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) { time.Sleep(20 * time.Millisecond) },
			},
			args: args{
				ctx:        context.Background(),
				ctxTimeout: 5 * time.Millisecond,
				succResp:   &message{},
				errResp:    &message{},
			},
			wantSuccResp:       &message{},
			wantErrResp:        &message{},
			wantErr:            true,
			withContextTimeout: true,
		},
		{
			name: "nil context error",
			fields: fields{
				handlerFunc: func(w http.ResponseWriter, r *http.Request) {},
			},
			args: args{
				ctx:      nil,
				succResp: &message{},
				errResp:  &message{},
			},
			wantSuccResp:       &message{},
			wantErrResp:        &message{},
			wantErr:            true,
			withContextTimeout: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(tt.fields.handlerFunc)
			defer ts.Close()

			client := NewClient(tt.fields.timeout)

			ctx := tt.args.ctx

			if tt.withContextTimeout {
				c, cancel := context.WithTimeout(tt.args.ctx, tt.args.ctxTimeout)
				ctx = c
				defer cancel()
			}

			u, err := url.Parse(ts.URL)
			require.NoError(t, err)

			_, err = client.Get(ctx, u, tt.args.succResp, tt.args.errResp)

			if tt.wantErr {
				fmt.Println(err)
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			require.Equal(t, tt.wantSuccResp, tt.args.succResp)
			require.Equal(t, tt.wantErrResp, tt.args.errResp)
		})
	}

	t.Run("invalid url", func(t *testing.T) {
		client := NewClient(time.Minute)

		sc, err := client.Get(context.Background(), nil, nil, nil)
		require.Error(t, err)
		require.Equal(t, -1, sc)
	})
}

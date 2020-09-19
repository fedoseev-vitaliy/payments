package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/fedoseev-vitaliy/payments/internal/controller"
	"github.com/fedoseev-vitaliy/payments/internal/providers"
)

const (
	appleAppURL   = "http://apple.store.com/myApp"
	androidAppURL = "http://google.store.com/myApp"
)

type Controller interface {
	GetPaymentsURL(ctx context.Context, productID string) (*controller.PaymentsURLs, error)
}

type Handler struct {
	l *logrus.Logger
	c Controller
}

type Response struct {
	GooglePayURL string `json:"g_url"`
	ApplePayURL  string `json:"a_url"`
}

type AppURLResponse struct {
	AppleAppURL  string `json:"apple_url"`
	GoogleAppURL string `json:"google_url"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewHandler(l *logrus.Logger, c Controller) *Handler {
	return &Handler{l: l, c: c}
}

func (h *Handler) GetPaymentsURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusForbidden)
		if err := json.NewEncoder(w).Encode(&ErrorResponse{
			Error: "Only GET method supported",
		}); err != nil {
			h.l.Error(err.Error())
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	pid := r.URL.Query().Get("productID")
	// this stuff for testing purpose
	switch pid {
	case "":
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(&ErrorResponse{
			Error: "orderID query params is missing",
		}); err != nil {
			h.l.Error(err.Error())
		}
		return
	case "panic":
		panic("panic test")
	case "fatal":
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(&ErrorResponse{
			Error: "test 500 error",
		}); err != nil {
			h.l.Error(err.Error())
		}
		return
	default:
	}

	pus, err := h.c.GetPaymentsURL(r.Context(), pid)
	switch errors.Cause(err) {
	case nil:
		if err := json.NewEncoder(w).Encode(&Response{
			ApplePayURL:  pus.APayURL,
			GooglePayURL: pus.GPayURL,
		}); err != nil {
			h.l.Error(err.Error())
		}
		return
	case providers.ErrInternalProvider:
		if err := json.NewEncoder(w).Encode(&AppURLResponse{
			AppleAppURL:  appleAppURL,
			GoogleAppURL: androidAppURL,
		}); err != nil {
			h.l.Error(err.Error())
		}
		return
	case providers.ErrNotOK:
		if err := json.NewEncoder(w).Encode(&AppURLResponse{
			AppleAppURL:  appleAppURL,
			GoogleAppURL: androidAppURL,
		}); err != nil {
			h.l.Error(err.Error())
		}
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(&ErrorResponse{
			Error: err.Error(),
		}); err != nil {
			h.l.Error(err.Error())
		}
	}
}

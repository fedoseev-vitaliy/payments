package apay

import (
	"context"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/fedoseev-vitaliy/payments/internal/providers"
	"github.com/fedoseev-vitaliy/payments/internal/utils"
)

type ApplePay struct {
	url    *url.URL
	client *utils.Client
}

type applePayResponse struct {
	PayButtonURL string `json:"some_url"`
}

type applePayError struct {
	Error string `json:"a_err"`
}

func New(cli *utils.Client, u *url.URL) *ApplePay {
	return &ApplePay{
		client: cli,
		url:    u,
	}
}

func (g *ApplePay) GetPayURL(ctx context.Context, productID string) (string, error) {
	res := &applePayResponse{}
	eres := &applePayError{}

	u := *g.url
	q := u.Query()
	q.Set("productID", productID)
	u.RawQuery = q.Encode()

	sc, err := g.client.Get(ctx, &u, res, eres)
	if err != nil {
		return "", errors.Wrapf(providers.ErrInternalProvider, "applePay err: %s", err.Error())
	}

	if sc != http.StatusOK {
		return "", errors.Wrapf(providers.ErrNotOK, "applePay status code:%d", sc)
	}

	return res.PayButtonURL, nil
}

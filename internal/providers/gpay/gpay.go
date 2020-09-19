package gpay

import (
	"context"
	"net/http"
	"net/url"

	"github.com/pkg/errors"

	"github.com/fedoseev-vitaliy/payments/internal/providers"
	"github.com/fedoseev-vitaliy/payments/internal/utils"
)

type GooglePay struct {
	url    *url.URL
	client *utils.Client
}

type googlePayResponse struct {
	PayButtonURL string `json:"p_url"`
}

type googlePayError struct {
	Error string `json:"err"`
}

func New(cli *utils.Client, u *url.URL) *GooglePay {
	return &GooglePay{
		client: cli,
		url:    u,
	}
}

func (g *GooglePay) GetPayURL(ctx context.Context, productID string) (string, error) {
	res := &googlePayResponse{}
	eres := &googlePayError{}

	u := *g.url
	q := u.Query()
	q.Set("productID", productID)
	u.RawQuery = q.Encode()

	sc, err := g.client.Get(ctx, &u, res, eres)
	if err != nil {
		return "", errors.Wrapf(providers.ErrInternalProvider, "googlePay err: %s", err.Error())
	}

	if sc != http.StatusOK {
		return "", errors.Wrapf(providers.ErrNotOK, "googlePay status code:%d", sc)
	}

	return res.PayButtonURL, nil
}

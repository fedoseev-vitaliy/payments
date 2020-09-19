package controller

import (
	"context"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

// PaymentsURLs model to store providers payment urls
type PaymentsURLs struct {
	GPayURL string
	APayURL string
}

// GetPaymentsURL call providers func to get payments urls
func (c *Controller) GetPaymentsURL(ctx context.Context, productID string) (*PaymentsURLs, error) {
	g, gctx := errgroup.WithContext(ctx)

	var gu string
	g.Go(func() error {
		u, err := c.gpay.GetPayURL(gctx, productID)
		if err != nil {
			return errors.WithStack(err)
		}
		gu = u
		return nil
	})

	var au string
	g.Go(func() error {
		u, err := c.apay.GetPayURL(gctx, productID)
		if err != nil {
			return errors.WithStack(err)
		}
		au = u
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, errors.WithStack(err)
	}

	return &PaymentsURLs{
		APayURL: au,
		GPayURL: gu,
	}, nil
}

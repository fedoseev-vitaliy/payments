package providers

import "context"

type Provider interface {
	GetPayURL(ctx context.Context, productID string) (string, error)
}

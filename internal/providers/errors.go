package providers

import "errors"

type ProviderErr error

var (
	ErrInternalProvider = errors.New("internal provider error")
	ErrNotOK            = errors.New("status code not OK")
)

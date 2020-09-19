package controller

import (
	"github.com/fedoseev-vitaliy/payments/internal/providers"
)

// Controller payment providers controller
type Controller struct {
	apay providers.Provider
	gpay providers.Provider
}

// New construct payment provider controller
func New(ap providers.Provider, gp providers.Provider) *Controller {
	return &Controller{gpay: gp, apay: ap}
}

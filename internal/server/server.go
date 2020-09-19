package server

import (
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/fedoseev-vitaliy/payments/internal/controller"
	"github.com/fedoseev-vitaliy/payments/internal/providers/apay"
	"github.com/fedoseev-vitaliy/payments/internal/providers/gpay"
	"github.com/fedoseev-vitaliy/payments/internal/utils"
)

// newRouter construct router
func newRouter(l *logrus.Logger, aPayURL *url.URL, gPayURL *url.URL) http.Handler {
	mux := http.NewServeMux()

	cli := utils.NewClient(time.Second * 5)

	ap := apay.New(cli, aPayURL)
	gp := gpay.New(cli, gPayURL)
	c := controller.New(ap, gp)
	h := NewHandler(l, c)

	mux.HandleFunc("/api/v1/payments/urls", h.GetPaymentsURLs)

	return newPanicRecoveryMiddleware(newLoggerMiddleware(newHeaderMiddleware(mux), l), l)
}

// NewServer construct server with handler
func NewServer(l *logrus.Logger, addr string, aPayURL *url.URL, gPayURL *url.URL) *http.Server {
	r := newRouter(l, aPayURL, gPayURL)

	return &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
}

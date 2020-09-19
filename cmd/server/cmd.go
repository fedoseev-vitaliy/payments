package server

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/fedoseev-vitaliy/payments/internal/providers/apay"
	"github.com/fedoseev-vitaliy/payments/internal/providers/gpay"
	"github.com/fedoseev-vitaliy/payments/internal/server"
	"github.com/fedoseev-vitaliy/payments/internal/utils"
)

var cfg Config

func init() {
	Cmd.Flags().AddFlagSet(cfg.Flags())
}

var Cmd = &cobra.Command{
	Use:   "server",
	Short: "Simple payments API",
	RunE: func(cmd *cobra.Command, args []string) error {
		bindEnv(cmd)
		addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
		l := logrus.New()
		l.SetFormatter(&logrus.JSONFormatter{})

		l.Infof("Starting server: %s", addr)
		apMock := utils.NewTestTLSServer(&apay.MockAPay{})
		aURL, err := url.Parse(apMock.URL)
		if err != nil {
			return errors.WithStack(err)
		}
		gpMock := utils.NewTestTLSServer(&gpay.MockGPay{})
		gURL, err := url.Parse(gpMock.URL)
		if err != nil {
			return errors.WithStack(err)
		}

		srv := server.NewServer(l, addr, aURL, gURL)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-quit

			defer func() {
				apMock.Close()
				gpMock.Close()
			}()

			l.Info("gracefully shutdown server...")
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
			defer cancel()

			if err := srv.Shutdown(ctx); err != nil {
				l.Errorf("failed to gracefully shutdown server. err:%s", err.Error())
			}

			l.Info("server shutdown completed")
		}()

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			return errors.WithStack(err)
		}

		return nil
	},
}

// bindEnv get config values from environment variables
// if no env params than cobra with try to take it from arguments otherwise defaults will be used
func bindEnv(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		envVar := strings.ToUpper(f.Name)

		if val := os.Getenv(envVar); val != "" {
			if err := cmd.Flags().Set(f.Name, val); err != nil {
				logrus.WithError(err).Error("failed to set flag")
			}
		}
	})
}

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/fedoseev-vitaliy/payments/cmd/server"
)

var RootCmd = &cobra.Command{
	Use:   "payments",
	Short: "Simple payments server",
}

func Execute() {
	l := logrus.New()
	l.SetFormatter(&logrus.JSONFormatter{})

	if err := RootCmd.Execute(); err != nil {
		l.Fatalf("something goes wrong. Err:%s", err.Error())
	}
}

func init() {
	RootCmd.AddCommand(server.Cmd)
}

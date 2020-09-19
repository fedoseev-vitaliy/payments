package server

import (
	"github.com/spf13/pflag"
)

type Config struct {
	Host string
	Port int
}

// Flags define default flag set
func (c *Config) Flags() *pflag.FlagSet {
	f := pflag.NewFlagSet("Config", pflag.PanicOnError)

	f.StringVar(&c.Host, "host", "0.0.0.0", "ip host")
	f.IntVar(&c.Port, "port", 80, "port")

	return f
}

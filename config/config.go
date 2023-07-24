package config

import (
	"flag"

	"gomarket/internal/logger"

	"github.com/caarlos0/env/v6"
)

const (
	serverHost  = "localhost:8080"
	accrualHost = ""
	dsn         = "postgres://gomark:gomark@localhost:5433/gomarket?sslmode=disable"
	signKey     = "secret"
)

type Config struct {
	ServerHost  string `env:"RUN_ADDRESS"`
	DSN         string `env:"DATABASE_URI"`
	AccrualHost string `env:"ACCRUAL_SYSTEM_ADDRESS"`
	SignKey     string `env:"KEY"`
}

func GetConfig() *Config {
	cfg := Config{}
	cfg.parseFlags()

	if err := env.Parse(&cfg); err != nil {
		logger.Log.Fatal().Msg("err parse environment variable to server config")
	}

	return &cfg
}

func (c *Config) parseFlags() {
	flag.StringVar(&c.ServerHost, "a", serverHost, "Listen server address (default - :8080)")
	flag.StringVar(&c.DSN, "d", dsn, "URI to database")
	flag.StringVar(&c.AccrualHost, "r", accrualHost, "Accrual system address")
	flag.StringVar(&c.SignKey, "k", signKey, "includes key signature using an algorithm SHA256")
	flag.Parse()
}

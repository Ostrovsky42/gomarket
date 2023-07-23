package main

import (
	"context"
	"gomarket/config"
	"gomarket/internal/logger"
	"os/signal"
	"syscall"
)

func main() {
	logger.InitLogger()

	cfg := config.GetConfig()
	logger.Log.Info().Interface("cfg", cfg).Msg("start with config")

	app := NewApp(cfg)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.Run(ctx)
}

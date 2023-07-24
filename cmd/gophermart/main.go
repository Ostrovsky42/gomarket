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

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	NewApp(cfg).Run(ctx)
}

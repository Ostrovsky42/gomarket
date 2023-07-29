package main

import (
	"context"
	"gomarket/config"
	"gomarket/internal/app"
	"gomarket/internal/logger"
	"os/signal"
	"syscall"
)

func main() {
	logger.InitLogger()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app.NewApp(config.GetConfig()).Run(ctx)
}

package main

import (
	"gomarket/config"
	"gomarket/internal/logger"
)

func main() {
	logger.InitLogger()

	cfg := config.GetConfig()
	logger.Log.Info().Interface("cfg", cfg).Msg("start with config")

	app := NewApp(cfg)

	app.Run()
}

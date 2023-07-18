package main

import (
	"gomarket/config"
	"gomarket/internal/controllers/handlers"
	"gomarket/internal/logger"
	"gomarket/internal/servises/hasher"
	"gomarket/internal/servises/jwt"
	"gomarket/internal/storage/accunts"
	"gomarket/internal/storage/db"
	"gomarket/internal/storage/orders"
	"net/http"
)

type Application struct {
	handlers   *handlers.Handlers
	serverHost string
	signKey    string
}

func NewApp(cfg *config.Config) Application {
	pg, err := db.NewPostgresDB(cfg.DSN)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed init storage")
	}

	accRepo := accunts.NewAccountPG(pg)
	orderRepo := orders.NewOrderPG(pg)
	hashServ := hasher.NewHashGenerator(cfg.SignKey) //todo отдельный ключ
	tokenServ := jwt.NewJWTService(cfg.SignKey, 600)

	return Application{
		handlers:   handlers.NewHandlers(hashServ, accRepo, orderRepo, tokenServ),
		serverHost: cfg.ServerHost,
		signKey:    cfg.SignKey,
	}
}

func (a Application) Run() {
	err := http.ListenAndServe(a.serverHost, a.NewRoutes(a.signKey))
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("Error start serve")
	}
}

func (a Application) Close() {
}

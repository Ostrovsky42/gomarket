package main

import (
	"context"
	"gomarket/config"
	"gomarket/internal/controllers/handlers"
	"gomarket/internal/logger"
	"gomarket/internal/servises/accrual"
	"gomarket/internal/servises/hasher"
	"gomarket/internal/servises/jwt"
	"gomarket/internal/storage/accunts"
	"gomarket/internal/storage/db"
	"gomarket/internal/storage/orders"
	"gomarket/internal/storage/withdraw"
	"net/http"
	"time"
)

type Application struct {
	httpServer *http.Server
	accrualCli *accrual.AccrualProcesser
	pg         *db.Postgres
}

func NewApp(cfg *config.Config) Application {
	pg, err := db.NewPostgresDB(cfg.DSN)
	if err != nil {
		logger.Log.Fatal().Err(err).Msg("failed init storage")
	}

	accRepo := accunts.NewAccountPG(pg)
	orderRepo := orders.NewOrderPG(pg)
	withdrawRepo := withdraw.NewAccountPG(pg)
	hashServ := hasher.NewHashGenerator(cfg.SignKey)
	accrualCli := accrual.NewAccrual(cfg.AccrualHost, orderRepo)
	tokenServ := jwt.NewJWTService(cfg.SignKey, 6000)

	accrualCli.Run()
	return Application{
		httpServer: &http.Server{
			Addr: cfg.ServerHost,
			Handler: NewRoutes(cfg.SignKey, handlers.NewHandlers(
				hashServ,
				accRepo,
				orderRepo,
				withdrawRepo,
				tokenServ,
			))},
		pg:         pg,
		accrualCli: accrualCli,
	}
}

func (a *Application) Run(ctx context.Context) {
	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatal().Err(err).Msg("Error start serve")
		}
	}()
	a.accrualCli.Run()

	<-ctx.Done()
	logger.Log.Info().Msg("server shutting down")

	if err := a.Shutdown(); err != nil {
		logger.Log.Fatal().Err(err).Msg("failed shutting down server")
	}

}

func (a *Application) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return err
	}

	a.pg.Close()

	return nil
}

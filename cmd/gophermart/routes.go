package main

import (
	"compress/gzip"
	"gomarket/internal/controllers/handlers"

	"github.com/go-chi/chi/v5"
	"gomarket/internal/middleware"
)

func NewRoutes(signKey string, handlers *handlers.Handlers) *chi.Mux {
	r := chi.NewRouter()

	zipMW := middleware.NewZipMiddleware(gzip.BestSpeed)

	auth := middleware.NewAuthMiddleware(signKey, 0)

	r.Use(middleware.WithLogging)

	r.Group(func(authRoute chi.Router) {
		authRoute.Use(middleware.WithLogging, zipMW.Zip, zipMW.UnZip)
		authRoute.Post(`/api/user/register`, handlers.RegisterHandler)
		authRoute.Post(`/api/user/login`, handlers.AuthHandler)

	})

	r.Group(func(apiRoute chi.Router) {
		apiRoute.Use(middleware.WithLogging, zipMW.Zip, zipMW.UnZip, auth.Auth)
		apiRoute.Post(`/api/user/orders`, handlers.LoadOrderHandler)
		apiRoute.Get(`/api/user/orders`, handlers.GetOrderHandler)
		apiRoute.Get(`/api/user/balance`, handlers.GetBalance)
		apiRoute.Post(`/api/user/balance/withdraw`, handlers.UsePoints)
		apiRoute.Get(`/api/user/withdrawals`, handlers.UsePointsInfo)
	})

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()

	return r
}

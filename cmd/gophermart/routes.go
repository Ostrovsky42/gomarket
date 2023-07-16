package main

import (
	"compress/gzip"

	"github.com/go-chi/chi/v5"
	"gomarket/internal/middleware"
)

func (a Application) NewRoutes(signKey string) *chi.Mux {
	r := chi.NewRouter()

	zipMW := middleware.NewZipMiddleware(gzip.BestSpeed)

	auth := middleware.NewAuthMiddleware(signKey, 0)

	r.Use(middleware.WithLogging)

	r.Group(func(authRoute chi.Router) {
		authRoute.Use(middleware.WithLogging, zipMW.Zip, zipMW.UnZip)
		authRoute.Post(`/api/user/register`, a.handlers.RegisterHandler)
		authRoute.Post(`/api/user/login`, a.handlers.AuthHandler)

	})

	r.Group(func(apiRoute chi.Router) {
		apiRoute.Use(middleware.WithLogging, zipMW.Zip, zipMW.UnZip, auth.Auth)
	})

	r.NotFoundHandler()
	r.MethodNotAllowedHandler()

	return r
}

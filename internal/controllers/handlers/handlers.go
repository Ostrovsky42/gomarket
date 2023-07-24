package handlers

import (
	"gomarket/internal/repositry"
	"gomarket/internal/servises"
)

type Handlers struct {
	repo *repositry.DataRepositories
	serv *servises.Serv
}

func NewHandlers(
	repo *repositry.DataRepositories,
	serv *servises.Serv,
) *Handlers {
	return &Handlers{
		repo: repo,
		serv: serv,
	}
}

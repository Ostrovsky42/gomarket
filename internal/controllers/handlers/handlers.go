package handlers

import (
	"gomarket/internal/repositry"
	"gomarket/internal/servises/hasher"
	"gomarket/internal/servises/jwt"
)

type Handlers struct {
	repo      repositry.DataRepositories
	hashServ  hasher.HashBuilder
	tokenServ jwt.TokenService
}

func NewHandlers(
	hashServ hasher.HashBuilder,
	repo repositry.DataRepositories,
	tokenServ jwt.TokenService,
) *Handlers {
	return &Handlers{
		repo:      repo,
		hashServ:  hashServ,
		tokenServ: tokenServ,
	}
}

package handlers

import (
	"gomarket/internal/servises/hasher"
	"gomarket/internal/servises/jwt"
	"gomarket/internal/storage/accunts"
)

type Handlers struct {
	accounts  accunts.AccountRepository
	hashServ  hasher.HashBuilder
	tokenServ jwt.TokenService
}

func NewHandlers(
	hashServ hasher.HashBuilder,
	accRepo accunts.AccountRepository,
	tokenServ jwt.TokenService,
) *Handlers {
	return &Handlers{
		accounts:  accRepo,
		hashServ:  hashServ,
		tokenServ: tokenServ,
	}
}

package handlers

import (
	"gomarket/internal/servises/hasher"
	"gomarket/internal/servises/jwt"
	"gomarket/internal/storage/accunts"
	"gomarket/internal/storage/orders"
)

type Handlers struct {
	accounts  accunts.AccountRepository
	orders    orders.OrderRepository
	hashServ  hasher.HashBuilder
	tokenServ jwt.TokenService
}

func NewHandlers(
	hashServ hasher.HashBuilder,
	accRepo accunts.AccountRepository,
	orderRepo orders.OrderRepository,
	tokenServ jwt.TokenService,
) *Handlers {
	return &Handlers{
		accounts:  accRepo,
		orders:    orderRepo,
		hashServ:  hashServ,
		tokenServ: tokenServ,
	}
}

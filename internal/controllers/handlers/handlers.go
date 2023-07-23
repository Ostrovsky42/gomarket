package handlers

import (
	"gomarket/internal/servises/hasher"
	"gomarket/internal/servises/jwt"
	"gomarket/internal/storage/accunts"
	"gomarket/internal/storage/orders"
	"gomarket/internal/storage/withdraw"
)

type Handlers struct {
	accounts  accunts.AccountRepository
	orders    orders.OrderRepository
	withdraw  withdraw.WithDrawRepository
	hashServ  hasher.HashBuilder
	tokenServ jwt.TokenService
}

func NewHandlers(
	hashServ hasher.HashBuilder,
	accRepo accunts.AccountRepository,
	orderRepo orders.OrderRepository,
	withdrawRepo withdraw.WithDrawRepository,
	tokenServ jwt.TokenService,
) *Handlers {
	return &Handlers{
		accounts:  accRepo,
		orders:    orderRepo,
		withdraw:  withdrawRepo,
		hashServ:  hashServ,
		tokenServ: tokenServ,
	}
}

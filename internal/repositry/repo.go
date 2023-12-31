package repositry

import (
	"gomarket/internal/storage/accounts"
	"gomarket/internal/storage/db"
	"gomarket/internal/storage/orders"
	"gomarket/internal/storage/withdraws"
)

type DataRepositories struct {
	Accounts  accounts.AccountRepository
	Orders    orders.OrderRepository
	Withdraws withdraws.WithDrawRepository
}

func NewRepo(pg *db.Postgres) *DataRepositories {
	return &DataRepositories{
		Accounts:  accounts.New(pg),
		Orders:    orders.New(pg),
		Withdraws: withdraws.New(pg),
	}
}

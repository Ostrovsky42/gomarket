// nolint:revive
package mock

import (
	"gomarket/internal/repositry"
	"gomarket/internal/storage_mock"
)

func NewMockRepo(
	account *storage_mock.MockAccountRepository,
	orders *storage_mock.MockOrderRepository,
	withdraw *storage_mock.MockWithDrawRepository,
) *repositry.DataRepositories {
	return &repositry.DataRepositories{
		Accounts:  account,
		Orders:    orders,
		Withdraws: withdraw,
	}
}

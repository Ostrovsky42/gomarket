// nolint:revive
package mocks

import (
	"gomarket/internal/repositry"
)

func NewMockRepo(
	account *MockAccountRepository,
	orders *MockOrderRepository,
	withdraw *MockWithDrawRepository,
) *repositry.DataRepositories {
	return &repositry.DataRepositories{
		Accounts:  account,
		Orders:    orders,
		Withdraws: withdraw,
	}
}

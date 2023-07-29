// nolint:revive
package repositry

import (
	"context"

	"gomarket/internal/entities"
	"gomarket/internal/errors"

	"github.com/stretchr/testify/mock"
)

func NewMockRepo(m *mock.Mock) *DataRepositories {
	return &DataRepositories{
		Accounts:  &AccountRepositoryMock{m},
		Orders:    &OrderRepositoryMock{m},
		Withdraws: &WithDrawRepositoryMock{m},
	}
}

type AccountRepositoryMock struct {
	*mock.Mock
}

func (m *AccountRepositoryMock) CreateAccount(ctx context.Context, login string, hashPass string) (string, *errors.ErrorApp) {
	args := m.Called(ctx, login, hashPass)
	return args.String(0), nil
}

func (m *AccountRepositoryMock) GetAccountByLogin(ctx context.Context, login string) (entities.Account, *errors.ErrorApp) {
	args := m.Called(ctx, login)
	return args.Get(0).(entities.Account), nil
}

func (m *AccountRepositoryMock) GetAccountBalance(ctx context.Context, accountID string) (float64, *errors.ErrorApp) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(float64), nil
}

func (m *AccountRepositoryMock) UpdateAccountBalance(ctx context.Context, accountID string, sum float64) *errors.ErrorApp {
	return nil
}

type OrderRepositoryMock struct {
	*mock.Mock
}

func (m *OrderRepositoryMock) CreateOrder(ctx context.Context, orderID string, accountID string) *errors.ErrorApp {
	return nil
}

func (m *OrderRepositoryMock) GetOrderByID(ctx context.Context, orderID string) (*entities.Order, *errors.ErrorApp) {
	args := m.Called(ctx, orderID)
	return nil, args.Get(1).(*errors.ErrorApp)
}

func (m *OrderRepositoryMock) GetOrdersByAccountID(ctx context.Context, accountID string) ([]entities.Order, *errors.ErrorApp) {
	args := m.Called(ctx, accountID)
	return args.Get(0).([]entities.Order), nil
}

func (m *OrderRepositoryMock) GetOrderIDsForAccrual(ctx context.Context) ([]string, *errors.ErrorApp) {
	args := m.Called(ctx)
	return args.Get(0).([]string), nil
}

func (m *OrderRepositoryMock) UpdateAfterAccrual(ctx context.Context, orderID string, status string, points float64) *errors.ErrorApp {
	args := m.Called(ctx, orderID, status, points)
	return args.Get(0).(*errors.ErrorApp)
}

type WithDrawRepositoryMock struct {
	*mock.Mock
}

func (m *WithDrawRepositoryMock) CreateWithdraw(ctx context.Context, accountID string, orderID string, points float64) *errors.ErrorApp {
	return nil
}

func (m *WithDrawRepositoryMock) GetWithdraw(ctx context.Context, accountID string) ([]entities.Withdraw, *errors.ErrorApp) {
	args := m.Called(ctx, accountID)
	return args.Get(0).([]entities.Withdraw), nil
}

func (m *WithDrawRepositoryMock) GetWithdrawSum(ctx context.Context, accountID string) (*float64, *errors.ErrorApp) {
	args := m.Called(ctx, accountID)
	return args.Get(0).(*float64), nil
}

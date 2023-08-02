// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/storage/accounts/accounts.go

// Package mock_accounts is a generated GoMock package.
package mocks

import (
	context "context"
	entities "gomarket/internal/entities"
	errors "gomarket/internal/errors"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockAccountRepository is a mock of AccountRepository interface.
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository.
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance.
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// CreateAccount mocks base method.
func (m *MockAccountRepository) CreateAccount(ctx context.Context, login, hashPass string) (string, *errors.ErrorApp) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccount", ctx, login, hashPass)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(*errors.ErrorApp)
	return ret0, ret1
}

// CreateAccount indicates an expected call of CreateAccount.
func (mr *MockAccountRepositoryMockRecorder) CreateAccount(ctx, login, hashPass interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccount", reflect.TypeOf((*MockAccountRepository)(nil).CreateAccount), ctx, login, hashPass)
}

// GetAccountBalance mocks base method.
func (m *MockAccountRepository) GetAccountBalance(ctx context.Context, accountID string) (float64, *errors.ErrorApp) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountBalance", ctx, accountID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(*errors.ErrorApp)
	return ret0, ret1
}

// GetAccountBalance indicates an expected call of GetAccountBalance.
func (mr *MockAccountRepositoryMockRecorder) GetAccountBalance(ctx, accountID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountBalance", reflect.TypeOf((*MockAccountRepository)(nil).GetAccountBalance), ctx, accountID)
}

// GetAccountByLogin mocks base method.
func (m *MockAccountRepository) GetAccountByLogin(ctx context.Context, login string) (entities.Account, *errors.ErrorApp) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountByLogin", ctx, login)
	ret0, _ := ret[0].(entities.Account)
	ret1, _ := ret[1].(*errors.ErrorApp)
	return ret0, ret1
}

// GetAccountByLogin indicates an expected call of GetAccountByLogin.
func (mr *MockAccountRepositoryMockRecorder) GetAccountByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountByLogin", reflect.TypeOf((*MockAccountRepository)(nil).GetAccountByLogin), ctx, login)
}

// UpdateAccountBalance mocks base method.
func (m *MockAccountRepository) UpdateAccountBalance(ctx context.Context, accountID string, sum float64) *errors.ErrorApp {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccountBalance", ctx, accountID, sum)
	ret0, _ := ret[0].(*errors.ErrorApp)
	return ret0
}

// UpdateAccountBalance indicates an expected call of UpdateAccountBalance.
func (mr *MockAccountRepositoryMockRecorder) UpdateAccountBalance(ctx, accountID, sum interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccountBalance", reflect.TypeOf((*MockAccountRepository)(nil).UpdateAccountBalance), ctx, accountID, sum)
}
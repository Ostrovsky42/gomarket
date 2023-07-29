package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gomarket/internal/accountctx"
	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/repositry"
	"gomarket/internal/servises"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterHandler(t *testing.T) {
	m := &mock.Mock{}
	serv := servises.NewService("secret")
	m.On("CreateAccount", context.Background(), "testuser", serv.Hash.GetHash("testpass")).
		Return("accountID123")
	r := repositry.NewMockRepo(m)
	h := NewHandlers(
		r,
		serv,
	)

	requestBody := []byte(`{"login": "testuser", "password": "testpass"}`)

	req, err := http.NewRequest("POST", "/api/user/register", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.RegisterHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	token, _ := serv.Token.GenerateToken("accountID123")
	assert.Equal(t, "Bearer "+token, rr.Header().Get("Authorization"))
}

func TestAuthHandler(t *testing.T) {
	m := &mock.Mock{}
	serv := servises.NewService("secret")
	m.On("GetAccountByLogin", context.Background(), "testuser").
		Return(entities.Account{
			ID:       "accountID123",
			HashPass: serv.Hash.GetHash("testpass"),
		}, nil)

	r := repositry.NewMockRepo(m)

	h := NewHandlers(
		r,
		serv,
	)

	requestBody := []byte(`{"login": "testuser", "password": "testpass"}`)

	req, err := http.NewRequest("POST", "/api/user/auth", bytes.NewBuffer(requestBody))
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.AuthHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	token, _ := serv.Token.GenerateToken("accountID123")
	assert.Equal(t, "Bearer "+token, rr.Header().Get("Authorization"))
}

var float = 3.14

func TestBalanceHandler(t *testing.T) {
	m := &mock.Mock{}
	ctx := accountctx.WithAccountID(context.Background(), "accountID123")
	serv := servises.NewService("secret")
	m.On("GetAccountBalance", ctx, "accountID123").Return(float, nil)
	m.On("GetWithdrawSum", ctx, "accountID123").Return(&float, nil)
	r := repositry.NewMockRepo(m)
	h := NewHandlers(
		r,
		serv,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/user/balance", nil)

	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.GetBalance(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `{"current":3.14,"withdrawn":3.14}`, rr.Body.String())
}

func TestLoadOrderHandler(t *testing.T) {
	m := &mock.Mock{}
	ctx := accountctx.WithAccountID(context.Background(), "accountID123")
	serv := servises.NewService("secret")
	m.On("GetOrderByID", ctx, "4111111111111111").Return(nil, errors.NewErrNotFound())
	m.On("CreateOrder", ctx, "4111111111111111", "accountID123").Return(nil)
	r := repositry.NewMockRepo(m)
	h := NewHandlers(
		r,
		serv,
	)
	requestBody := []byte(`4111111111111111`)

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/user/orders", bytes.NewBuffer(requestBody))

	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.LoadOrderHandler(rr, req)

	assert.Equal(t, http.StatusAccepted, rr.Code)
}

func TestGetOrderHandler(t *testing.T) {
	m := &mock.Mock{}
	ctx := accountctx.WithAccountID(context.Background(), "accountID123")
	serv := servises.NewService("secret")
	m.On("GetOrdersByAccountID", ctx, "accountID123").
		Return([]entities.Order{{
			ID:         "4111111111111111",
			AccountID:  "accountID123",
			Status:     entities.New,
			UploadedAt: time.Time{},
		}}, nil)
	r := repositry.NewMockRepo(m)
	h := NewHandlers(
		r,
		serv,
	)
	req, err := http.NewRequestWithContext(ctx, "GET", "/api/user/orders", nil)

	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.GetOrderHandler(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUsePointsHandler(t *testing.T) {
	m := &mock.Mock{}
	ctx := accountctx.WithAccountID(context.Background(), "accountID123")
	serv := servises.NewService("secret")
	m.On("UpdateAccountBalance", ctx, "accountID123", transferToNegative(float))
	m.On("CreateWithdraw", ctx, "accountID123", "4111111111111111", float)

	r := repositry.NewMockRepo(m)
	h := NewHandlers(
		r,
		serv,
	)
	requestBody := []byte(`{"order":"4111111111111111","sum":3.14}`)

	req, err := http.NewRequestWithContext(ctx, "POST", "/api/user/balance/withdraw", bytes.NewBuffer(requestBody))

	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.UsePoints(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestUsePointsInfoHandler(t *testing.T) {
	m := &mock.Mock{}

	ctx := accountctx.WithAccountID(context.Background(), "accountID123")
	serv := servises.NewService("secret")
	m.On("GetWithdraw", ctx, "accountID123").
		Return([]entities.Withdraw{{
			OrderID:     "4111111111111111",
			AccountID:   "accountID123",
			Sum:         float,
			ProcessedAt: time.Time{},
		}}, nil)

	r := repositry.NewMockRepo(m)
	h := NewHandlers(
		r,
		serv,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", "/api/user/withdrawals", nil)

	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	h.UsePointsInfo(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, `[{"order":"4111111111111111","sum":3.14,"processed_at":"0001-01-01T00:00:00Z"}]`, rr.Body.String())
}

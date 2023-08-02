package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/mocks"
	"gomarket/internal/repositry"
	"gomarket/internal/servises"
)

func TestHandlers_RegisterHandler(t *testing.T) {
	var ctrl *gomock.Controller
	serv := servises.NewService("secret")
	ctx := context.Background()
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		repo     *repositry.DataRepositories
		args     args
		wantCode int
	}{
		{
			name: "Successful registration",
			repo: func() *repositry.DataRepositories {
				ctrl = gomock.NewController(t)
				re := mocks.NewMockAccountRepository(ctrl)
				re.EXPECT().CreateAccount(ctx, "testuser", serv.Hash.GetHash("testpass")).
					Return("accountID123", nil).Times(1)
				return mocks.NewMockRepo(re, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"login": "testuser", "password": "testpass"}`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Conflict registration",
			repo: func() *repositry.DataRepositories {
				ctrl = gomock.NewController(t)
				re := mocks.NewMockAccountRepository(ctrl)
				re.EXPECT().CreateAccount(ctx, "testuser", serv.Hash.GetHash("testpass")).
					Return("", errors.NewErrUniquenessViolation(nil)).Times(1)
				return mocks.NewMockRepo(re, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"login": "testuser", "password": "testpass"}`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusConflict,
		},
		{
			name: "Failed validation",
			repo: func() *repositry.DataRepositories {
				ctrl = gomock.NewController(t)
				re := mocks.NewMockAccountRepository(ctrl)
				return mocks.NewMockRepo(re, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"login": "testuser=", "password": "testpass"}`)

					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				repo: tt.repo,
				serv: serv,
			}
			h.RegisterHandler(tt.args.w, tt.args.r)
		})
		assert.Equal(t, tt.wantCode, tt.args.w.Code)
		if http.StatusOK == tt.args.w.Code {
			wantToken, _ := serv.Token.GenerateToken("accountID123")
			token := tt.args.w.Header().Get("Authorization")
			assert.Equal(t, "Bearer "+wantToken, token)
		}
	}
	ctrl.Finish()
}

func TestHandlers_AuthHandler(t *testing.T) {
	var ctrl *gomock.Controller
	serv := servises.NewService("secret")
	ctx := context.Background()
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		repo     *repositry.DataRepositories
		args     args
		wantCode int
	}{
		{
			name: "Successful auth",
			repo: func() *repositry.DataRepositories {
				ctrl = gomock.NewController(t)
				re := mocks.NewMockAccountRepository(ctrl)
				re.EXPECT().GetAccountByLogin(ctx, "testuser").
					Return(entities.Account{
						ID:       "accountID123",
						HashPass: serv.Hash.GetHash("testpass"),
					}, nil).Times(1)
				return mocks.NewMockRepo(re, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"login": "testuser", "password": "testpass"}`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "Wrong pass",
			repo: func() *repositry.DataRepositories {
				ctrl = gomock.NewController(t)
				re := mocks.NewMockAccountRepository(ctrl)
				re.EXPECT().GetAccountByLogin(ctx, "testuser").
					Return(entities.Account{
						ID:       "accountID123",
						HashPass: serv.Hash.GetHash("wrongPassword"),
					}, nil).Times(1)
				return mocks.NewMockRepo(re, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"login": "testuser", "password": "testpass"}`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusUnauthorized,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				repo: tt.repo,
				serv: serv,
			}
			h.AuthHandler(tt.args.w, tt.args.r)
		})
		assert.Equal(t, tt.wantCode, tt.args.w.Code)
		if http.StatusOK == tt.args.w.Code {
			wantToken, _ := serv.Token.GenerateToken("accountID123")
			token := tt.args.w.Header().Get("Authorization")
			assert.Equal(t, "Bearer "+wantToken, token)
		}
	}
	ctrl.Finish()
}

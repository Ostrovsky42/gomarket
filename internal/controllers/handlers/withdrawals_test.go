package handlers

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"gomarket/internal/accountctx"
	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/mocks"
	"gomarket/internal/repositry"
	"gomarket/internal/servises"
)

const sum = 3.14

func TestHandlers_UsePoints(t *testing.T) {
	serv := servises.NewService("secret")
	ctx := accountctx.WithAccountID(context.Background(), "accountID123")
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
			name: "Successful Use Points",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				acc := mocks.NewMockAccountRepository(ctrl)
				acc.EXPECT().UpdateAccountBalance(ctx, "accountID123", transferToNegative(sum)).Return(nil).AnyTimes()
				with := mocks.NewMockWithDrawRepository(ctrl)
				with.EXPECT().CreateWithdraw(ctx, "accountID123", "4111111111111111", sum).Return(nil).AnyTimes()
				return mocks.NewMockRepo(acc, nil, with)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"order":"4111111111111111","sum":3.14}`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusOK,
		},
		{
			name: "InsufficientFunds",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				acc := mocks.NewMockAccountRepository(ctrl)
				acc.EXPECT().UpdateAccountBalance(ctx, "accountID123", transferToNegative(sum)).Return(errors.NewErrInsufficientFunds()).AnyTimes()
				return mocks.NewMockRepo(acc, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"order":"4111111111111111","sum":3.14}`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusPaymentRequired,
		},
		{
			name: "Unauthorized",
			repo: func() *repositry.DataRepositories {
				return mocks.NewMockRepo(nil, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`{"order":"4111111111111111","sum":3.14}`)
					req, _ := http.NewRequestWithContext(context.Background(), "", "", bytes.NewBuffer(requestBody))
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
			h.UsePoints(tt.args.w, tt.args.r)
		})
		assert.Equal(t, tt.wantCode, tt.args.w.Code)
	}
}

func TestHandlers_UsePointsInfo(t *testing.T) {
	serv := servises.NewService("secret")
	ctx := accountctx.WithAccountID(context.Background(), "accountID123")
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		repo     *repositry.DataRepositories
		args     args
		wantCode int
		wantBody string
	}{
		{
			name: "Successful Use Points Info ",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				with := mocks.NewMockWithDrawRepository(ctrl)
				with.EXPECT().GetWithdraw(ctx, "accountID123").
					Return([]entities.Withdraw{
						{
							OrderID: "1",
							Sum:     1.23,
						}, {
							OrderID: "2",
							Sum:     2.34,
						},
					}, nil).AnyTimes()
				return mocks.NewMockRepo(nil, nil, with)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req, _ := http.NewRequestWithContext(ctx, "", "", nil)
					return req
				}(),
			},
			wantBody: `[{"order":"1","sum":1.23,"processed_at":"0001-01-01T00:00:00Z"},{"order":"2","sum":2.34,"processed_at":"0001-01-01T00:00:00Z"}]`,
			wantCode: http.StatusOK,
		},
		{
			name: "No content Use Points Info ",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				with := mocks.NewMockWithDrawRepository(ctrl)
				with.EXPECT().GetWithdraw(ctx, "accountID123").
					Return([]entities.Withdraw{}, nil).AnyTimes()
				return mocks.NewMockRepo(nil, nil, with)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req, _ := http.NewRequestWithContext(ctx, "", "", nil)
					return req
				}(),
			},
			wantCode: http.StatusNoContent,
		},
		{
			name: "Unauthorized",
			repo: func() *repositry.DataRepositories {
				return mocks.NewMockRepo(nil, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req, _ := http.NewRequestWithContext(context.Background(), "", "", nil)
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
			h.UsePointsInfo(tt.args.w, tt.args.r)
		})
		assert.Equal(t, tt.wantCode, tt.args.w.Code)
		assert.Equal(t, tt.wantBody, tt.args.w.Body.String())
	}
}

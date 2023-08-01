package handlers

import (
	"bytes"
	"context"
	"gomarket/internal/repositry/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"gomarket/internal/accountctx"
	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/repositry"
	"gomarket/internal/servises"
	"gomarket/internal/storage_mock"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandlers_GetOrderHandler(t *testing.T) {
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
			name: "Successful get order",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				order := storage_mock.NewMockOrderRepository(ctrl)
				order.EXPECT().GetOrdersByAccountID(ctx, "accountID123").
					Return([]entities.Order{
						{
							ID:     "378282246310005",
							Status: entities.New,
						}, {
							ID:     "6011111111111117",
							Status: entities.New,
						},
					}, nil).AnyTimes()
				return mock.NewMockRepo(nil, order, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req, _ := http.NewRequestWithContext(ctx, "", "", nil)
					return req
				}(),
			},
			wantCode: http.StatusOK,
			wantBody: `[{"number":"378282246310005","status":"NEW","uploaded_at":"0001-01-01T00:00:00Z"},{"number":"6011111111111117","status":"NEW","uploaded_at":"0001-01-01T00:00:00Z"}]`,
		},
		{
			name: "Successful get order",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				order := storage_mock.NewMockOrderRepository(ctrl)
				order.EXPECT().GetOrdersByAccountID(ctx, "accountID123").Return([]entities.Order{}, nil).AnyTimes()
				return mock.NewMockRepo(nil, order, nil)
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				repo: tt.repo,
				serv: serv,
			}
			h.GetOrderHandler(tt.args.w, tt.args.r)
		})
		assert.Equal(t, tt.wantCode, tt.args.w.Code)
		assert.Equal(t, tt.wantBody, tt.args.w.Body.String())
	}
}

func TestHandlers_LoadOrderHandler(t *testing.T) {
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
			name: "Successful Load order",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				order := storage_mock.NewMockOrderRepository(ctrl)
				order.EXPECT().GetOrderByID(ctx, "30569309025904").Return(nil, errors.NewErrNotFound()).AnyTimes()
				order.EXPECT().CreateOrder(ctx, "30569309025904", "accountID123").Return(nil).AnyTimes()
				return mock.NewMockRepo(nil, order, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`30569309025904`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusAccepted,
		},
		{
			name: "Conflict Load order",
			repo: func() *repositry.DataRepositories {
				ctrl := gomock.NewController(t)
				defer ctrl.Finish()
				order := storage_mock.NewMockOrderRepository(ctrl)
				order.EXPECT().GetOrderByID(ctx, "30569309025904").Return(&entities.Order{AccountID: "Conflict"}, nil).AnyTimes()
				return mock.NewMockRepo(nil, order, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`30569309025904`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusConflict,
		},
		{
			name: "Failed validate order",
			repo: func() *repositry.DataRepositories {
				return mock.NewMockRepo(nil, nil, nil)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					requestBody := []byte(`555`)
					req, _ := http.NewRequestWithContext(ctx, "", "", bytes.NewBuffer(requestBody))
					return req
				}(),
			},
			wantCode: http.StatusUnprocessableEntity,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{
				repo: tt.repo,
				serv: serv,
			}
			h.LoadOrderHandler(tt.args.w, tt.args.r)
		})
		assert.Equal(t, tt.wantCode, tt.args.w.Code)
	}
}

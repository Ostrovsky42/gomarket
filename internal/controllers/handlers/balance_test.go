package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"gomarket/internal/accountctx"
	"gomarket/internal/mocks"
	"gomarket/internal/repositry"
	"gomarket/internal/servises"
)

var float = 3.14

func TestHandlers_GetBalance(t *testing.T) {
	var ctrl *gomock.Controller
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
			name: "Successful Received Balance",
			repo: func() *repositry.DataRepositories {
				ctrl = gomock.NewController(t)
				acc := mocks.NewMockAccountRepository(ctrl)
				acc.EXPECT().GetAccountBalance(ctx, "accountID123").Return(float, nil).Times(1)
				with := mocks.NewMockWithDrawRepository(ctrl)
				with.EXPECT().GetWithdrawSum(ctx, "accountID123").Return(&float, nil).Times(1)
				return mocks.NewMockRepo(acc, nil, with)
			}(),
			args: args{
				w: httptest.NewRecorder(),
				r: func() *http.Request {
					req, _ := http.NewRequestWithContext(ctx, "", "", nil)
					return req
				}(),
			},
			wantCode: http.StatusOK,
			wantBody: `{"current":3.14,"withdrawn":3.14}`,
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
			h.GetBalance(tt.args.w, tt.args.r)
		})
		assert.Equal(t, tt.wantCode, tt.args.w.Code)
		assert.Equal(t, tt.wantBody, tt.args.w.Body.String())
	}
	ctrl.Finish()
}

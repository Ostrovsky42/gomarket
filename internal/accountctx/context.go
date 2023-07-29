package accountctx

import (
	"context"
	"gomarket/internal/errors"
)

type ctxKey string

const accountID ctxKey = "account_id"

func WithAccountID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, accountID, id)
}

func GetAccountID(ctx context.Context) (string, *errors.ErrorApp) {
	if temp := ctx.Value(accountID); temp != nil {
		if id, ok := temp.(string); ok {
			return id, nil
		}
	}

	return "", errors.NewError("account_id is empty", nil)
}

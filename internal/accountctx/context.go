package accountctx

import (
	"context"
	"gomarket/internal/errors"
)

type accountIDKey struct{}

func WithAccountID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, accountIDKey{}, id)
}

func GetAccountID(ctx context.Context) (string, *errors.ErrorApp) {
	if temp := ctx.Value(accountIDKey{}); temp != nil {
		if id, ok := temp.(string); ok {
			return id, nil
		}
	}

	return "", errors.NewError("account_id is empty", nil)
}

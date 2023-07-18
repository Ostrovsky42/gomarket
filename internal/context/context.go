package context

import "context"

const accountID = "account_id"

func WithAccountID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, accountID, id)
}

func AccountID(ctx context.Context) string {
	if temp := ctx.Value(accountID); temp != nil {
		if id, ok := temp.(string); ok {
			return id
		}
	}

	return ""
}

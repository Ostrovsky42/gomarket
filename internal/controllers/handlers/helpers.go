package handlers

import (
	contextStd "context"
	"net/http"

	"gomarket/internal/context"
	"gomarket/internal/errors"
)

const (
	ContentType = "Content-Type"
	JSON        = "application/json"
)

func setJSONContentType(w http.ResponseWriter) {
	w.Header().Set(ContentType, JSON)
}

func getAccountID(ctx contextStd.Context) (string, *errors.ErrorApp) {
	if accountID := context.AccountID(ctx); accountID != "" {
		return accountID, nil
	}

	return "", errors.NewError("account_id is empty", nil)
}

func transferToNegative(val float64) float64 {
	return -1 * val
}

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

func transferToNegativeCoins(val float64) int {
	return int(-100 * val)
}

func transferToCoins(val float64) int {
	return int(100 * val)
}

func transferFromCoins(coins int) float64 {
	return float64(coins) / 100
}

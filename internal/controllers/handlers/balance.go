package handlers

import (
	"encoding/json"
	"gomarket/internal/logger"
	"net/http"
)

type BalanceResponse struct {
	Current   float64 `json:"current"`
	Withdrawn *int    `json:"withdrawn,omitempty"`
}

func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID, errApp := getAccountID(ctx)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account_id")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	current, errApp := h.accounts.GetAccountBalance(ctx, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account balance")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	withdrawSum, errApp := h.withdraw.GetWithdrawSum(ctx, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get withdraw sum")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	balanceResponse := BalanceResponse{
		Current:   transferFromCoins(current),
		Withdrawn: withdrawSum,
	}

	jsonData, err := json.Marshal(balanceResponse)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to marshal JSON response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

package handlers

import (
	"gomarket/internal/context"
	"gomarket/internal/logger"
	"net/http"
)

type BalanceResponse struct {
	Current   float64  `json:"current"`
	Withdrawn *float64 `json:"withdrawn,omitempty"`
}

func (h *Handlers) GetBalance(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID, errApp := context.GetAccountID(ctx)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account_id")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	current, errApp := h.repo.Accounts.GetAccountBalance(ctx, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account balance")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	withdrawSum, errApp := h.repo.Withdraws.GetWithdrawSum(ctx, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get withdraw sum")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	writeOKWithJSON(w, BalanceResponse{
		Current:   current,
		Withdrawn: withdrawSum,
	})
}

package handlers

import (
	"encoding/json"
	"gomarket/internal/context"
	"gomarket/internal/entities"
	"gomarket/internal/errors"
	"gomarket/internal/logger"

	"net/http"
)

func (h *Handlers) UsePoints(w http.ResponseWriter, r *http.Request) {
	var req entities.Withdraw
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed decode body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	ctx := r.Context()
	accountID, errApp := context.GetAccountID(ctx)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account_id")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	if errApp = ValidateUsePoints(req.OrderID, req.Sum); errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	errApp = h.repo.Accounts.UpdateAccountBalance(r.Context(), accountID, transferToNegative(req.Sum))
	if errApp != nil {
		if errApp.Description() == errors.InsufficientFunds {
			w.WriteHeader(http.StatusPaymentRequired)

			return
		}

		logger.Log.Error().Err(errApp).Msg("failed update balance")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	errApp = h.repo.Withdraws.CreateWithdraw(ctx, accountID, req.OrderID, req.Sum)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed create withdraw")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) UsePointsInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID, errApp := context.GetAccountID(ctx)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account_id")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	withdraw, errApp := h.repo.Withdraws.GetWithdraw(ctx, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get withdrawals")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(withdraw) == 0 {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	writeOKWithJSON(w, withdraw)
}

func transferToNegative(val float64) float64 {
	return -1 * val
}

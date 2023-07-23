package handlers

import (
	"encoding/json"
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
	accountID, errApp := getAccountID(ctx)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account_id")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	if errApp = ValidateLoadOrder(req.OrderID); errApp != nil { //todo  use points check negative
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	errApp = h.accounts.UpdateAccountBalance(r.Context(), accountID, transferToNegative(req.Sum))
	if errApp != nil {
		if errApp.Description() == errors.InsufficientFunds {
			w.WriteHeader(http.StatusPaymentRequired)

			return
		}

		logger.Log.Error().Err(errApp).Msg("failed update balance")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	errApp = h.withdraw.CreateWithdraw(ctx, accountID, req.OrderID, req.Sum)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed create withdraw")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	//todo create order?
}

func (h *Handlers) UsePointsInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID, errApp := getAccountID(ctx)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account_id")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	withdraw, errApp := h.withdraw.GetWithdraw(ctx, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get withdrawals")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(withdraw) == 0 {
		w.WriteHeader(http.StatusNoContent)

		return
	}

	jsonData, err := json.Marshal(withdraw)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to marshal JSON response")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	setJSONContentType(w)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

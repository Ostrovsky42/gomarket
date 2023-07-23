package handlers

import (
	"encoding/json"
	"gomarket/internal/errors"
	"gomarket/internal/logger"
	"net/http"
	"strconv"
)

func (h *Handlers) LoadOrderHandler(w http.ResponseWriter, r *http.Request) {
	var req int
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

	orderID := strconv.Itoa(req)
	if errApp = ValidateLoadOrder(orderID); errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	order, errApp := h.orders.GetOrderByID(ctx, orderID)
	if errApp != nil {
		if errApp.Description() != errors.NotFound {
			logger.Log.Error().Err(errApp).Msg("failed get order")
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
	}

	if order != nil {
		if order.AccountID != accountID {
			w.WriteHeader(http.StatusConflict)

			return
		}
		w.WriteHeader(http.StatusOK)

		return
	}

	errApp = h.orders.CreateOrder(ctx, orderID, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed create order")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h *Handlers) GetOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	accountID, errApp := getAccountID(ctx)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get account_id")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	orders, errApp := h.orders.GetOrdersByAccountID(ctx, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed get orders by account_id")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)

		return

	}

	jsonData, err := json.Marshal(orders)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed to marshal JSON response")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	setJSONContentType(w)

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

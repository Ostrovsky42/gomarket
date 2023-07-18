package handlers

import (
	"gomarket/internal/context"
	"gomarket/internal/errors"
	"gomarket/internal/logger"
	"io"
	"net/http"
)

func (h *Handlers) LoadOrderHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed read body")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	accountID := context.AccountID(ctx)

	orderID := string(body)
	if errApp := ValidateLoadOrder(orderID); errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusUnprocessableEntity)

		return
	}

	order, errApp := h.orders.GetOrdersByID(ctx, orderID)
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
		w.WriteHeader(http.StatusOK) /*
			200 — номер заказа уже был загружен этим пользователем;
			202 — новый номер заказа принят в обработку;
			409 — номер заказа уже был загружен другим пользователем
		*/

		return
	}

	logger.Log.Debug().Interface(orderID, order).Msg("oredr")

	errApp = h.orders.CreateOrder(ctx, orderID, accountID)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed create order")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

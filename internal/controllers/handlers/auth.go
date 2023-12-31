package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gomarket/internal/errors"
	"gomarket/internal/logger"
)

type AccountAuth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func (h *Handlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req AccountAuth
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed decode body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if errApp := ValidationAuth(req.Login, req.Password); errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	id, errApp := h.repo.Accounts.CreateAccount(r.Context(), req.Login, h.serv.Hash.GetHash(req.Password))
	if errApp != nil {
		if errApp.Description() == errors.UniquenessViolation {
			w.WriteHeader(http.StatusConflict)

			return
		}
		logger.Log.Error().Err(errApp).Msg("failed create account")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if err = h.setJWT(w, id); err != nil {
		logger.Log.Error().Err(err).Msg("failed generate jwt token")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) AuthHandler(w http.ResponseWriter, r *http.Request) {
	var req AccountAuth
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed decode body")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	if errApp := ValidationAuth(req.Login, req.Password); errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	account, errApp := h.repo.Accounts.GetAccountByLogin(r.Context(), req.Login)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed create account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hash := h.serv.Hash.GetHash(req.Password)
	if account.HashPass != hash {
		logger.Log.Info().Msg("failed compare hash password")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	if err = h.setJWT(w, account.ID); err != nil {
		logger.Log.Error().Err(err).Msg("failed generate jwt token")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) setJWT(w http.ResponseWriter, accountID string) error {
	jwt, err := h.serv.Token.GenerateToken(accountID)
	if err != nil {
		return err
	}
	w.Header().Set("Authorization", fmt.Sprintf("Bearer %s", jwt))

	return nil
}

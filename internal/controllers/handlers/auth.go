package handlers

import (
	"context"
	"encoding/json"
	"gomarket/internal/errors"
	"gomarket/internal/logger"
	"net/http"
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

	if errApp := ValidationAuth(req.Login, req.Password); err != nil {
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	id, errApp := h.accounts.CreateAccount(context.Background(), req.Login, h.hashServ.GetHash(req.Password))
	if errApp != nil {
		if errApp.Description() == errors.UniquenessViolation {
			w.WriteHeader(http.StatusConflict)

			return
		}
		logger.Log.Error().Err(errApp).Msg("failed create account")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	jwt, err := h.tokenServ.GenerateToken(id)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed generate jwt token")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	setJWT(w, jwt)

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

	if errApp := ValidationAuth(req.Login, req.Password); err != nil {
		logger.Log.Error().Err(errApp).Msg("failed validation")
		w.WriteHeader(http.StatusBadRequest)

		return
	}

	account, errApp := h.accounts.GetAccountByLogin(context.Background(), req.Login)
	if errApp != nil {
		logger.Log.Error().Err(errApp).Msg("failed create account")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if account.HashPass == h.hashServ.GetHash(req.Password) {
		logger.Log.Info().Msg("failed compare hash password")
		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	jwt, err := h.tokenServ.GenerateToken(account.ID)
	if err != nil {
		logger.Log.Error().Err(err).Msg("failed generate jwt token")
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	setJWT(w, jwt)

	w.WriteHeader(http.StatusOK)
}

func setJWT(w http.ResponseWriter, jwt string) {
	cookie := &http.Cookie{
		Name:     "jwt",
		Value:    jwt,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	}

	http.SetCookie(w, cookie)
}

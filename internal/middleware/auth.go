package middleware

import (
	"gomarket/internal/accountctx"
	"gomarket/internal/servises/jwt"
	"net/http"
	"strings"
)

const Bearer = "Bearer "

type Auth struct {
	tokenService jwt.TokenService
}

func NewAuthMiddleware(secretKey string, ttlSec int) *Auth {
	return &Auth{tokenService: jwt.NewJWTService(secretKey, ttlSec)}
}

func (a *Auth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rawJWT := r.Header.Get("Authorization")
		token := strings.TrimPrefix(rawJWT, Bearer)
		accountID, err := a.tokenService.VerifyToken(token)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		ctx := accountctx.WithAccountID(r.Context(), accountID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package middleware

import (
	"gomarket/internal/logger"
	"gomarket/internal/servises/jwt"
	"net/http"
)

type Auth struct {
	tokenService jwt.TokenService
}

func NewAuthMiddleware(secretKey string, ttlSec int) *Auth {
	return &Auth{tokenService: jwt.NewJWTService(secretKey, ttlSec)}
}

func (a *Auth) Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return
		}

		accountID, err := a.tokenService.VerifyToken(cookie.Value)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)

			return

		}
		logger.Log.Info().Msg(accountID) //todo context

		next.ServeHTTP(w, r)
	})
}

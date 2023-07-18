package middleware

import (
	"gomarket/internal/context"
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

		ctx := context.WithAccountID(r.Context(), accountID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

package jwt

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"gomarket/internal/errors"
)

type TokenService interface {
	GenerateToken(accountID string) (string, error)
	VerifyToken(tokenString string) (string, error)
}

type ServiceJWT struct {
	secretKey       []byte
	tokenExpiration time.Duration
}

func NewJWTService(secretKey string, tokenExpirationSec int) TokenService {
	return &ServiceJWT{
		secretKey:       []byte(secretKey),
		tokenExpiration: time.Duration(tokenExpirationSec) * time.Second,
	}
}

func (s *ServiceJWT) GenerateToken(accountID string) (string, error) {
	claims := jwt.MapClaims{
		"account_id": accountID,
		"exp":        time.Now().Add(s.tokenExpiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *ServiceJWT) VerifyToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.NewErrUnauthorized()
		}

		return s.secretKey, nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.NewErrUnauthorized()
	}

	accountID, ok := claims["account_id"].(string)
	if !ok {
		return "", errors.NewErrUnauthorized()
	}

	return accountID, nil
}

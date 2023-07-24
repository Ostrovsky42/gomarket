package jwt

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const secretKey = "secret-key"

func TestGenerateToken(t *testing.T) {
	tokenExpirationSec := 3600

	tokenService := NewJWTService(secretKey, tokenExpirationSec)

	accountID := "test-account-id"
	token, err := tokenService.GenerateToken(accountID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	verifiedAccountID, err := tokenService.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, accountID, verifiedAccountID)
}

func TestVerifyInvalidToken(t *testing.T) {
	secretKey := "secret-key"
	tokenExpirationSec := 3600

	tokenService := NewJWTService(secretKey, tokenExpirationSec)

	invalidToken := "invalid-token"
	_, err := tokenService.VerifyToken(invalidToken)
	assert.Error(t, err)
}

func TestTokenExpired(t *testing.T) {
	secretKey := "secret-key"
	tokenExpirationSec := 1

	tokenService := NewJWTService(secretKey, tokenExpirationSec)

	accountID := "test-account-id"
	token, err := tokenService.GenerateToken(accountID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	time.Sleep(2 * time.Second)

	_, err = tokenService.VerifyToken(token)
	assert.Error(t, err)
}

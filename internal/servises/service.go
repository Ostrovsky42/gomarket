package servises

import (
	"gomarket/internal/servises/hasher"
	"gomarket/internal/servises/jwt"
)

const tokenLife = 600

type Serv struct {
	Hash  hasher.HashBuilder
	Token jwt.TokenService
}

func NewService(signKey string) *Serv {
	return &Serv{
		Hash:  hasher.NewHashGenerator(signKey),
		Token: jwt.NewJWTService(signKey, tokenLife),
	}
}

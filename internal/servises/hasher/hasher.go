package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

type HashBuilder interface {
	GetHash(data string) string
}

type HashGenerator struct {
	signKey []byte
}

func NewHashGenerator(key string) *HashGenerator {
	return &HashGenerator{
		signKey: []byte(key),
	}
}

func (h *HashGenerator) GetHash(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	signedHash := hash.Sum(h.signKey)
	return hex.EncodeToString(signedHash)
}

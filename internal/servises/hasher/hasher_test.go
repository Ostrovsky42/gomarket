package hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetHash(t *testing.T) {
	key := "secret-key"
	hashGenerator := NewHashGenerator(key)

	data := "pass"
	expectedHash := "7365637265742d6b6579d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1"

	actualHash := hashGenerator.GetHash(data)
	assert.Equal(t, expectedHash, actualHash)
}

func TestEmptyData(t *testing.T) {
	key := "secret-key"
	hashGenerator := NewHashGenerator(key)

	data := ""
	expectedHash := "7365637265742d6b6579e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	actualHash := hashGenerator.GetHash(data)
	assert.Equal(t, expectedHash, actualHash)
}

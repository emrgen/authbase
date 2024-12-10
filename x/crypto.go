package x

// taken form gotrue

import (
	"crypto/rand"
	"encoding/base64"
	"io"
	"strings"
)

func GenerateToken() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}

	return removePadding(base64.URLEncoding.EncodeToString(b))
}

// RefreshToken creates a new random token
func RefreshToken() string {
	b := make([]byte, 32)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}
	return removePadding(base64.URLEncoding.EncodeToString(b))
}

func removePadding(token string) string {
	return strings.TrimRight(token, "=")
}

func Keygen() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}

	return removePadding(base64.URLEncoding.EncodeToString(b))
}

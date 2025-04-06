package x

// taken form gotrue

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// Keygen creates a new random key
func Keygen() string {
	return generateSecureToken(64)
}

// KeygenSize creates a new random key
func KeygenSize(size int) string {
	return generateSecureToken(size)
}

func GenerateSalt() string {
	return generateSecureToken(32)
}

func GenerateClientSecret() string {
	return generateSecureToken(40)
}

func verificationToken() string {
	return generateSecureToken(32)
}

func generateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	token := hex.EncodeToString(b)
	return token
}

func base64EncodeStripped(s string) string {
	encoded := base64.StdEncoding.EncodeToString([]byte(s))
	return strings.TrimRight(encoded, "=")
}

func base64DecodeStripped(s string) (string, error) {
	if i := len(s) % 4; i != 0 {
		s += strings.Repeat("=", 4-i)
	}
	decoded, err := base64.StdEncoding.DecodeString(s)
	return string(decoded), err
}

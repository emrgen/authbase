package x

// taken form gotrue

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"io"
	"strings"
)

func GenerateCode() string {
	return base64EncodeStripped(randomString(6))
}

func accessKey() string {
	return base64EncodeStripped(randomString(32))
}

// VerificationToken creates a new random token
func verificationToken() string {
	return base64EncodeStripped(randomString(32))
}

// Keygen creates a new random key
func Keygen() string {
	return base64EncodeStripped(randomString(20))
}

func randomString(n int) string {
	b := make([]byte, n)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(err.Error()) // rand should never fail
	}

	return string(b)
}

func generateSecureToken(length int) (string, int) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", 0
	}
	token := hex.EncodeToString(b)
	return token, len(token)
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

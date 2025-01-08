package x

import (
	"golang.org/x/crypto/argon2"
)

// HashPassword hashes the password with the salt
func HashPassword(password, salt string) []byte {
	return argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)
}

// CompareHashAndPassword compares the hash with the password and salt
func CompareHashAndPassword(password, salt, hash string) bool {
	return string(argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)) == hash
}

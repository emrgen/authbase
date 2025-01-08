package x

import (
	"fmt"
	"golang.org/x/crypto/argon2"
)

// HashPassword hashes the password with the salt
func HashPassword(password, salt string) ([]byte, error) {
	return argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32), nil
}

// CompareHashAndPassword compares the hash with the password and salt
func CompareHashAndPassword(hash, password, salt string) bool {
	return fmt.Sprintf("%s", argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 4, 32)) == hash
}

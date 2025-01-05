package x

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password with the salt
func HashPassword(password, salt string) ([]byte, error) {
	token := fmt.Sprintf("%s%s", password, salt)
	return bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
}

// CompareHashAndPassword compares the hash with the password and salt
func CompareHashAndPassword(hash, password, salt string) bool {
	token := fmt.Sprintf("%s%s", password, salt)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

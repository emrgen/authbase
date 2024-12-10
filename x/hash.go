package x

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password, secret string) ([]byte, error) {
	token := fmt.Sprintf("%s%s", password, secret)
	return bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
}

func CompareHashAndPassword(hash, password, secret string) bool {
	token := fmt.Sprintf("%s%s", password, secret)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(token))
	return err == nil
}

package x

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestCompareHashAndPassword(t *testing.T) {
	password := "password"
	salt := "salt"
	hash := HashPassword(password, salt)
	if !CompareHashAndPassword(password, salt, string(hash)) {
		t.Errorf("CompareHashAndPassword() = false; want true")
	}

}

func TestHashPassword(t *testing.T) {
	password := "bd6b1652fd6d8e6120f660a28e28f3563aeeec1e63cb0665af1735208a21f408af8e4408b"
	salt := "0172eefdb0b3cfeb55b4081e3f9bfe8cd127f2bf390c60847c9e96ab9e2157d9"
	hash := HashPassword(password, salt)
	logrus.Infof("hash: %v", string(hash))
	logrus.Infof("hash: %v", hash)
	if !CompareHashAndPassword(password, salt, string(hash)) {
		t.Errorf("CompareHashAndPassword() = false; want true")
	}
}

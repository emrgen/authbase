package x

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestResJwt function to test JWT with RSA key pair
func TestResJwt(t *testing.T) {
	private, public, err := GeneratePemKeyPair(2048)
	assert.NoError(t, err)
	assert.NotNil(t, private)
	assert.NotNil(t, public)

	priv, err := jwt.ParseRSAPrivateKeyFromPEM(private)
	assert.NoError(t, err)
	assert.NotNil(t, priv)

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, nil)
	tokenString, err := token.SignedString(priv)
	assert.NoError(t, err)
	fmt.Println("Signed Token:", tokenString)

	pub, err := jwt.ParseRSAPublicKeyFromPEM(public)
	assert.NoError(t, err)

	token, err = jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return pub, nil
	})
	assert.NoError(t, err)
}

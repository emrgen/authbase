package x

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type staticVerifier struct {
	key []byte
}

func newStaticVerifier(key []byte) *staticVerifier {
	return &staticVerifier{key: key}
}

func (v *staticVerifier) Verify(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return v.key, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to get claims")
	}

	return claims, nil
}

type staticSigner struct {
	key []byte
}

func newStaticSigner(key []byte) *staticSigner {
	return &staticSigner{key: []byte(key)}
}

func (s *staticSigner) Sign(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.key)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type StaticKeyProvider struct {
	key []byte
}

func NewStaticKeyProvider(key string) *StaticKeyProvider {
	return &StaticKeyProvider{key: []byte(key)}
}

func (r *StaticKeyProvider) GetSigner(id string) (JWTSigner, error) {
	return newStaticSigner(r.key), nil
}

func (r *StaticKeyProvider) GetVerifier(id string) (JWTVerifier, error) {
	return newStaticVerifier(r.key), nil
}

// UnverifiedKeyProvider is a key provider that does not verify the key.
type UnverifiedKeyProvider struct{}

// NewUnverifiedKeyProvider creates a new UnverifiedKeyProvider.
func NewUnverifiedKeyProvider() *UnverifiedKeyProvider {
	return &UnverifiedKeyProvider{}
}

func (r *UnverifiedKeyProvider) GetSigner(id string) (JWTSigner, error) {
	return nil, nil
}

func (r *UnverifiedKeyProvider) GetVerifier(id string) (JWTVerifier, error) {
	return nil, nil
}

type unverifiedVerifier struct{}

func newUnverifiedVerifier() *unverifiedVerifier {
	return &unverifiedVerifier{}
}

func (v *unverifiedVerifier) Verify(tokenString string) (jwt.MapClaims, error) {
	// FIXME: ParseUnverified is not safe, it should be used only for testing
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return nil, err
	}

	claims := token.Claims.(jwt.MapClaims)

	return claims, nil
}

type unverifiedSigner struct{}

func newUnverifiedSigner() *unverifiedSigner {
	return &unverifiedSigner{}
}

func (s *unverifiedSigner) Sign(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, err := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

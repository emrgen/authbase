package x

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

// ref: https://gist.github.com/goliatone/e9c13e5f046e34cef6e150d06f20a34c

func GenerateKeyPair(bitSize int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// GeneratePemKeyPair function to generate RSA private key in PEM format
func GeneratePemKeyPair(bitSize int) ([]byte, []byte, error) {
	// Generate RSA key pair
	privateKey, publicKey, err := GenerateKeyPair(bitSize)
	if err != nil {
		return nil, nil, err
	}

	// Encode private key to PEM format
	privKeyPEM := EncodePrivateKeyToPEM(privateKey)

	// Encode public key to PEM format
	pubKeyPEM := EncodePublicKeyToPEM(publicKey)

	return privKeyPEM, pubKeyPEM, nil
}

// EncodePrivateKeyToPEM function to encode RSA private key to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	}

	// Private key in PEM format
	return pem.EncodeToMemory(&privBlock)
}

// EncodePublicKeyToPEM function to encode RSA public key to PEM format
func EncodePublicKeyToPEM(publicKey *rsa.PublicKey) []byte {
	// Get ASN.1 DER format
	pubDER, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil
	}

	// pem.Block
	pubBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubDER,
	}

	// Public key in PEM format
	return pem.EncodeToMemory(&pubBlock)
}

func DecodeRSAPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// DecodeRSAPublicKey function to load RSA public key from a file
func DecodeRSAPublicKey(key []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(key)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not an RSA public key")
	}
	return rsaPubKey, nil
}

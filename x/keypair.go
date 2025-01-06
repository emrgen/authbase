package x

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
)

// ref: https://gist.github.com/goliatone/e9c13e5f046e34cef6e150d06f20a34c

// GenerateKeyPair generates a new key pair of specified byte size
// returns private key and public key in PEM format
func GenerateKeyPair(bitSize int) ([]byte, []byte, error) {
	privateKey, err := GeneratePrivateKey(bitSize)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBytes, err := GeneratePublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, nil, err
	}

	privateKeyBytes := EncodePrivateKeyToPEM(privateKey)

	return privateKeyBytes, publicKeyBytes, nil
}

// GeneratePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

// GeneratePublicKey take a rsa.PublicKey and return bytes suitable for writing to .pub file
// returns in the format "ssh-rsa ..."
func GeneratePublicKey(private *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(private)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	return pubKeyBytes, nil
}

// EncodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	pkey := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privateBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   pkey,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privateBlock)

	return privatePEM
}

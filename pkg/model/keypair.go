package model

import (
	"gorm.io/gorm"
	"time"
)

// Keypair represents a keypair used for token generation and verification
type Keypair struct {
	gorm.Model
	ClientID   string    `gorm:"uuid;primaryKey"`
	PrivateKey string    // used for token generation
	PublicKey  string    // used for token verification
	ExpiresAt  time.Time `gorm:"index"`
}

package model

import (
	"gorm.io/gorm"
	"time"
)

// Keypair represents a keypair used for token generation and verification
type Keypair struct {
	gorm.Model
	ProjectID  string    `gorm:"uuid;primaryKey"`
	PublicKey  string    // used for token verification
	PrivateKey string    // used for token generation
	ExpiresAt  time.Time `gorm:"index"`
}

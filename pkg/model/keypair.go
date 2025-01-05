package model

import "gorm.io/gorm"

type Keypair struct {
	gorm.Model
	OrganizationID string `gorm:"uuid;primaryKey"`
	PublicKey      string
	PrivateKey     string
}

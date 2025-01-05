package model

import "gorm.io/gorm"

type Keypair struct {
	gorm.Model
	ProjectID  string `gorm:"uuid;primaryKey"`
	PublicKey  string
	PrivateKey string
}

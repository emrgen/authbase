package model

import (
	"time"

	"gorm.io/gorm"
)

// AccessKey IssuedAt is same as gorm.Model.CreatedAt
type AccessKey struct {
	gorm.Model
	ID        string `gorm:"primaryKey;type:uuid"`
	Name      string
	ProjectID string `gorm:"uuid;foreignKey:ProjectID;OnDelete:CASCADE"`
	AccountID string `gorm:"uuid;foreignKey:AccountID;OnDelete:CASCADE"`
	Token     string // hashed token
	Scopes    string
	ExpireAt  time.Time
}

// TableName returns the table name of the model
func (AccessKey) TableName() string {
	return tableName("access_keys")
}

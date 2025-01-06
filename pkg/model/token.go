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
	ProjectID string `gorm:"type:uuid"`
	AccountID string `gorm:"type:uuid"`
	Token     string // hashed token
	Scopes    string
	ExpireAt  time.Time
}

// TableName returns the table name of the model
func (AccessKey) TableName() string {
	return tableName("access_keys")
}

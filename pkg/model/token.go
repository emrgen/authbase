package model

import (
	"time"

	"gorm.io/gorm"
)

// IssuedAt is same as gorm.Model.CreatedAt
type Token struct {
	gorm.Model
	ID             string `gorm:"primaryKey;type:uuid"`
	Name           string
	OrganizationID string `gorm:"type:uuid"`
	UserID         string `gorm:"type:uuid"`
	Token          string // hashed token
	Hash           string // hashed token
	ExpireAt       time.Time
}

// TableName returns the table name of the model
func (Token) TableName() string {
	return tableName("tokens")
}

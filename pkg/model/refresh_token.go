package model

import (
	"time"

	"gorm.io/gorm"
)

// RefreshToken represents a refresh token.
// It is used to generate new access tokens when the access token expires.
type RefreshToken struct {
	gorm.Model
	Token     string    `gorm:"primaryKey;not null"`
	ProjectID string    `gorm:"not null"`
	AccountID string    `gorm:"not null"`
	ExpireAt  time.Time `gorm:"not null"`
	IssuedAt  time.Time `gorm:"not null"`
}

func (RefreshToken) TableName() string {
	return tableName("refresh_tokens")
}

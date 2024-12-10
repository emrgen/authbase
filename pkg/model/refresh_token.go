package model

import (
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	gorm.Model
	Token          string    `gorm:"primaryKey;not null"`
	OrganizationID string    `gorm:"not null"`
	UserID         string    `gorm:"not null"`
	ExpireAt       time.Time `gorm:"not null"`
	IssuedAt       time.Time `gorm:"not null"`
}

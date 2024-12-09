package model

import (
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	gorm.Model
	ID       string    `gorm:"primaryKey;uuid;not null;"`
	Token    string    `gorm:"not null"` // hash of the token
	ExpireAt time.Time `gorm:"not null"`
}

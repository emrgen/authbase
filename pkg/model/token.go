package model

import (
	"gorm.io/gorm"
	"time"
)

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

func (_ Token) TableName() string {
	return tableName("tokens")
}

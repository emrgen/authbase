package model

import (
	"gorm.io/gorm"
	"time"
)

type VerificationCode struct {
	gorm.Model
	ID             string `gorm:"primaryKey;type:uuid"`
	UserID         string `gorm:"type:uuid"`
	OrganizationID string `gorm:"type:uuid"`
	Code           string
	ExpiresAt      time.Time
	Medium         string
	CallbackURL    string
}

func (_ VerificationCode) TableName() string {
	return tableName("verification_codes")
}

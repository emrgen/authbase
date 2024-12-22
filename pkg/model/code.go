package model

import (
	"time"

	"gorm.io/gorm"
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

func (VerificationCode) TableName() string {
	return tableName("verification_codes")
}

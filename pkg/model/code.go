package model

import (
	"time"

	"gorm.io/gorm"
)

type VerificationCode struct {
	gorm.Model
	ID          string `gorm:"primaryKey;type:uuid"`
	AccountID   string `gorm:"type:uuid"`
	PoolID      string `gorm:"type:uuid"`
	ProjectID   string `gorm:"type:uuid"`
	Code        string
	ExpiresAt   time.Time
	Medium      string
	CallbackURL string
}

func (VerificationCode) TableName() string {
	return tableName("verification_codes")
}

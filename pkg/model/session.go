package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID        string `gorm:"primaryKey"`
	AccountID string
	Account   *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE;"`
	ProjectID string
	ExpiredAt time.Time `gorm:"default:null"`
}

func (Session) TableName() string {
	return tableName("sessions")
}

package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID        string   `gorm:"primaryKey"`
	AccountID string   `gorm:"not null;index"`
	Account   *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE;"`
	PoolID    string   `gorm:"not null;index"`
	Pool      *Pool    `gorm:"foreignKey:PoolID;OnDelete:CASCADE;"`
	ProjectID string
	ExpiredAt time.Time `gorm:"default:null"`
}

func (Session) TableName() string {
	return tableName("sessions")
}

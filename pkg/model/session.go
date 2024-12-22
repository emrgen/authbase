package model

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	ID             string `gorm:"primaryKey"`
	UserID         string
	User           *User `gorm:"foreignKey:UserID;OnDelete:CASCADE;"`
	OrganizationID string
	ExpiredAt      time.Time `gorm:"default:null"`
}

func (Session) TableName() string {
	return tableName("sessions")
}

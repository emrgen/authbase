package model

import (
	"gorm.io/gorm"
	"time"
)

type Session struct {
	gorm.Model
	ID             string `gorm:"primaryKey"`
	UserID         string
	User           *User `gorm:"foreignKey:UserID;OnDelete:CASCADE;"`
	OrganizationID string
	ExpireAt       time.Time `gorm:"default:null"`
}

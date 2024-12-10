package model

import (
	"gorm.io/gorm"
	"time"
)

// User is the database model for users.
type User struct {
	gorm.Model
	ID             string `gorm:"primaryKey;uuid;not null;"`
	Username       string `gorm:"not null"`
	Email          string `gorm:"not null;uniqueIndex:compositeIndex;"`
	Password       string // hash of the password
	SassAdmin      bool   `gorm:"not null;default:false"`
	Member         bool   `gorm:"not null;default:false"`
	OrganizationID string `gorm:"not null;uniqueIndex:compositeIndex;"`
	Organization   *Organization
	Verified       bool `gorm:"not null;default:false"`
	VerifiedAt     *time.Time
	Disabled       bool `gorm:"not null;default:false"`
	DisabledAt     *time.Time
	Recovered      bool `gorm:"not null;default:false"`
	RecoveredAt    *time.Time
	RecoveredBy    string `gorm:"uuid;"`
}

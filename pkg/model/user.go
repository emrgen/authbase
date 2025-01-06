package model

import (
	"time"

	"gorm.io/gorm"
)

// Account is the database model for users.
// one project can have multiple users and one user can belong to multiple projects.
type Account struct {
	gorm.Model
	ID            string `gorm:"primaryKey;uuid;not null;"`
	Username      string `gorm:"not null;uniqueIndex:userProjectIndex;"`
	Email         string `gorm:"not null;uniqueIndex:compositeIndex;"`
	VisibleName   string
	Password      string // hash of the password
	Salt          string
	SassAdmin     bool   `gorm:"not null;default:false"`
	ProjectMember bool   `gorm:"not null;default:false"`
	ProjectID     string `gorm:"not null;uniqueIndex:compositeIndex;uniqueIndex:userProjectIndex;"`
	Project       *Project
	Verified      bool      `gorm:"not null;default:false"`
	VerifiedAt    time.Time `gorm:"default:null"`
	Disabled      bool      `gorm:"not null;default:false"`
	DisabledAt    time.Time `gorm:"default:null"`
	Recovered     bool      `gorm:"not null;default:false"`
	RecoveredAt   time.Time `gorm:"default:null"`
	RecoveredBy   string    `gorm:"uuid;"`
}

func (Account) TableName() string {
	return tableName("accounts")
}

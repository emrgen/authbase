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
	Username      string `gorm:"not null;uniqueIndex:compositeIndex;"`
	Email         string `gorm:"not null;uniqueIndex:compositeIndex;"`
	VisibleName   string
	PasswordHash  string // hash of the password
	Salt          string
	SassAdmin     bool      `gorm:"not null;default:false"`
	ProjectMember bool      `gorm:"not null;default:false"`
	PoolID        string    `gorm:"not null;uniqueIndex:compositeIndex;"`
	ProjectID     string    `gorm:"uuid"`
	Project       *Project  `gorm:"foreignKey:ProjectID;constraint:OnDelete:NO ACTION"`
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

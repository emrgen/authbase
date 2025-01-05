package model

import (
	"time"

	"gorm.io/gorm"
)

// User is the database model for users.
// one project can have multiple users and one user can belong to multiple projects.
type User struct {
	gorm.Model
	ID          string `gorm:"primaryKey;uuid;not null;"`
	Username    string `gorm:"not null;uniqueIndex:userProjectIndex;"`
	Email       string `gorm:"not null;uniqueIndex:compositeIndex;"`
	Password    string // hash of the password
	Salt        string
	SassAdmin   bool   `gorm:"not null;default:false"`
	Member      bool   `gorm:"not null;default:false"`
	ProjectID   string `gorm:"not null;uniqueIndex:compositeIndex;uniqueIndex:userProjectIndex;"`
	Project     *Project
	Verified    bool      `gorm:"not null;default:false"`
	VerifiedAt  time.Time `gorm:"default:null"`
	Disabled    bool      `gorm:"not null;default:false"`
	DisabledAt  time.Time `gorm:"default:null"`
	Recovered   bool      `gorm:"not null;default:false"`
	RecoveredAt time.Time `gorm:"default:null"`
	RecoveredBy string    `gorm:"uuid;"`
}

func (User) TableName() string {
	return tableName("users")
}

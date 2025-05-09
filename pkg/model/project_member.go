package model

import (
	"gorm.io/gorm"
	"time"
)

type ProjectMember struct {
	ProjectID  string   `gorm:"primaryKey;not null"`
	AccountID  string   `gorm:"primaryKey;not null"`
	Permission uint32   `gorm:"not null;default:0"`
	Project    *Project `gorm:"foreignKey:ProjectID;OnDelete:CASCADE"` // delete permissions when project is deleted
	Account    *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE"` // delete permissions when user is deleted

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (ProjectMember) TableName() string {
	return tableName("project_members")
}

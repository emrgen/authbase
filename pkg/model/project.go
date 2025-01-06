package model

import "gorm.io/gorm"

// PasswordPolicy represents a password policy for an project
type PasswordPolicy struct {
	MinLength int `json:"min_length"`
	MaxLength int `json:"max_length"`
	MinUpper  int `json:"min_upper"`
	MinLower  int `json:"min_lower"`
	MinDigit  int `json:"min_digit"`
	MinSymbol int `json:"min_symbol"`
}

// Project represents an project
type Project struct {
	gorm.Model
	ID                string         `gorm:"primaryKey;uuid;not null;"`
	ParentProjectID   string         `gorm:"uuid;default:null"`
	Name              string         `gorm:"not null;unique;index:idx_project_name"` // unique project name
	DisplayName       string         `gorm:"not null"`
	OwnerID           string         `gorm:"not null"`
	Owner             *Account       `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Master            bool           `gorm:"not null;default:false"`
	AllowedDomains    string         `gorm:"not null;default:''"`
	EmailVerification bool           `gorm:"not null;default:false"`
	PasswordPolicy    PasswordPolicy `gorm:"embedded;embeddedPrefix:password_policy_"`
}

// TableName returns the table name for the project model
func (Project) TableName() string {
	return tableName("projects")
}

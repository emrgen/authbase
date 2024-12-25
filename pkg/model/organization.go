package model

import "gorm.io/gorm"

// PasswordPolicy represents a password policy for an organization
type PasswordPolicy struct {
	MinLength int `json:"min_length"`
	MaxLength int `json:"max_length"`
	MinUpper  int `json:"min_upper"`
	MinLower  int `json:"min_lower"`
	MinDigit  int `json:"min_digit"`
	MinSymbol int `json:"min_symbol"`
}

// Organization represents an organization
type Organization struct {
	gorm.Model
	ID                string         `gorm:"primaryKey;uuid;not null;"`
	Name              string         `gorm:"not null;unique;index:idx_organization_name"`
	OwnerID           string         `gorm:"not null"`
	ProjectID         string         `gorm:"uuid;default:null"`                              // filled when running in multistore mode
	Owner             *User          `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"` // filled when running in multistore mode
	Master            bool           `gorm:"not null;default:false"`
	AllowedDomains    string         `gorm:"not null;default:''"`
	EmailVerification bool           `gorm:"not null;default:false"`
	PasswordPolicy    PasswordPolicy `gorm:"embedded;embeddedPrefix:password_policy_"`
}

// TableName returns the table name for the organization model
func (Organization) TableName() string {
	return tableName("organizations")
}

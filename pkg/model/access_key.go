package model

import (
	"time"

	"gorm.io/gorm"
)

// AccessKey IssuedAt is same as gorm.Model.CreatedAt
type AccessKey struct {
	gorm.Model
	ID          string `gorm:"primaryKey;type:uuid"`
	Name        string
	ProjectID   string   `gorm:"uuid"`
	Project     *Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	AccountID   string   `gorm:"uuid"`
	Account     *Account `gorm:"foreignKey:AccountID;constraint:OnDelete:CASCADE"`
	PoolID      string   `gorm:"uuid"`
	Pool        *Pool    `gorm:"foreignKey:PoolID;constraint:OnDelete:CASCADE"`
	Token       string   // hashed token
	Scopes      string   `gorm:"default:[]"`
	Permissions string   `gorm:"default:[]"`
	ExpireAt    time.Time
}

// TableName returns the table name of the model
func (AccessKey) TableName() string {
	return tableName("access_keys")
}

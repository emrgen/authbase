package model

import (
	"gorm.io/gorm"
	"time"
)

type PoolMember struct {
	AccountID  string `gorm:"primaryKey;not null"`
	PoolID     string `gorm:"primaryKey;not null"`
	Permission uint32 `gorm:"not null;default:0"`

	Account *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE"`
	Pool    *Pool    `gorm:"foreignKey:PoolID;OnDelete:CASCADE"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (PoolMember) TableName() string {
	return tableName("pool_members")
}

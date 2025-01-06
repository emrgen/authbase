package model

import "gorm.io/gorm"

type PoolMember struct {
	gorm.Model
	AccountID  string `gorm:"primaryKey;not null"`
	PoolID     string `gorm:"primaryKey;not null"`
	Permission uint32 `gorm:"not null;default:0"`

	Account *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE"`
	Pool    *Pool    `gorm:"foreignKey:PoolID;OnDelete:CASCADE"`
}

func (PoolMember) TableName() string {
	return tableName("pool_members")
}

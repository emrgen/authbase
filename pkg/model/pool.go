package model

import "gorm.io/gorm"

// Pool is the database model for pools.
// A pool is a group of accounts with a specific permission level.
type Pool struct {
	gorm.Model
	ID        string `gorm:"primaryKey;uuid;not null;"`                  // ID of the pool
	Name      string `gorm:"not null;index:idx_name_project_id,unique;"` // Name of the pool
	ProjectID string `gorm:"not null;index:idx_name_project_id,unique;"` // Project ID
	Default   bool   `gorm:"not null;default:false;"`                    // Default pool
}

func (Pool) TableName() string {
	return tableName("pools")
}

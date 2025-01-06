package model

import "gorm.io/gorm"

// Pool is the database model for pools.
// A pool is a group of accounts with a specific permission level.
type Pool struct {
	gorm.Model
	ID        string `gorm:"primaryKey;uuid;not null;"`                   // ID of the pool
	Name      string `gorm:"not null;uniqueIndex:pool_project_id__name;"` // Name of the pool
	ProjectID string `gorm:"not null;uniqueIndex:ool_project_id__name;"`  // Project ID
	Default   bool   `gorm:"not null;default:false;"`                     // Default pool
}

func (Pool) TableName() string {
	return tableName("pools")
}

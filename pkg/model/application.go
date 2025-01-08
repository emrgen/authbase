package model

import "gorm.io/gorm"

type Application struct {
	gorm.Model
	ID     string `gorm:"uuid;primaryKey"`
	Name   string `gorm:"not null;uniqueIndex:idx_app_name_pool_id"`
	PoolID string `gorm:"not null;uniqueIndex:idx_app_name_pool_id"`
}

func (a *Application) TableName() string {
	return tableName("applications")
}

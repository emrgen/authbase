package model

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ID     string `gorm:"uuid;primaryKey"`
	PoolID string `gorm:"uuid"`
	Name   string
	Secret string
}

func (c *Client) TableName() string {
	return tableName("clients")
}

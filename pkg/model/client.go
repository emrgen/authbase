package model

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ID        string `gorm:"uuid;primaryKey"`
	ProjectID string `gorm:"uuid;index"`
	Name      string
	Secret    string
}

func (c *Client) TableName() string {
	return tableName("clients")
}

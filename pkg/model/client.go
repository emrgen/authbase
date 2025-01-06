package model

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ID               string `gorm:"uuid;primaryKey"`
	PoolID           string `gorm:"uuid;index:idx_client_pool_id_name,unique"`
	Pool             *Pool  `gorm:"foreignKey:PoolID;constraint:OnDelete:CASCADE"`
	Name             string `gorm:"index:idx_client_pool_id_name,unique"`
	Secret           string
	CreatedByID      string   `gorm:"uuid"`
	CreatedByAccount *Account `gorm:"foreignKey:CreatedByID"`
}

func (c *Client) TableName() string {
	return tableName("clients")
}

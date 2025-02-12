package model

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	ID         string `gorm:"uuid;primaryKey"`
	PoolID     string `gorm:"uuid;index:idx_client_pool_id_name,unique"`
	Pool       *Pool  `gorm:"foreignKey:PoolID;constraint:OnDelete:CASCADE"`
	Name       string `gorm:"index:idx_client_pool_id_name,unique"`
	SecretHash string
	// TODO: This is dangerous, we should not store the secret in the database.
	// We should only store the secret a secret manager.
	// This is needed if users needs to check the client secret for integration.
	Secret           string
	Salt             string
	CreatedByID      string   `gorm:"uuid"`
	CreatedByAccount *Account `gorm:"foreignKey:CreatedByID"`
	Default          bool     `gorm:"default:false"`
}

func (c *Client) TableName() string {
	return tableName("clients")
}

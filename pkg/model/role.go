package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	ID         string      `gorm:"primaryKey;uuid;not null;"`              // ID of the role
	Name       string      `gorm:"uniqueIndex:idx_name_pool_id;not null;"` // Name of the role
	PoolID     string      `gorm:"uniqueIndex:idx_name_pool_id;not null;"` // Pool ID
	Pool       *Pool       `gorm:"foreignKey:PoolID;references:ID;OnDelete:CASCADE"`
	Groups     []*Group    `gorm:"many2many:group_roles;constraint:OnDelete:CASCADE"`
	Attributes interface{} `gorm:"type:jsonb;"` // Attributes of the role
}

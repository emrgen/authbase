package model

import "gorm.io/gorm"

const (
	// GroupScopePrefix is the prefix for scopes that are given to account from a group.
	groupScopePrefix = "x:"
)

// Group represents a group of accounts with a common purpose (scopes).
type Group struct {
	gorm.Model
	ID     string  `gorm:"uuid;not null"`
	Name   string  `gorm:"uuid;not null;uniqueIndex:idx_pool_id_group_name"`
	PoolID string  `gorm:"uuid;not null;uniqueIndex:idx_pool_id_group_name"`
	Pool   *Pool   `gorm:"foreignKey:PoolID;OnDelete:CASCADE"`
	Roles  []*Role `gorm:"many2many:group_roles;constraint:OnDelete:CASCADE"`
}

// GroupMember represents a member of a group.
type GroupMember struct {
	GroupID   string `gorm:"uuid;not null;primaryKey"`
	AccountID string `gorm:"uuid;not null;primaryKey"`

	Group   *Group   `gorm:"foreignKey:GroupID;OnDelete:CASCADE"`
	Account *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE"`
}

package model

import "gorm.io/gorm"

const (
	// GroupScopePrefix is the prefix for scopes that are given to account from a group.
	groupScopePrefix = "x:"
)

// Group represents a group of accounts with a common purpose (scopes).
type Group struct {
	gorm.Model
	ID     string `gorm:"uuid;not null"`
	Name   string `gorm:"uuid;not null;uniqueIndex:idx_pool_id_group_name"`
	PoolID string `gorm:"uuid;not null;uniqueIndex:idx_pool_id_group_name"`
	Pool   *Pool  `gorm:"foreignKey:PoolID;OnDelete:CASCADE"`
	Scopes string `gorm:"not null"`
}

// GroupMember represents a member of a group.
type GroupMember struct {
	GroupID   string `gorm:"uuid;not null"`
	AccountID string `gorm:"uuid;not null"`

	Group   *Group   `gorm:"foreignKey:GroupID;OnDelete:CASCADE"`
	Account *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE"`
}

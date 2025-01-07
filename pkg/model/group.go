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
	Name   string `gorm:"unique;not null"`
	PoolID string `json:"unique;not null"`
	Pool   *Pool  `gorm:"foreignKey:PoolID;OnDelete:CASCADE"`
	Scopes string `json:"scopes"`
}

// GroupMember represents a member of a group.
type GroupMember struct {
	GroupID   string `gorm:"uuid;not null"`
	AccountID string `gorm:"uuid;not null"`

	Group   *Group   `gorm:"foreignKey:GroupID;OnDelete:CASCADE"`
	Account *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE"`
}

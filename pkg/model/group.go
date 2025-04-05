package model

import "gorm.io/gorm"

// Group represents a group of accounts with a common purpose (scopes).
// Groups have associated roles that define the permissions for the group.
// Accounts with group membership inherit the roles of the group.
type Group struct {
	gorm.Model
	ID       string  `gorm:"uuid;not null"`
	Name     string  `gorm:"uuid;not null;uniqueIndex:idx_pool_id_group_name"`
	PoolID   string  `gorm:"uuid;not null;uniqueIndex:idx_pool_id_group_name"`
	Pool     *Pool   `gorm:"foreignKey:PoolID;OnDelete:CASCADE"`
	Roles    []*Role `gorm:"many2many:group_roles;constraint:OnDelete:CASCADE"`
	Internal bool    `gorm:"not null;default:false"` // Internal groups are not allowed to be deleted
}

// GroupMember represents a member of a group.
type GroupMember struct {
	GroupID   string `gorm:"uuid;not null;primaryKey"`
	AccountID string `gorm:"uuid;not null;primaryKey"`

	Group   *Group   `gorm:"foreignKey:GroupID;OnDelete:CASCADE"`
	Account *Account `gorm:"foreignKey:AccountID;OnDelete:CASCADE"`
}

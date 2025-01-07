package model

type Role struct {
	Name       string      `gorm:"primaryKey;not null;"`          // Name of the role
	PoolID     string      `gorm:"primaryKey;not null;not null;"` // Pool ID
	Pool       *Pool       `gorm:"foreignKey:PoolID;references:ID;OnDelete:CASCADE"`
	Groups     []*Group    `gorm:"many2many:group_roles"`
	Attributes interface{} `gorm:"type:jsonb;"` // Attributes of the role
}

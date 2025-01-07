package model

type Role struct {
	Name       string      `gorm:"primaryKey"`
	PoolID     string      `gorm:"primaryKey"`
	Attributes interface{} `gorm:"type:jsonb"`
	Groups     []*Group    `gorm:"many2many:group_roles;foreignKey:Name;joinForeignKey:RoleName;references:Name;joinReferences:GroupID;OnDelete:CASCADE"`
}

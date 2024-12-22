package model

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	ID        string `gorm:"primaryKey;uuid;not null;"`
	Name      string `gorm:"not null;unique;index:idx_organization_name"`
	OwnerID   string `gorm:"not null"`
	ProjectID string `gorm:"uuid;default:null"`                              // filled when running in multistore mode
	Owner     *User  `gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"` // filled when running in multistore mode
	Config    string `gorm:"type:json"`
	Master    bool   `gorm:"not null;default:false"`
}

func (Organization) TableName() string {
	return tableName("organizations")
}

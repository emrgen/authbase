package model

import "gorm.io/gorm"

type Organization struct {
	gorm.Model
	ID      string `gorm:"primaryKey;uuid;not null;"`
	Name    string `gorm:"not null;unique"`
	OwnerID string `gorm:"not null"`
	Owner   *User  `gorm:"foreignKey:OwnerID"`
	Config  string `gorm:"type:json"`
	Master  bool   `gorm:"not null;default:false"`
}

package model

import "gorm.io/gorm"

type Secret struct {
	gorm.Model
	ID    string `gorm:"primaryKey"`
	Value string
}

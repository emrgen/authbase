package model

import "gorm.io/gorm"

// User is the database model for users.
type User struct {
	gorm.Model
	ID    string `gorm:"primaryKey;uuid;not null;"`
	Name  string `gorm:"not null"`
	Email string `gorm:"unique;not null"`
}

package model

import (
	"github.com/emrgen/authbase/pkg/config"
	"gorm.io/gorm"
)

// Migrate creates the tables in the database
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&RefreshToken{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Permission{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&OauthProvider{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Organization{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Token{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&VerificationCode{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Session{}); err != nil {
		return err
	}

	return nil
}

// tableName returns the table name for the given model depending on the database type
func tableName(name string) string {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	if cfg.DB.Type == "sqlite3" {
		return name
	} else {
		return "authbase." + name
	}
}

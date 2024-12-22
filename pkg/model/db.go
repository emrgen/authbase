package model

import (
	"github.com/emrgen/authbase/pkg/config"
	"gorm.io/gorm"
)

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

func tableName(name string) string {
	cfg := config.GetConfig()
	if cfg.DB.Type == "sqlite3" {
		return name
	} else {
		return "authbase." + name
	}
}

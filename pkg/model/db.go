package model

import (
	"github.com/emrgen/authbase/pkg/config"
	"gorm.io/gorm"
)

// Migrate creates the tables in the database
func Migrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&Account{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&RefreshToken{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&ProjectMember{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&OauthProvider{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Project{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&AccessKey{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&VerificationCode{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Session{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Keypair{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Client{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Pool{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&PoolMember{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Group{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&GroupMemberAccount{}); err != nil {
		return err
	}
	
	if err := db.AutoMigrate(&GroupMemberAccessKey{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Role{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Application{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&Secret{}); err != nil {
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

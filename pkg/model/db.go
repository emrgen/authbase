package model

import "gorm.io/gorm"

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

	if err := db.AutoMigrate(&Provider{}); err != nil {
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

	return nil
}

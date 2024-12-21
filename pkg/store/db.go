package store

import (
	"github.com/emrgen/authbase/pkg/config"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB() AuthBaseStore {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(err)
	}

	switch cfg.DB.Type {
	case "sqlite3":
		logrus.Infof("connecting to sqlite database: %s", cfg.DB.FilePath)
		db, err := gorm.Open(sqlite.Open(cfg.DB.FilePath), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		return NewGormStore(db)
	}

	panic("unknown database type")
}

package config

import (
	"github.com/emrgen/authbase/pkg/store"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetDB() store.AuthBaseStore {
	cfg, err := FromEnv()
	if err != nil {
		panic(err)
	}

	switch cfg.DB.Type {
	case "sqlite":
		logrus.Infof("connecting to sqlite database: %s", cfg.DB.FilePath)
		db, err := gorm.Open(sqlite.Open(cfg.DB.FilePath), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		return store.NewGormStore(db)
	}

	panic("unknown database type")
}

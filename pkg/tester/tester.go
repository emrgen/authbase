package tester

import (
	"github.com/emrgen/authbase/pkg/cache"
	"github.com/emrgen/authbase/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

const (
	testPath = "../.test/"
)

var (
	db *gorm.DB
)

func Setup() {
	_ = os.Setenv("ENV", "test")

	err := os.MkdirAll(testPath+"/db", os.ModePerm)
	if err != nil {
		panic(err)
	}

	db, err = gorm.Open(sqlite.Open(testPath+"db/authbase.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	err = model.Migrate(db)
	if err != nil {
		panic(err)
	}
}

func TestDB() *gorm.DB {
	return db
}

func RemoveDBFile() {
	err := os.RemoveAll(testPath)
	if err != nil {
		panic(err)
	}
}

func CleanUp() {
	db.Exec("DELETE FROM accounts")
	db.Exec("DELETE FROM projects")
	db.Exec("DELETE FROM project_members")
	db.Exec("DELETE FROM pools")
	db.Exec("DELETE FROM pool_members")
	db.Exec("DELETE FROM clients")
	db.Exec("DELETE FROM access_keys")
	db.Exec("DELETE FROM refresh_tokens")
	db.Exec("DELETE FROM sessions")
}

func TestRedis() *cache.Redis {
	return cache.NewRedisClient()
}

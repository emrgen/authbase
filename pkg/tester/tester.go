package tester

import (
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
	DropUsers()
}

func DropUsers() {
	db.Exec("DELETE FROM users")
}

package storage

import (
	"os"

	"github.com/7onn/croupier/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func Connect() {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "db/.sqlite"
	}
	var err error
	db, err = gorm.Open(sqlite.Open(dbURL), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to database")
	}

	db.AutoMigrate(
		&models.User{},
		&models.Table{},
	)
}

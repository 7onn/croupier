package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
)

var db *gorm.DB

func init() {

	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	fmt.Println(os.Getenv("DB_URI"))

	conn, err := gorm.Open("postgres", os.Getenv("DB_URI"))
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{})
}

//GetDB !
func GetDB() *gorm.DB {
	return db
}

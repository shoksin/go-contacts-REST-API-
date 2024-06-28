package models

import (
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *gorm.DB

func init() {
	err := godotenv.Load() //загружает переменные окружения из файла .env в среду выполнения приложения
	if err != nil {
		fmt.Print(err)
	}

	username := os.Getenv("db_user")
	password := os.Getenv("db_password")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, username, dbName, password)

	conn, err := gorm.Open("postgres", dbUri)
	if err != nil {
		fmt.Print(err)
	}

	db = conn
	db.Debug().AutoMigrate(&Account{}, &Contact{}) //Миграция базы данных

}

func GetDB() *gorm.DB {
	return db
}

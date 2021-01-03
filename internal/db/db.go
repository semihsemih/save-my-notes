package db

import (
	"fmt"
	"github.com/semihsemih/save-my-notes/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var db *gorm.DB

func ConnectDB() *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	err = db.AutoMigrate(&models.User{}, &models.List{}, &models.Note{})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
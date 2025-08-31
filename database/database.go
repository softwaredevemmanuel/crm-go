package database


import (
	"fmt"
	"log"
	// "os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"crm-go/models"
)

var DB *gorm.DB


func Connect() {
	dsn := "host=localhost user=postgres password=root dbname=go_crm port=5432 sslmode=disable"
	
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Run migrations
	db.AutoMigrate(&models.User{})

	DB = db
	fmt.Println("Database connected!")
}

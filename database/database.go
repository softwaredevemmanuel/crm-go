package database

import (
	"fmt"
	"log"

	"crm-go/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"crm-go/models"
)

var cfg = config.LoadEnv()


func MigrateDatabase() {
	// Connect to database
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	   
	// Run migrations to database
    db.AutoMigrate(&models.User{})
    db.AutoMigrate(&models.PasswordReset{})
	db.AutoMigrate(&models.Course{})

	log.Println("✅ Database migrated successfully")

}

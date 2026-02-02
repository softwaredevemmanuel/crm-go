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
	db.AutoMigrate(&models.Product{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.CourseProductTable{})
	db.AutoMigrate(&models.CourseCategoryTable{})
	db.AutoMigrate(&models.UserSession{})
	db.AutoMigrate(&models.Enrollment{})
	db.AutoMigrate(&models.ActivityLog{})
	db.AutoMigrate(&models.Announcement{})
	db.AutoMigrate(&models.Assignment{})
	db.AutoMigrate(&models.AssignmentSubmission{})
	db.AutoMigrate(&models.Chapter{})
	db.AutoMigrate(&models.Lessons{})
	db.AutoMigrate(&models.CourseMaterial{})
	db.AutoMigrate(&models.DeletedRecord{})
	db.AutoMigrate(&models.Topic{})
	db.AutoMigrate(&models.Grade{})

	log.Println("✅ Database migrated successfully")

}

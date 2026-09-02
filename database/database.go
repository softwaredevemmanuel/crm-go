package database

import (
	"fmt"
	"log"

	"crm-go/config"
	"crm-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
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
	db.AutoMigrate(&models.ActivityLog{})
	db.AutoMigrate(&models.Announcement{})
	db.AutoMigrate(&models.Assignment{})
	db.AutoMigrate(&models.AssignmentSubmission{})
	db.AutoMigrate(&models.CourseMaterial{})
	db.AutoMigrate(&models.DeletedRecord{})
	db.AutoMigrate(&models.Lesson{})
	db.AutoMigrate(&models.Grade{})
	db.AutoMigrate(&models.ObjectiveQuestion{})
	db.AutoMigrate(&models.QuestionOption{})
	db.AutoMigrate(&models.Subject{})
	db.AutoMigrate(&models.ClassGrade{})
	db.AutoMigrate(&models.Department{})
	db.AutoMigrate(&models.Arm{})
	db.AutoMigrate(&models.Guardian{})
	db.AutoMigrate(&models.Address{})
	db.AutoMigrate(&models.AcademicSession{})
	db.AutoMigrate(&models.GradeSubject{})
	db.AutoMigrate(&models.StudentEnrollment{})
	db.AutoMigrate(&models.Notification{})
	db.AutoMigrate(&models.PushSubscription{})
	db.AutoMigrate(&models.TeacherSubjectAssignment{})
	db.AutoMigrate(&models.SchemeOfWork{})
	db.AutoMigrate(&models.ArmClassTeacher{})
	db.AutoMigrate(&models.Term{})
	db.AutoMigrate(&models.SchemeOfWork{})
	db.AutoMigrate(&models.Module{})
	db.AutoMigrate(&models.Topic{})
	db.AutoMigrate(&models.Lesson{})
	


	log.Println("✅ Database migrated successfully")

}

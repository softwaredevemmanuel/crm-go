package models

import (
	"time"
	"github.com/google/uuid"


)

type StudentProfile struct {
    ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
    UserID            uuid.UUID      `gorm:"type:uuid;unique;not null" json:"user_id"` // Foreign key to users table
    
    // Personal Information
    FirstName         string         `gorm:"type:varchar(100);not null" json:"first_name"`
    LastName          string         `gorm:"type:varchar(100);not null" json:"last_name"`
    ProfilePicture    string         `gorm:"type:varchar(500)" json:"profile_picture,omitempty"`
    Bio               string         `gorm:"type:text" json:"bio,omitempty"` // e.g., "Aspiring Web Developer"
    DateOfBirth       *time.Time     `json:"date_of_birth,omitempty"`
    Gender            string         `gorm:"type:varchar(20);check:gender IN ('male', 'female', 'other', 'prefer_not_to_say')" json:"gender,omitempty"`
    
    // Contact Information
    PhoneNumber       string         `gorm:"type:varchar(50)" json:"phone_number,omitempty"`
    
    // Location
    Country           string         `gorm:"type:varchar(100)" json:"country,omitempty"`
    City              string         `gorm:"type:varchar(100)" json:"city,omitempty"`
    State             string         `gorm:"type:varchar(100)" json:"state,omitempty"`
    Timezone          string         `gorm:"type:varchar(50)" json:"timezone,omitempty"`
    Language          string         `gorm:"type:varchar(10);default:'en'" json:"language"`
    
    
    // Academic Progress
    TotalCoursesEnrolled    int      `gorm:"default:0" json:"total_courses_enrolled"`
    TotalCoursesCompleted   int      `gorm:"default:0" json:"total_courses_completed"`
    TotalLearningHours      int      `gorm:"default:0" json:"total_learning_hours"` // in hours
    CurrentStreak           int      `gorm:"default:0" json:"current_streak"` // consecutive days
    LongestStreak           int      `gorm:"default:0" json:"longest_streak"`
    LastActiveDate          *time.Time `json:"last_active_date,omitempty"`
    
    // Achievements
    CertificatesCount       int      `gorm:"default:0" json:"certificates_count"`
    BadgesCount             int      `gorm:"default:0" json:"badges_count"`
    CompletionRate          float64  `gorm:"type:decimal(5,2);default:0" json:"completion_rate"` // percentage
    AverageGrade            float64  `gorm:"type:decimal(5,2);default:0" json:"average_grade"`
    
    // Preferences & Settings
    EmailNotifications      bool     `gorm:"default:true" json:"email_notifications"`
    SMSNotifications        bool     `gorm:"default:false" json:"sms_notifications"`
    PushNotifications       bool     `gorm:"default:true" json:"push_notifications"`
    NewsletterSubscribed    bool     `gorm:"default:true" json:"newsletter_subscribed"`
    
    // Account Status
    ProfileCompletion       int      `gorm:"default:0;check:profile_completion >= 0 AND profile_completion <= 100" json:"profile_completion"` // percentage
    IsVerified              bool     `gorm:"default:false" json:"is_verified"`
    AccountStatus           string   `gorm:"type:varchar(20);default:'active';check:account_status IN ('active', 'inactive', 'suspended', 'deleted')" json:"account_status"`
    
    // Metadata
    CreatedAt               time.Time `json:"created_at"`
    UpdatedAt               time.Time `json:"updated_at"`
  
    
    // Relationships
    User                    User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Enrollments             []Enrollment `gorm:"foreignKey:StudentID" json:"enrollments,omitempty"`
    Certificates            []Certificate `gorm:"foreignKey:StudentID" json:"certificates,omitempty"`
}
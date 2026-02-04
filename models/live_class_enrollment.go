// models/live_class_enrollment.go
package models

import (
    "time"
    "github.com/google/uuid"
)

type LiveClassEnrollment struct {
    ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Relationships
    LiveClassID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_enrollment_live_class"`
    StudentID    uuid.UUID  `gorm:"type:uuid;not null;index:idx_enrollment_student"`
    CourseID     uuid.UUID  `gorm:"type:uuid;not null;index"` // Denormalized for faster queries
    TutorID      uuid.UUID  `gorm:"type:uuid;not null;index"` // Denormalized
    
    // Enrollment Details
    Status       string     `gorm:"type:varchar(20);default:'pending';not null;check:status IN ('pending', 'confirmed', 'waitlisted', 'cancelled', 'attended', 'absent')"`
    JoinedAt     *time.Time `gorm:"index"` // When they actually joined the class
    LeftAt       *time.Time // When they left the class
    
    // Payment & Access
    PaidAmount   float64    `gorm:"default:0"` // If class has separate fee
    PaymentID    *uuid.UUID `gorm:"type:uuid"` // Link to payment if paid separately
    DiscountCode string     `gorm:"type:varchar(50)"`
    
    // Invitation & Approval
    InvitedBy    *uuid.UUID `gorm:"type:uuid"` // Who invited this student
    ApprovedBy   *uuid.UUID `gorm:"type:uuid"` // Who approved if requires approval
    InvitedAt    *time.Time
    
    // Access Details
    AccessToken  string     `gorm:"type:varchar(100);index"` // Unique token for joining
    MeetingURL   string     `gorm:"type:varchar(500)"` // Personalized join URL if any
    
    // Attendance & Participation
    Duration     int        `gorm:"default:0"` // Minutes attended
    PollAnswers  []PollAnswer `gorm:"foreignKey:EnrollmentID"` // For polls during class
    Questions    []Question   `gorm:"foreignKey:EnrollmentID"` // Questions asked
    
    // Feedback
    Rating       int        `gorm:"check:rating >= 0 AND rating <= 5"` // 1-5 stars
    Feedback     string     `gorm:"type:text"`
    
    // Timestamps
    EnrolledAt   time.Time  `gorm:"not null;index"`
    UpdatedAt    time.Time
    
    // Relationships
    LiveClass    LiveClass  `gorm:"foreignKey:LiveClassID"`
    Student      User       `gorm:"foreignKey:StudentID"`
    Tutor        User       `gorm:"foreignKey:TutorID"`
    Course       Course     `gorm:"foreignKey:CourseID"`
}

// TableName specifies the table name
func (LiveClassEnrollment) TableName() string {
    return "live_class_enrollments"
}

// PollAnswer model for class polls
type PollAnswer struct {
    ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    EnrollmentID uuid.UUID  `gorm:"type:uuid;not null;index"`
    PollID       uuid.UUID  `gorm:"type:uuid;not null;index"`
    Answer       string     `gorm:"type:text;not null"`
    AnsweredAt   time.Time
}

// Question model for student questions during class
type Questions struct {
    ID           uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    EnrollmentID uuid.UUID  `gorm:"type:uuid;not null;index"`
    Question     string     `gorm:"type:text;not null"`
    Answer       string     `gorm:"type:text"`
    AskedAt      time.Time
    AnsweredAt   *time.Time
    IsAnonymous  bool       `gorm:"default:false"`
}
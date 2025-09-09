package models

import (
	"time"
	"github.com/google/uuid"


)

type AssignmentSubmission struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Core Relationships
    AssignmentID       uuid.UUID      `gorm:"type:uuid;not null;index"`
    StudentID          uuid.UUID      `gorm:"type:uuid;not null;index"`
    CourseID           uuid.UUID      `gorm:"type:uuid;not null;index"`
    EnrollmentID       uuid.UUID      `gorm:"type:uuid;not null;index"`
    
    // Submission Content
    Title              string         `gorm:"type:varchar(255);not null"`              // Submission title
    Content            string         `gorm:"type:text"`                               // Text content/description
    SubmissionType     string         `gorm:"type:varchar(50);default:'text';check:submission_type IN ('text', 'file', 'url', 'code', 'video', 'audio', 'image', 'document', 'presentation', 'multiple')"`
    
    // Submission Status
    Status             string         `gorm:"type:varchar(20);default:'submitted';check:status IN ('draft', 'submitted', 'under_review', 'graded', 'rejected', 'resubmitted', 'late', 'excused')"`
    SubmissionNumber   int            `gorm:"default:1"`                               // Attempt number
    IsFinal            bool           `gorm:"default:false"`                           // Final submission
    IsPlagiarized      bool           `gorm:"default:false"`                           // Plagiarism detected
    PlagiarismScore    float64        `gorm:"type:decimal(5,2);default:0"`             // Plagiarism percentage
    

 
    // Timing & Deadlines
    SubmittedAt        time.Time      `gorm:"default:CURRENT_TIMESTAMP;index"`         // When submitted
    DueDate            time.Time      `gorm:"not null;index"`                          // Assignment due date
    TimeExtension      int            `gorm:"default:0"`                               // Extension in hours
    IsLate             bool           `gorm:"default:false"`                           // Submitted after deadline
    LatePenalty        float64        `gorm:"type:decimal(5,2);default:0"`             // Penalty percentage
    DaysLate           int            `gorm:"default:0"`                               // Number of days late
    TimeSpent          int            `gorm:"default:0"`                               // Time spent in minutes
    
                                     // AI-suggested improvements
    
    // Relationships
    Assignment         Assignment     `gorm:"foreignKey:AssignmentID"`
    Student            User           `gorm:"foreignKey:StudentID"`
    Course             Course         `gorm:"foreignKey:CourseID"`
    Enrollment         Enrollment     `gorm:"foreignKey:EnrollmentID"`
    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
    DraftSavedAt       *time.Time     `gorm:"index"`                                   // Last draft save
    ResubmittedAt      *time.Time     `gorm:"index"`                                   // When resubmitted
}

// TableName specifies the table name
func (AssignmentSubmission) TableName() string {
    return "assignment_submissions"
}
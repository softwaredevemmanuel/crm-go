package models

import (
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
	"encoding/json"
)


type AssignmentSubmission struct {
    ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

    // Relations
    AssignmentID uuid.UUID `gorm:"type:uuid;not null;index"`
    StudentID    uuid.UUID `gorm:"type:uuid;not null;index"`

    // Submission content
    SubmissionType string `gorm:"type:varchar(50);not null;check:submission_type IN ('text', 'file', 'url', 'code', 'video', 'audio', 'image', 'document', 'presentation', 'multiple')"`

    TextContent string `gorm:"type:text"`        // for text answers
    FileURL     string `gorm:"type:text"`        // file storage link
    ExternalURL string `gorm:"type:text"`        // github, google doc, etc
    CodeRepoURL string `gorm:"type:text"`        // optional for code
    Metadata    datatypes.JSON `gorm:"type:jsonb"` // flexible (screenshots, multiple files, etc)

    // Status & grading
    Status string `gorm:"type:varchar(20);default:'submitted';check:status IN ('draft','submitted','late','under_review','graded','rejected')"`

    // Timestamps
    SubmittedAt time.Time `gorm:"index"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
    DeletedAt   gorm.DeletedAt `gorm:"index"`

    // Foreign keys
    Assignment Assignment `gorm:"foreignKey:AssignmentID"`
    Student    User       `gorm:"foreignKey:StudentID"`
}


// Create DTO
type CreateAssignmentSubmissionRequest struct {
    AssignmentID   uuid.UUID       `json:"assignment_id" binding:"required"`
    StudentID      uuid.UUID       `json:"student_id" binding:"required"`
    SubmissionType string          `json:"submission_type" binding:"required,oneof=text file url code video audio image document presentation multiple"`
    TextContent    string          `json:"text_content,omitempty"`
    FileURL        string          `json:"file_url,omitempty"`
    ExternalURL    string          `json:"external_url,omitempty"`
    CodeRepoURL    string          `json:"code_repo_url,omitempty"`
    Metadata       json.RawMessage `json:"metadata,omitempty"` // Raw JSON
    Status         string          `json:"status,omitempty"`   // Defaults to "submitted"
}

// Response DTO
type AssignmentSubmissionResponse struct {
    ID             uuid.UUID       `json:"id"`
    AssignmentID   uuid.UUID       `json:"assignment_id"`
    StudentID      uuid.UUID       `json:"student_id"`
    SubmissionType string          `json:"submission_type"`
    TextContent    string          `json:"text_content,omitempty"`
    FileURL        string          `json:"file_url,omitempty"`
    ExternalURL    string          `json:"external_url,omitempty"`
    CodeRepoURL    string          `json:"code_repo_url,omitempty"`
    Metadata       json.RawMessage `json:"metadata,omitempty"`
    Status         string          `json:"status"`
    SubmittedAt    time.Time       `json:"submitted_at"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
}
package models

import (
	"time"
	"github.com/google/uuid"


)

type Assignment struct {
    ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

    CourseID  uuid.UUID  `gorm:"type:uuid;not null;index"`
    ChapterID *uuid.UUID `gorm:"type:uuid;index"`
    LessonID  *uuid.UUID `gorm:"type:uuid;index"`
    PublisherID  *uuid.UUID `gorm:"type:uuid;index"`

    Title       string   `gorm:"type:varchar(255);not null"`
    Slug        string   `gorm:"type:varchar(300);uniqueIndex;not null"`
    Description string   `gorm:"type:text"`
	Content            string         `gorm:"type:text"`                               // Text content/description
    SubmissionType     string         `gorm:"type:varchar(50);default:'text';check:submission_type IN ('text', 'file', 'url', 'code', 'video', 'audio', 'image', 'document', 'presentation', 'multiple')"`    
    Status             string         `gorm:"type:varchar(20);default:'submitted';check:status IN ('draft', 'submitted', 'under_review', 'rejected')"`

    Type    string    `gorm:"type:varchar(50);default:'homework';check:type IN ('homework','project','essay','quiz','exam','lab','presentation','discussion','peer_review','group','research','creative')"`
    DueDate time.Time `gorm:"not null;index"`

    // Foreign key for Publisher
    Publisher   User      `gorm:"foreignKey:PublisherID"`
    Course      Course     `gorm:"foreignKey:CourseID"`
    Chapter     Chapter    `gorm:"foreignKey:ChapterID"`
    Lesson      Lesson     `gorm:"foreignKey:LessonID"`

    CreatedAt   time.Time
    UpdatedAt   time.Time
    ApprovedAt  *time.Time
    PublishedAt *time.Time
    ArchivedAt  *time.Time
}

type AssignmentInput struct {
	CourseID  uuid.UUID  `json:"course_id" binding:"required"`
	ChapterID *uuid.UUID `json:"chapter_id,omitempty"`
	LessonID  *uuid.UUID `json:"lesson_id,omitempty"`
	PublisherID  *uuid.UUID `json:"publisher_id,omitempty"`

	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description,omitempty"`

	DueDate time.Time `json:"due_date" binding:"required"`

	Content          string `json:"content,omitempty"`
	Type             string `json:"type,omitempty"`
	SubmissionType   string `json:"submission_type,omitempty"`
	Status           string `json:"status,omitempty"`
}


type AssignmentListResponse struct {
    ID         uuid.UUID  `json:"id"`
	CourseID   uuid.UUID  `json:"course_id"`
	ChapterID  *uuid.UUID `json:"chapter_id,omitempty"`
	LessonID   *uuid.UUID `json:"lesson_id,omitempty"`
	PublisherID *uuid.UUID `json:"publisher_id,omitempty"`

	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Content     string    `json:"content,omitempty"`
	Type            string    `json:"type"`
	SubmissionType   string `json:"submission_type,omitempty"`
	Status           string `json:"status,omitempty"`
	DueDate         time.Time `json:"due_date"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	ApprovedAt  *time.Time `json:"approved_at,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	ArchivedAt  *time.Time `json:"archived_at,omitempty"`
}

type AssignmentViewResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Type        string    `json:"type"`
	SubmissionType   string `json:"submission_type,omitempty"`
	Content     string    `json:"content,omitempty"`
	Status           string `json:"status,omitempty"`
	DueDate     time.Time `json:"due_date"`


	Course      CourseResponse                   `json:"course"`
	Chapter     *ChapterResponse                 `json:"chapter,omitempty"`
	Lesson      *LessonResponse                  `json:"lesson,omitempty"`
	Publisher   UserResponse                     `json:"publisher"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AssignmentUpdateInput struct {
	CourseID  uuid.UUID  `json:"course_id" binding:"required"`
	ChapterID *uuid.UUID `json:"chapter_id,omitempty"`
	LessonID  *uuid.UUID `json:"lesson_id,omitempty"`

	Title       string `json:"title" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Description string `json:"description,omitempty"`

	DueDate time.Time `json:"due_date" binding:"required"`

	Content          string `json:"content,omitempty"`
	Type             string `json:"type,omitempty"`
	SubmissionType   string `json:"submission_type,omitempty"`
	Status           string `json:"status,omitempty"`
}
// TableName specifies the table name
func (Assignment) TableName() string {
    return "assignments"
}
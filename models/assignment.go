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

    Title       string   `gorm:"type:varchar(255);not null"`
    Slug        string   `gorm:"type:varchar(300);uniqueIndex;not null"`
    Description string   `gorm:"type:text"`
    LearningObjectives []string `gorm:"type:text[]"`

    Type    string    `gorm:"type:varchar(50);default:'homework';check:type IN ('homework','project','essay','quiz','exam','lab','presentation','discussion','peer_review','group','research','creative')"`
    DueDate time.Time `gorm:"not null;index"`

    // Foreign key for Publisher
    PublishedBy uuid.UUID `gorm:"type:uuid;not null" json:"published_by"`
    Publisher   User      `gorm:"foreignKey:PublishedBy" json:"publisher,omitempty"`

    Course      Course     `gorm:"foreignKey:CourseID"`
    Chapter     Chapter    `gorm:"foreignKey:ChapterID"`
    Lesson      Lesson     `gorm:"foreignKey:LessonID"`
    Submissions []AssignmentSubmission `gorm:"foreignKey:AssignmentID"`

    CreatedAt   time.Time
    UpdatedAt   time.Time
    ApprovedAt  *time.Time
    PublishedAt *time.Time
    ArchivedAt  *time.Time
}


// TableName specifies the table name
func (Assignment) TableName() string {
    return "assignments"
}
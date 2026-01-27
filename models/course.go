package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type Course struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Image       string    `gorm:"size:255" json:"image"`
	VideoURL    string    `gorm:"size:255" json:"video_url"`
	TutorID     uuid.UUID `gorm:"type:uuid;not null" json:"tutor_id"`
	LearningOutcomes datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"learning_outcomes"`
	Requirements datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"requirements"`

	// Relationships
    Products []CourseProductTable `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
    Categories []CourseCategoryTable `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CourseResponse struct {
	ID               uuid.UUID      `json:"id"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	Image            string         `json:"image"`
	VideoURL         string         `json:"video_url"`
	TutorID          uuid.UUID      `json:"tutor_id"`
	LearningOutcomes datatypes.JSON `json:"learning_outcomes"`
	Requirements     datatypes.JSON `json:"requirements"`
}
type CourseMiniResponse struct {
	ID               uuid.UUID      `json:"id"`
	Title            string         `json:"title"`

}





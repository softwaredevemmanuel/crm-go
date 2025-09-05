package models

import (
	"time"
	"github.com/google/uuid"
)

type Course struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Image       string    `gorm:"size:255" json:"image"`
	VideoURL    string    `gorm:"size:255" json:"video_url"`
	TutorID     uuid.UUID `gorm:"type:uuid;not null" json:"tutor_id"`
	 // Relationships
    Products []CourseProduct `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

package models

import (
	"time"
)

type Course struct {
	ID          string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Image       string    `gorm:"size:255" json:"image"`
	VideoURL    string    `gorm:"size:255" json:"video_url"`
	TutorID     string     `gorm:"type:varchar(255)" json:"tutor_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

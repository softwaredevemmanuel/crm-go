package models

import (
	"time"

	"github.com/google/uuid"
)

type TutorRating struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TutorID   uuid.UUID `gorm:"type:uuid;not null" json:"tutor_id"`
	StudentID uuid.UUID `gorm:"type:uuid;not null" json:"student_id"`
	CourseID  *uuid.UUID `gorm:"type:uuid" json:"course_id,omitempty"`
	Rating    int       `gorm:"type:int;not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Review    string    `gorm:"type:text" json:"review"`
	Anonymous bool      `gorm:"default:false" json:"anonymous"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TutorRating) TableName() string {
	return "tutor_ratings"
}	
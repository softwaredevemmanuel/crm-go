package models

import (
	"time"
	"github.com/google/uuid"


)

type CourseRating struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID    uuid.UUID `gorm:"type:uuid;not null"`  // student who gives feedback
    CourseID  uuid.UUID `gorm:"type:uuid"`           // optional: feedback on a course
    Rating    int       `gorm:"not null"`            // 1–5 stars
    Comment   string    `gorm:"type:text"`
    CreatedAt time.Time
}


func (CourseRating) TableName() string {
	return "course_ratings"
}

package models

import (
	"time"
	"github.com/google/uuid"


)

type StudentBadge struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    StudentID uuid.UUID `gorm:"type:uuid;not null"`  // reference to student
    BadgeID   uuid.UUID `gorm:"type:uuid;not null"`  // reference to badge
    AwardedAt time.Time `gorm:"autoCreateTime"`      // when badge was awarded
}



func (StudentBadge) TableName() string {
	return "student_badges"
}

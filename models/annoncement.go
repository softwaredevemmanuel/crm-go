package models

import (
	"time"

	"github.com/google/uuid"
)

type Announcement struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title     string    `gorm:"type:varchar(255);not null" json:"title"`
	Message   string    `gorm:"type:text;not null" json:"message"`
	Type      string    `gorm:"type:varchar(50);not null" json:"type"` // general, update, maintenance, urgent
	Audience  string    `gorm:"type:varchar(50);not null" json:"audience"` // all, students, tutors, admins
	CreatedBy uuid.UUID `gorm:"type:uuid;not null" json:"created_by"`

	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`

	IsPinned bool `gorm:"default:false" json:"is_pinned"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Announcement) TableName() string {
	return "announcements"
}	
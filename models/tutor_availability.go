package models

import (
	"time"

	"github.com/google/uuid"
)

type TutorAvailability struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TutorID    uuid.UUID `gorm:"type:uuid;not null" json:"tutor_id"`
	DayOfWeek  string    `gorm:"type:varchar(20);not null" json:"day_of_week"`
	StartTime  string    `gorm:"type:varchar(10);not null" json:"start_time"`
	EndTime    string    `gorm:"type:varchar(10);not null" json:"end_time"`
	Timezone   string    `gorm:"type:varchar(50);not null" json:"timezone"`
	IsRecurring bool     `gorm:"default:true" json:"is_recurring"`
	IsAvailable bool     `gorm:"default:true" json:"is_available"`
	Notes      string    `gorm:"type:text" json:"notes"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TutorAvailability) TableName() string {
	return "tutor_availabilities"
}	
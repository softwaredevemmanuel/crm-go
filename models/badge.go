package models

import (
	"time"
	"github.com/google/uuid"


)
type Badge struct {
    ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    Name        string    `gorm:"type:varchar(100);not null;unique"` // e.g., "Top Scorer"
    Description string    `gorm:"type:text"`                         // what the badge means
    IconURL     string    `gorm:"type:varchar(255)"`                 // optional icon
    CreatedAt   time.Time
    UpdatedAt   time.Time
}


func (Badge) TableName() string {
	return "badges"
}

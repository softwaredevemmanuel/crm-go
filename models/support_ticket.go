package models

import (
	"time"
	"github.com/google/uuid"


)

type SupportTicket struct {
    ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    UserID      uuid.UUID `gorm:"type:uuid;not null"`  // who raised it (student/tutor)
    Title       string    `gorm:"type:varchar(150);not null"`
    Description string    `gorm:"type:text;not null"`
    Status      string    `gorm:"type:varchar(50);default:'open'"` // open, in-progress, resolved
    Priority    string    `gorm:"type:varchar(50);default:'medium'"` // low, medium, high
    CreatedAt   time.Time
    UpdatedAt   time.Time
}



func (SupportTicket) TableName() string {
	return "support_tickets"
}
package models

import (
	"time"
	"github.com/google/uuid"


)

type SupportResponse struct {
    ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    TicketID  uuid.UUID `gorm:"type:uuid;not null"`
    ResponderID uuid.UUID `gorm:"type:uuid;not null"` // admin or tutor responding
    Message   string    `gorm:"type:text;not null"`
    CreatedAt time.Time
}



func (SupportResponse) TableName() string {
	return "support_responses"
}
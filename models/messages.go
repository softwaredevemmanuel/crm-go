package models

import (
	"time"
	"github.com/google/uuid"


)

type Message struct {
    ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
    SenderID   uuid.UUID `gorm:"type:uuid;not null"`   // student or tutor
    ReceiverID uuid.UUID `gorm:"type:uuid;not null"`
    Content    string    `gorm:"type:text;not null"`
    IsRead     bool      `gorm:"default:false"`
    CreatedAt  time.Time
}



func (Message) TableName() string {
	return "messages"
}
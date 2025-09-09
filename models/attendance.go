package models

import (
	"time"
	"github.com/google/uuid"

)
type LiveClassAttendance struct {
    ID             uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    SessionID      uuid.UUID  `gorm:"type:uuid;not null;index"`
    UserID         uuid.UUID  `gorm:"type:uuid;not null;index"`
    JoinTime       time.Time
    LeaveTime      *time.Time
    Duration       int        `gorm:"default:0"` // seconds
    ConnectionType string     `gorm:"type:varchar(20)"` // audio, video, both
    // ... other fields
}
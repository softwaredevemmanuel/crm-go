package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DeletedRecord struct {
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`

	EntityType string    `gorm:"type:varchar(100);index"` // "topic", "module", "course_material"
	EntityID   uuid.UUID `gorm:"type:uuid;index"`

	Data datatypes.JSON `gorm:"type:jsonb"` // Full snapshot of the record

	DeletedBy *uuid.UUID `gorm:"type:uuid"` // optional (user/admin)
	Reason    string     `gorm:"type:text"`

	DeletedAt time.Time `gorm:"autoCreateTime"`
}

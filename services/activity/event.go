package activity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	UserID     uuid.UUID
	Action     string
	EntityID   uuid.UUID
	EntityType string
	Details    string
	Metadata   map[string]interface{}

	IPAddress string
	UserAgent string
	Tx        *gorm.DB
}

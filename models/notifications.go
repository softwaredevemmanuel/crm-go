// models/notification.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Notification struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Type        string         `gorm:"type:varchar(50);not null;index" json:"type"`
	Title       string         `gorm:"type:varchar(255);not null" json:"title"`
	Message     string         `gorm:"type:text;not null" json:"message"`
	Link        string         `gorm:"type:varchar(500)" json:"link"`
	IsRead      bool           `gorm:"default:false;index" json:"is_read"`
	IsDismissed bool           `gorm:"default:false" json:"is_dismissed"`
	Priority    string         `gorm:"type:varchar(20);default:'normal'" json:"priority"`
	SentAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;index" json:"sent_at"`
	ReadAt      *time.Time     `json:"read_at,omitempty"`
	DismissedAt *time.Time     `json:"dismissed_at,omitempty"`
	ExpiresAt   *time.Time     `json:"expires_at,omitempty"`
	Metadata    string         `gorm:"type:jsonb" json:"metadata"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	User    User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Creator User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName specifies the table name
func (Notification) TableName() string {
	return "notifications"
}
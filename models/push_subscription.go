// models/push_subscription.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PushSubscription struct {
	ID        uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	Endpoint  string         `gorm:"type:text;not null" json:"endpoint"`
	P256dh    string         `gorm:"type:text;not null" json:"p256dh"`
	Auth      string         `gorm:"type:text;not null" json:"auth"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PushSubscription) TableName() string {
	return "push_subscriptions"
}
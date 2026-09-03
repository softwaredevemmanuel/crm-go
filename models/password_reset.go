// models/password_reset.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)
type PasswordReset struct {
	ID         uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"user_id"`
	FirstName  string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName   string         `gorm:"type:varchar(100);not null" json:"last_name"`
	Token      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"token"`
	Email      string         `gorm:"type:varchar(255);not null" json:"email"`
	ExpiresAt  time.Time      `gorm:"not null" json:"expires_at"`
	Used       bool           `gorm:"default:false" json:"used"`
	CreatedAt  time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PasswordReset) TableName() string {
	return "password_resets"
}
// models/session.go
package models

import (
	"time"
	
	"github.com/google/uuid"
)

type UserSession struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	SessionToken string         `gorm:"type:varchar(500);uniqueIndex;not null"`
	UserID       uuid.UUID      `gorm:"type:uuid;not null;index"`
	UserAgent    string         `gorm:"type:text"`
	UserIP       string         `gorm:"type:varchar(45)"`
	DeviceType   string         `gorm:"type:varchar(50)"`
	DeviceOS     string         `gorm:"type:varchar(100)"`
	Browser      string         `gorm:"type:varchar(100)"`
	IsActive     bool           `gorm:"default:true;index"`
	LoginType    string         `gorm:"type:varchar(50);default:'password'"`
	IssuedAt     time.Time      `gorm:"not null"`
	ExpiresAt    time.Time      `gorm:"not null;index"`
	LastUsedAt   time.Time      `gorm:"default:CURRENT_TIMESTAMP;index"`
	CreatedAt    time.Time
    LoggedOutAt  *time.Time     `gorm:"index"`  // ✅ Changed to pointer

	User         User           `gorm:"foreignKey:UserID"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
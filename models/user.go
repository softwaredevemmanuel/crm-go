package models

import (
	"time"
	"github.com/google/uuid"

)

type Role string

const (
	Student Role = "student"
	Tutor   Role = "tutor"
	Admin   Role = "admin"
)

type User struct {
	ID        uuid.UUID  `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName string    `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName  string    `gorm:"type:varchar(100);not null" json:"last_name"`
	Email     string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string    `gorm:"type:text" json:"-"` // - means exclude from JSON
	Picture   string    `gorm:"type:text" json:"picture,omitempty"`
	Provider  string    `gorm:"type:varchar(50);default:'local'" json:"provider"`
	Role      Role      `gorm:"type:varchar(10);not null"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	
}

type PasswordReset struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"type:uuid;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}


func (User) TableName() string {
	return "users"
}


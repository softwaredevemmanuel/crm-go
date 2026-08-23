package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"


)



type User struct {
	ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	FirstName  string         `gorm:"type:varchar(100);not null" json:"first_name"`
	LastName   string         `gorm:"type:varchar(100);not null" json:"last_name"`
	MiddleName string         `gorm:"type:varchar(100)" json:"middle_name"`
	Email      string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password   string         `gorm:"type:text" json:"-"`
	LoginID    string         `gorm:"type:text" json:"login_id,omitempty"`
	Picture    string         `gorm:"type:text" json:"picture,omitempty"`
	Provider   string         `gorm:"type:varchar(50);default:'local'" json:"provider"`
	Role       string         `gorm:"type:varchar(20);default:'user'" json:"role"`
	Phone      string         `gorm:"type:varchar(20)" json:"phone"`
	DOB 		*time.Time `gorm:"type:date" json:"dob,omitempty"`
	IsVerified bool        `gorm:"default:false" json:"is_verified"`
	IsActive   bool        `gorm:"default:false" json:"is_active"`
	Location   string      `gorm:"type:varchar(255)" json:"location"`
	LastLoginAt *time.Time `json:"last_login_at"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}


type PasswordReset struct {
	ID        string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	UserID    string    `gorm:"type:uuid;not null"`
	Token     string    `gorm:"uniqueIndex;not null"`
	ExpiresAt time.Time `gorm:"not null"`
}


type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Role      string      `json:"role"`
	DOB        string `json:"dob"` // Format: YYYY-MM-DD

}


func (User) TableName() string {
	return "users"
}


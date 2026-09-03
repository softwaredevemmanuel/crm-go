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
	Position   string         `gorm:"type:varchar(20);default:'user'" json:"position"`
	Phone      string         `gorm:"type:varchar(20)" json:"phone"`
	DOB 		*time.Time `gorm:"type:date" json:"dob,omitempty"`
	IsVerified bool        `gorm:"default:false" json:"is_verified"`
	IsActive   bool        `gorm:"default:false" json:"is_active"`
	Location   string      `gorm:"type:varchar(255)" json:"location"`
	LastLoginAt *time.Time `json:"last_login_at"`
	EmailVerifiedAt   *time.Time     `json:"email_verified_at,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}






func (User) TableName() string {
	return "users"
}


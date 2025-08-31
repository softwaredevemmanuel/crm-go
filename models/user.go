package models

import (
	"time"
)

type Role string

const (
	Student Role = "student"
	Tutor   Role = "tutor"
	Admin   Role = "admin"
)

type User struct {
	ID string `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name      string    `gorm:"not null"`
	Email     string    `gorm:"uniqueIndex;not null"`
	Password  string    `gorm:"not null"`
	Role      Role      `gorm:"type:varchar(10);not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}



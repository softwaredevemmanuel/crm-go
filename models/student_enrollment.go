// models/student_enrollment.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentEnrollment struct {
	ID             uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StudentID      uuid.UUID      `gorm:"type:uuid;not null;index:idx_student_arm_unique,unique" json:"student_id"`
	ArmID          uuid.UUID      `gorm:"type:uuid;not null;index:idx_student_arm_unique,unique" json:"arm_id"`
	Status         string         `gorm:"type:varchar(20);not null;default:'active';check:status IN ('active', 'inactive', 'transferred', 'graduated', 'withdrawn')" json:"status"`
	GraduationDate *time.Time     `gorm:"type:date" json:"graduation_date,omitempty"`
	Notes          string         `gorm:"type:text" json:"notes"`
	IsVerified     bool           `gorm:"not null;default:false" json:"is_verified"`
	CreatedBy      uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Student User `gorm:"foreignKey:StudentID" json:"student,omitempty"`
	Arm     Arm  `gorm:"foreignKey:ArmID" json:"arm,omitempty"`
}

func (StudentEnrollment) TableName() string {
	return "student_enrollments"
}
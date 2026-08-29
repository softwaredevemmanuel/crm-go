// models/arm_class_teacher.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArmClassTeacher struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ArmID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"arm_id"`
	TeacherID uuid.UUID      `gorm:"type:uuid;not null;index" json:"teacher_id"`
	Status    string         `gorm:"type:varchar(20);default:'active'" json:"status"` // active, inactive
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Arm     Arm       `gorm:"foreignKey:ArmID" json:"arm,omitempty"`
	Teacher User      `gorm:"foreignKey:TeacherID" json:"teacher,omitempty"`
}

func (ArmClassTeacher) TableName() string {
	return "arm_class_teachers"
}
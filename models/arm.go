// models/arm.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Arm struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Code        string         `gorm:"type:varchar(50)" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	GradeID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"grade_id"`
	Capacity    int            `gorm:"type:int;default:0" json:"capacity"`
	Status      string         `gorm:"type:varchar(20);default:'active'" json:"status"`
	CreatedBy   uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Grade              *ClassGrade          `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
	ClassTeacher       *ArmClassTeacher     `gorm:"foreignKey:ArmID" json:"class_teacher,omitempty"`
	StudentEnrollments []StudentEnrollment  `gorm:"foreignKey:ArmID" json:"student_enrollments,omitempty"`
}

func (Arm) TableName() string {
	return "arms"
}
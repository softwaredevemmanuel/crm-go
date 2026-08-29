// models/arm.go
package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Arm struct {
	ID          uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Code        string         `gorm:"type:varchar(50)" json:"code"`
	GradeID     uuid.UUID      `gorm:"type:uuid;not null;index" json:"grade_id"`
	Status      string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'inactive', 'archived')" json:"status"`
	Capacity    int            `gorm:"default:30" json:"capacity"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Grade ClassGrade `gorm:"foreignKey:GradeID" json:"grade,omitempty"`
}

// TableName specifies the table name
func (Arm) TableName() string {
	return "arms"
}
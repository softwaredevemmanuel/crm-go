// models/class_grade.go
package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type ClassGrade struct {
	ID                uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	Code              string         `gorm:"type:varchar(50);not null;" json:"code"`
	Level             int            `gorm:"not null" json:"level"` // 1, 2, 3, 4, 5, 6
	Description       string         `gorm:"type:text" json:"description"`
	Status            string         `gorm:"type:varchar(20);default:'active';check:status IN ('active', 'inactive', 'completed')" json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

}

// TableName specifies the table name
func (ClassGrade) TableName() string {
	return "class_grades"
}

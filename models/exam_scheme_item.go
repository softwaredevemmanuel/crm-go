// models/exam_scheme_item.go
package models

import (
	"github.com/google/uuid"
)

type ExamSchemeItem struct {
	ExamID             uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"exam_id"`
	SchemeOfWorkItemID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"scheme_of_work_item_id"`

	// Relationships
	Exam             Exam             `gorm:"foreignKey:ExamID" json:"exam,omitempty"`
	SchemeOfWorkItem SchemeOfWorkItem `gorm:"foreignKey:SchemeOfWorkItemID" json:"scheme_of_work_item,omitempty"`
}

func (ExamSchemeItem) TableName() string {
	return "exam_scheme_items"
}
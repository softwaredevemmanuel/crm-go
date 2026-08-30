// models/test_scheme_item.go
package models

import (
	"github.com/google/uuid"
)

type TestSchemeItem struct {
	TestID             uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"test_id"`
	SchemeOfWorkItemID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"scheme_of_work_item_id"`

	// Relationships
	Test             Test             `gorm:"foreignKey:TestID" json:"test,omitempty"`
	SchemeOfWorkItem SchemeOfWorkItem `gorm:"foreignKey:SchemeOfWorkItemID" json:"scheme_of_work_item,omitempty"`
}

func (TestSchemeItem) TableName() string {
	return "test_scheme_items"
}
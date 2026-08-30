// models/lesson_plan.go
package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonPlan struct {
	ID                   uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	LessonID             uuid.UUID      `gorm:"type:uuid;unique;not null;index" json:"lesson_id"`
	PreviousKnowledge    string         `gorm:"type:text" json:"previous_knowledge"`
	BehaviouralObjectives string         `gorm:"type:text" json:"behavioural_objectives"`
	TeachingAids         string         `gorm:"type:text" json:"teaching_aids"`
	Introduction         string         `gorm:"type:text" json:"introduction"`
	LessonContent        string         `gorm:"type:text" json:"lesson_content"`
	TeacherActivities    string         `gorm:"type:text" json:"teacher_activities"`
	StudentActivities    string         `gorm:"type:text" json:"student_activities"`
	Conclusion           string         `gorm:"type:text" json:"conclusion"`
	Evaluation           string         `gorm:"type:text" json:"evaluation"`
	CreatedBy            uuid.UUID      `gorm:"type:uuid;not null" json:"created_by"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            time.Time      `json:"updated_at"`
	DeletedAt            gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Lesson  Lesson `gorm:"foreignKey:LessonID" json:"lesson,omitempty"`
	Creator User   `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

func (LessonPlan) TableName() string {
	return "lesson_plans"
}
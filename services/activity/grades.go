package activity

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GradeActivity struct {
	logger *Logger
}

func (a *GradeActivity) Created(
	tx *gorm.DB,
	userID uuid.UUID,
	grade models.Grade,
) error {

	metadata := map[string]interface{}{
		"grade_id":   grade.ID,
		"tutor_id":   grade.TutorID,
		"course_id":  grade.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:   userID,
			Action:   models.ActionGradeCreate,
			EntityID: grade.ID,
			EntityType: "grades",
			Details:  fmt.Sprintf("Created grade for student: %s", grade.StudentID),
			Metadata: metadata,

		},
	)
}
	

func (a *GradeActivity) Updated(
	tx *gorm.DB,
	userID uuid.UUID,
	grade models.Grade,
) error {

	metadata := map[string]interface{}{
		"grade_id":   grade.ID,
		"student_id": grade.StudentID,
		"course_id":  grade.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:   userID,
			Action:   models.ActionGradeUpdate,
			EntityID: grade.ID,
			EntityType: "grade",
			Details:  fmt.Sprintf("Updated grade for student: %s", grade.StudentID),
			Metadata: metadata,
		},
	)
}

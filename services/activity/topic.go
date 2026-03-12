package activity

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LessonActivity struct {
	logger *Logger
}

func (a *LessonActivity) Created(
	tx *gorm.DB,
	userID uuid.UUID,
	lesson models.Lesson,
) error {

	metadata := map[string]interface{}{
		"lesson_id": lesson.ID,
		"course_id": lesson.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:     userID,
			Action:     models.ActionLessonCreate,
			EntityID:   lesson.ID,
			EntityType: "lesson",
			Details:    fmt.Sprintf("Created lesson: %s", lesson.Title),
			Metadata:   metadata,
		},
	)
}

func (a *LessonActivity) Updated(
	tx *gorm.DB,
	userID uuid.UUID,
	lesson models.Lesson,
) error {

	metadata := map[string]interface{}{
		"lesson_id": lesson.ID,
		"course_id": lesson.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:     userID,
			Action:     models.ActionLessonUpdate,
			EntityID:   lesson.ID,
			EntityType: "lesson",
			Details:    fmt.Sprintf("Updated lesson: %s", lesson.Title),
			Metadata:   metadata,
		},
	)
}

package activity

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ObjectiveActivity struct {
	logger *Logger
}

func (a *ObjectiveActivity) Created(
	tx *gorm.DB,
	userID uuid.UUID,
	objective models.ObjectiveQuestion,
) error {

	metadata := map[string]interface{}{
		"objective_id": objective.ID,
		"creator_id":   objective.CreatedBy,
		"course_id":    objective.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:   userID,
			Action:   models.ActionObjectiveCreate,
			EntityID: objective.ID,
			EntityType: "objectives",
			Details:  fmt.Sprintf("Created objective question: %s", objective.QuestionText),
			Metadata: metadata,

		},
	)
}

	
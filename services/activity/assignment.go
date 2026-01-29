package activity

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AssignmentActivity struct {
	logger *Logger
}

func (a *AssignmentActivity) Submitted(
	tx *gorm.DB,
	studentID uuid.UUID,
	assignment models.Assignment,
	submission models.AssignmentSubmission,
	extra map[string]interface{},
) error {

	metadata := map[string]interface{}{
		"assignment_id": assignment.ID,
		"course_id":     assignment.CourseID,
		"submission_id": submission.ID,
		"status":        submission.Status,
	}

	for k, v := range extra {
		metadata[k] = v
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:     studentID,
			Action:     models.ActionAssignmentSubmit,
			EntityID:   submission.ID,
			EntityType: "assignment_submission",
			Details:    fmt.Sprintf("Submitted assignment: %s", assignment.Title),
			Metadata:   metadata,
		},
	)
}

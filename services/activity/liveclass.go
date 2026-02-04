package activity

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type LiveClassActivity struct {
	logger *Logger
}

func (a *LiveClassActivity) Created(
	tx *gorm.DB,
	userID uuid.UUID,
	liveClass models.LiveClass,
) error {

	metadata := map[string]interface{}{
		"live_class_id": liveClass.ID,
		"tutor_id":      liveClass.TutorID,
		"course_id":     liveClass.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:   userID,
			Action:   models.ActionLiveClassCreate,
			EntityID: liveClass.ID,
			EntityType: "live_classes",
			Details:  fmt.Sprintf("Created live class: %s", liveClass.Title),
			Metadata: metadata,

		},
	)
}
	

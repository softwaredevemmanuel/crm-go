package activity

import (
	"context"
	"fmt"

	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TopicActivity struct {
	logger *Logger
}

func (a *TopicActivity) Created(
	tx *gorm.DB,
	userID uuid.UUID,
	topic models.Topic,
) error {

	metadata := map[string]interface{}{
		"topic_id": topic.ID,
		"course_id": topic.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:   userID,
			Action:   models.ActionTopicCreate,
			EntityID: topic.ID,
			EntityType: "topic",
			Details:  fmt.Sprintf("Created topic: %s", topic.Title),
			Metadata: metadata,
		},
	)
}


func (a *TopicActivity) Updated(
	tx *gorm.DB,
	userID uuid.UUID,
	topic models.Topic,
) error {

	metadata := map[string]interface{}{
		"topic_id": topic.ID,
		"course_id": topic.CourseID,
	}

	return a.logger.LogWithTx(
		context.Background(),
		tx,
		Event{
			UserID:   userID,
			Action:   models.ActionTopicUpdate,
			EntityID: topic.ID,
			EntityType: "topic",
			Details:  fmt.Sprintf("Updated topic: %s", topic.Title),
			Metadata: metadata,
		},
	)
}

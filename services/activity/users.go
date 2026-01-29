package activity

import (
	"crm-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserActivity struct {
	logger *Logger
}

func (u *UserActivity) Login(
	c *gin.Context,
	userID uuid.UUID,
	success bool,
) error {

	details := "User logged in"
	if !success {
		details = "Failed login attempt"
	}

	return u.logger.LogFromRequest(
		c,
		Event{
			UserID:     userID,
			Action:     models.ActionLogin,
			EntityID:   userID,
			EntityType: "user",
			Details:    details,
			Metadata: map[string]interface{}{
				"success": success,
			},
		},
	)
}

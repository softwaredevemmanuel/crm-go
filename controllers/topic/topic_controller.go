package controllers

import (
	"net/http"

	"crm-go/models"
	"crm-go/services"
    "crm-go/services/activity"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TopicController struct {
	db       *gorm.DB
	topicService *services.TopicService
	activity *activity.Service

}

func NewTopicController(db *gorm.DB, topicService *services.TopicService, activitySvc *activity.Service) *TopicController {
	return &TopicController{
		db:             db,
		topicService:   topicService,
		activity:       activitySvc,
	}
}

// CreateTopic creates a new topic
// @Summary Create a new topic
// @Description Create a new topic
// @Tags topics
// @Accept json
// @Produce json
// @Param topic body models.TopicInput true "Topic"
// @Success 201 {object} models.SuccessResponse "Topic created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/topics [post]
// @Security BearerAuth
func (ctl *TopicController) CreateTopic(c *gin.Context) {
	var req models.TopicInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

tx := ctl.db.Begin()

topic, err := ctl.topicService.CreateTopic(req)
if err != nil {
	tx.Rollback()
	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	return
}

	// 🔥 ACTIVITY LOG — CLEAN & REUSABLE
	_ = ctl.activity.Topics.Created(
		tx,
		req.TutorID,	
		*topic,	
	)

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Topic created successfully",
		"data":    topic,
	})
}

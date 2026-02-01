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
    db           *gorm.DB
    topicService *services.TopicService
    activity     *activity.Service
}

func NewTopicController(db *gorm.DB, topicService *services.TopicService, activitySvc *activity.Service) *TopicController {
    return &TopicController{
        db:           db,
        topicService: topicService,
        activity:     activitySvc,
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

    
    // Start transaction
    tx := ctl.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()
    
    // Create topic - need to use transaction version
    response, err := ctl.topicService.CreateTopicWithTx(tx, req)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    // Convert response back to Topic model for activity logging
    topicModel := models.Topic{
        ID:          response.ID,
        CourseID:    response.CourseID,
        ChapterID:   response.ChapterID,
        TutorID:     response.TutorID,
        Title:       response.Title,
        Description: response.Description,
        Order:       response.Order,
        CreatedAt:   response.CreatedAt,
        UpdatedAt:   response.UpdatedAt,
    }
    
    // Activity logging with error handling
    if err := ctl.activity.Topics.Created(tx, req.TutorID, topicModel); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to log activity: " + err.Error(),
        })
        return
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save changes: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Topic created successfully",
        "data":    response,
    })
}


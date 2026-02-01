package controllers

import (
  "net/http"
    
    "crm-go/models"
    "github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateTopic updates an existing topic
//@Summary Update a topic
//@Description Update a topic by ID
//@Tags topics
//@Accept json
//@Produce json
//@Param id path string true "Topic ID"
//@Param topic body models.TopicInput true "Topic"
//@Success 200 {object} models.TopicResponse "Topic updated successfully"
//@Failure 400 {object} models.ErrorResponse
//@Failure 404 {object} models.ErrorResponse
//@Router /api/topics/{id} [put]
//@Security BearerAuth
func (ctl *TopicController) UpdateTopic(ctx *gin.Context) {
    topicID, err := uuid.Parse(ctx.Param("id"))
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid topic ID",
        })
        return
    }
    
    var req models.TopicInput
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
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
    
    // Update topic - need to use transaction version
    updatedTopic, err := ctl.updateTopicService.UpdateTopicWithTx(tx, topicID, req)
    if err != nil {
        tx.Rollback()
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    topicModel := models.Topic{
        ID:          updatedTopic.ID,
        CourseID:    updatedTopic.CourseID,
        ChapterID:   updatedTopic.ChapterID,
        Title:       updatedTopic.Title,
        Description: updatedTopic.Description,
        Order:       updatedTopic.Order,
        CreatedAt:   updatedTopic.CreatedAt,
        UpdatedAt:   updatedTopic.UpdatedAt,
    }
    
    // Log activity with error handling
    if err := ctl.activity.Topics.Updated(tx, req.TutorID, topicModel); err != nil {
        tx.Rollback()
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to log activity: " + err.Error(),
        })
        return
    }
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to commit changes: " + err.Error(),
        })
        return
    }
    
    ctx.JSON(http.StatusOK, gin.H{
        "message": "Topic updated successfully",
        "data":    updatedTopic,
    })
}
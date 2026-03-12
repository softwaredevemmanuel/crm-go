package controllers

import (
	"net/http"

	"crm-go/models"

	"github.com/gin-gonic/gin"
)

// CreateLesson creates a new lesson
// @Summary Create a new lesson
// @Description Create a new lesson
// @Tags lessons
// @Accept json
// @Produce json
// @Param lesson body models.LessonInput true "Lesson"
// @Success 201 {object} models.SuccessResponse "Lesson created successfully"
// @Failure 400 {object} models.ErrorResponse "Bad request"
// @Router /api/lessons [post]
// @Security BearerAuth
func (ctl *LessonController) CreateLesson(c *gin.Context) {
	var req models.LessonInput

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

	// Create lesson - need to use transaction version
	response, err := ctl.createLessonService.CreateLessonWithTx(tx, req)
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert response back to Lesson model for activity logging
	lessonModel := models.Lesson{
		ID:          response.ID,
		CourseID:    response.CourseID,
		ModuleID:    response.ModuleID,
		TutorID:     response.TutorID,
		Title:       response.Title,
		Description: response.Description,
		Order:       response.Order,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	}

	// Activity logging with error handling
	if err := ctl.activity.Lessons.Created(tx, req.TutorID, lessonModel); err != nil {
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
		"message": "Lesson created successfully",
		"data":    response,
	})
}

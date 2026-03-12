package controllers

import (
	"net/http"

	"crm-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateLesson updates an existing lesson
// @Summary Update a lesson
// @Description Update a lesson by ID
// @Tags lessons
// @Accept json
// @Produce json
// @Param id path string true "Lesson ID"
// @Param lesson body models.LessonInput true "Lesson"
// @Success 200 {object} models.LessonResponse "Lesson updated successfully"
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/lessons/{id} [put]
// @Security BearerAuth
func (ctl *LessonController) UpdateLesson(ctx *gin.Context) {
	lessonID, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid lesson ID",
		})
		return
	}

	var req models.LessonInput
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

	// Update lesson - need to use transaction version
	updatedLesson, err := ctl.updateLessonService.UpdateLessonWithTx(tx, lessonID, req)
	if err != nil {
		tx.Rollback()
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	lessonModel := models.Lesson{
		ID:          updatedLesson.ID,
		CourseID:    updatedLesson.CourseID,
		ModuleID:    updatedLesson.ModuleID,
		Title:       updatedLesson.Title,
		Description: updatedLesson.Description,
		Order:       updatedLesson.Order,
		CreatedAt:   updatedLesson.CreatedAt,
		UpdatedAt:   updatedLesson.UpdatedAt,
	}

	// Log activity with error handling
	if err := ctl.activity.Lessons.Updated(tx, req.TutorID, lessonModel); err != nil {
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
		"message": "Lesson updated successfully",
		"data":    updatedLesson,
	})
}

package lesson

import (
	"crm-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"crm-go/config"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateLesson godoc
// @Summary      Create a lesson
// @Description  Create a new lesson under a chapter
// @Tags         Lessons
// @Accept       json
// @Produce      json
// @Param        lesson body models.LessonInput true "Lesson payload"
// @Success      201 {object} models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      409 {object} models.ConflictResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/lessons [post]
//@Security BearerAuth	
func CreateLesson(c *gin.Context) {
	var input models.LessonInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
		return
	}

	// 🔒 Duplicate check
	var existing models.Lesson
	if err := config.DB.
		Where("course_id = ? AND chapter_id = ? AND title = ?",
			input.CourseID, input.ChapterID, input.Title).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Lesson already exists for this chapter",
		})
		return
	}

	lesson := models.Lesson{
		ID:          uuid.New(),
		CourseID:    input.CourseID,
		ChapterID:   input.ChapterID,
		Title:       input.Title,
		ContentType: input.ContentType,
		ContentURL:  input.ContentURL,
	}

	if err := config.DB.Create(&lesson).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Lesson created successfully",
	})
}

// GetLessons godoc
// @Summary      List lessons
// @Description  Get all lessons
// @Tags         Lessons
// @Produce      json
// @Success      200 {array} models.LessonResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /lessons [get]
func GetLessons(c *gin.Context) {
	var lessons []models.Lesson

	if err := config.DB.
		Preload("Chapter", func(db *gorm.DB) *gorm.DB {
			return db.
				Select("id", "course_id", "title", "slug", "description",
					"chapter_number", "is_free", "status",
					"estimated_time", "total_lessons", "total_duration",
					"created_at", "updated_at").
				Preload("Course", func(db *gorm.DB) *gorm.DB {
					return db.Select(
						"id", "title", "description", "image",
						"video_url", "tutor_id",
						"learning_outcomes", "requirements",
						"created_at", "updated_at",
					)
				})
		}).
		Select(
			"id", "chapter_id", "course_id",
			"title", "content_type", "content_url",
			"created_at", "updated_at",
		).
		Order("created_at DESC").
		Find(&lessons).Error; err != nil {

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to fetch lessons",
		})
		return
	}

	c.JSON(http.StatusOK, lessons)
}



// GetLessonByID godoc
// @Summary      Get lesson details
// @Description  Get lesson with chapter and course details
// @Tags         Lessons
// @Produce      json
// @Param        id path string true "Lesson ID"
// @Success      200 {object} models.LessonResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Router       /lessons/{id} [get]
func GetLessonByID(c *gin.Context) {
	id := c.Param("id")

	lessonID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid lesson ID",
		})
		return
	}

	var lesson models.Lesson

	if err := config.DB.
		Preload("Chapter").
		Preload("Course").
		First(&lesson, "id = ?", lessonID).Error; err != nil {

		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Lesson not found",
		})
		return
	}

	c.JSON(http.StatusOK, lesson)
}


// UpdateLesson godoc
// @Summary      Update lesson
// @Description  Update lesson details
// @Tags         Lessons
// @Accept       json
// @Produce      json
// @Param        id path string true "Lesson ID"
// @Param        lesson body models.LessonInput true "Lesson payload"
// @Success      200 {object} models.LessonResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/lessons/{id} [put]
//@Security BearerAuth	
func UpdateLesson(c *gin.Context) {
	id := c.Param("id")

	lessonID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid lesson ID",
		})
		return
	}

	var lesson models.Lesson
	if err := config.DB.First(&lesson, "id = ?", lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Lesson not found",
		})
		return
	}

	var input models.LessonInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	if err := config.DB.Model(&lesson).Updates(models.Lesson{
		ChapterID:   input.ChapterID,
		CourseID:    input.CourseID,
		Title:       input.Title,
		ContentType: input.ContentType,
		ContentURL:  input.ContentURL,
	}).Error; err != nil {

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// reload relations
	config.DB.
		Preload("Chapter").
		Preload("Course").
		First(&lesson, "id = ?", lessonID)

	c.JSON(http.StatusOK, lesson)
}


// DeleteLesson godoc
// @Summary      Delete lesson
// @Description  Delete lesson by ID
// @Tags         Lessons
// @Produce      json
// @Param        id path string true "Lesson ID"
// @Success      200 {object} models.SuccessResponse
// @Failure      400 {object} models.ErrorResponse
// @Failure      404 {object} models.NotFoundResponse
// @Failure      500 {object} models.FailureResponse
// @Router       /api/lessons/{id} [delete]
//@Security BearerAuth	
func DeleteLesson(c *gin.Context) {
	id := c.Param("id")

	lessonID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "Invalid lesson ID",
		})
		return
	}

	result := config.DB.Delete(&models.Lesson{}, "id = ?", lessonID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to delete lesson",
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Lesson not found",
		})
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Lesson deleted successfully",
	})
}

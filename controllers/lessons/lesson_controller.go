package lesson

import (
	"crm-go/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"crm-go/config"
	"github.com/google/uuid"
	// "gorm.io/gorm"
	
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
	var existing models.Lessons
	if err := config.DB.
		Where("course_id = ? AND chapter_id = ? AND title = ?",
			input.CourseID, input.ChapterID, input.Title).
		First(&existing).Error; err == nil {

		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error: "Lesson already exists for this chapter",
		})
		return
	}

	lesson := models.Lessons{
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
// @Summary      Get lessons
// @Description  Get all lessons, optionally filtered by course or chapter
// @Tags         Lessons
// @Produce      json
// @Param        course_id   query   string  false  "Course ID"
// @Param        chapter_id  query   string  false  "Chapter ID"
// @Success      200 {array}  models.LessonResponse
// @Failure      400 {object} models.ErrorResponse
// @Router       /lessons [get]
func GetAllLessons(c *gin.Context) {
	var lessons []models.Lessons

	query := config.DB

	if courseID := c.Query("course_id"); courseID != "" {
		if uid, err := uuid.Parse(courseID); err == nil {
			query = query.Where("course_id = ?", uid)
		}
	}

	if chapterID := c.Query("chapter_id"); chapterID != "" {
		if uid, err := uuid.Parse(chapterID); err == nil {
			query = query.Where("chapter_id = ?", uid)
		}
	}

	if err := query.
		Order("created_at DESC").
		Find(&lessons).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch lessons",
		})
		return
	}

	responses := make([]models.LessonResponse, 0, len(lessons))
	for _, lesson := range lessons {
		responses = append(responses, models.LessonResponse{
			ID:          lesson.ID,
			ChapterID:   lesson.ChapterID,
			CourseID:    lesson.CourseID,
			Title:       lesson.Title,
			ContentType: lesson.ContentType,
			ContentURL:  lesson.ContentURL,
			CreatedAt:   lesson.CreatedAt,
			UpdatedAt:   lesson.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, responses)
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

	var lesson models.Lessons

	// Fetch lesson with relations
	if err := config.DB.
		Preload("Chapter").
		Preload("Chapter.Course").
		First(&lesson, "id = ?", lessonID).
		Error; err != nil {

		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Lesson not found",
		})
		return
	}

	// Map to DTO
	response := models.LessonViewResponse{
		ID:          lesson.ID,
		ChapterID:   lesson.ChapterID,
		CourseID:    lesson.CourseID,
		Title:       lesson.Title,
		ContentType: lesson.ContentType,
		ContentURL:  lesson.ContentURL,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,
		Course:      models.CourseMiniResponse{
			ID:    lesson.Chapter.Course.ID,
			Title: lesson.Chapter.Course.Title,
		},
		Chapter: &models.ChapterMiniResponse{
			ID:    lesson.Chapter.ID,
			Title: lesson.Chapter.Title,
			ChapterNumber: lesson.Chapter.ChapterNumber,
		},
	}

	c.JSON(http.StatusOK, response)
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

	var lesson models.Lessons

	// 1️⃣ Check if lesson exists
	if err := config.DB.First(&lesson, "id = ?", lessonID).Error; err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "Lesson not found",
		})
		return
	}

	// 2️⃣ Bind update payload
	var input models.LessonUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// 3️⃣ Apply updates safely
	if err := config.DB.Model(&lesson).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to update lesson",
		})
		return
	}

	// 4️⃣ Reload lesson with relations
	if err := config.DB.
		Preload("Chapter").
		Preload("Chapter.Course").
		First(&lesson, "id = ?", lessonID).Error; err != nil {

		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Error: "Failed to load updated lesson",
		})
		return
	}

	// 5️⃣ Map to LessonViewResponse DTO
	response := models.LessonViewResponse{
		ID:          lesson.ID,
		ChapterID:   lesson.ChapterID,
		CourseID:    lesson.CourseID,
		Title:       lesson.Title,
		ContentType: lesson.ContentType,
		ContentURL:  lesson.ContentURL,
		CreatedAt:   lesson.CreatedAt,
		UpdatedAt:   lesson.UpdatedAt,

		Course: models.CourseMiniResponse{
			ID:    lesson.Chapter.Course.ID,
			Title: lesson.Chapter.Course.Title,
		},

		Chapter: &models.ChapterMiniResponse{
			ID:            lesson.Chapter.ID,
			Title:         lesson.Chapter.Title,
			ChapterNumber: lesson.Chapter.ChapterNumber,
		},
	}

	// 6️⃣ Return clean DTO
	c.JSON(http.StatusOK, response)
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

	result := config.DB.Delete(&models.Lessons{}, "id = ?", lessonID)

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

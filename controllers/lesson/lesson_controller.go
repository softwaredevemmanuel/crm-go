// handlers/lesson_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"crm-go/dto"
	"crm-go/services/lesson"
)

type LessonHandler struct {
	lessonService *services.LessonService
}

func NewLessonHandler(lessonService *services.LessonService) *LessonHandler {
	return &LessonHandler{
		lessonService: lessonService,
	}
}

// CreateLesson handles the creation of a new lesson
// @Summary Create a lesson
// @Description Create a new lesson
// @Tags Lessons
// @Accept json
// @Produce json
// @Param request body dto.CreateLessonRequest true "Lesson request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons [post]
func (h *LessonHandler) CreateLesson(c *gin.Context) {
	var req dto.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	lesson, err := h.lessonService.CreateLesson(&req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lesson created successfully",
		"lesson":  lesson,
	})
}

// BulkCreateLessons handles bulk creation of lessons
// @Summary Bulk create lessons
// @Description Create multiple lessons
// @Tags Lessons
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateLessonsRequest true "Bulk lesson request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/bulk [post]
func (h *LessonHandler) BulkCreateLessons(c *gin.Context) {
	var req dto.BulkCreateLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.lessonService.BulkCreateLessons(&req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Lessons created successfully",
		"data":    result,
	})
}

// GetAllLessons handles fetching all lessons
// @Summary Get all lessons
// @Description Get a paginated list of all lessons
// @Tags Lessons
// @Accept json
// @Produce json
// @Param scheme_of_work_item_id query string false "Filter by scheme of work item ID"
// @Param class_id query string false "Filter by class ID"
// @Param arm_id query string false "Filter by arm ID"
// @Param status query string false "Filter by status"
// @Param week query int false "Filter by week"
// @Param search query string false "Search by title"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.LessonListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons [get]
func (h *LessonHandler) GetAllLessons(c *gin.Context) {
	var params dto.LessonQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.lessonService.GetAllLessons(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lessons retrieved successfully",
		"data":    response,
	})
}

// GetLessonByID handles fetching a single lesson by ID
// @Summary Get lesson by ID
// @Description Get a single lesson by its ID
// @Tags Lessons
// @Accept json
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/{id} [get]
func (h *LessonHandler) GetLessonByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson ID is required",
		})
		return
	}

	lesson, err := h.lessonService.GetLessonByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson retrieved successfully",
		"lesson":  lesson,
	})
}



// GetLessonsByClass handles fetching all lessons for a class
// @Summary Get lessons by class
// @Description Get all lessons for a specific class
// @Tags Lessons
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/class/{class_id} [get]
func (h *LessonHandler) GetLessonsByClass(c *gin.Context) {
	classID := c.Param("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Class ID is required",
		})
		return
	}

	lessons, err := h.lessonService.GetLessonsByClass(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lessons retrieved successfully",
		"lessons": lessons,
		"total":   len(lessons),
	})
}

// UpdateLesson handles updating a lesson
// @Summary Update a lesson
// @Description Update an existing lesson
// @Tags Lessons
// @Accept json
// @Produce json
// @Param id path string true "Lesson ID"
// @Param request body dto.UpdateLessonRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/{id} [put]
func (h *LessonHandler) UpdateLesson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson ID is required",
		})
		return
	}

	var req dto.UpdateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	lesson, err := h.lessonService.UpdateLesson(id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson updated successfully",
		"lesson":  lesson,
	})
}

// DeleteLesson handles deleting a lesson (soft delete)
// @Summary Delete a lesson
// @Description Soft delete a lesson
// @Tags Lessons
// @Accept json
// @Produce json
// @Param id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/{id} [delete]
func (h *LessonHandler) DeleteLesson(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson ID is required",
		})
		return
	}

	err := h.lessonService.DeleteLesson(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson deleted successfully",
	})
}
// handlers/lesson_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/lessons"
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
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user ID not found",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req dto.CreateLessonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	lesson, err := h.lessonService.CreateLesson(&req, userID)
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
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user ID not found",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req dto.BulkCreateLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.lessonService.BulkCreateLessons(&req, userID)
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
// @Param scheme_of_work_id query string false "Filter by scheme of work ID"
// @Param module_id query string false "Filter by module ID"
// @Param topic_id query string false "Filter by topic ID"
// @Param status query string false "Filter by status"
// @Param week query int false "Filter by week"
// @Param search query string false "Search by title or description"
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

// GetLessonsBySchemeOfWork handles fetching all lessons for a scheme of work
// @Summary Get lessons by scheme of work
// @Description Get all lessons for a specific scheme of work
// @Tags Lessons
// @Accept json
// @Produce json
// @Param scheme_of_work_id path string true "Scheme of Work ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/scheme/{scheme_of_work_id} [get]
func (h *LessonHandler) GetLessonsBySchemeOfWork(c *gin.Context) {
	schemeID := c.Param("scheme_of_work_id")
	if schemeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme of work ID is required",
		})
		return
	}

	lessons, err := h.lessonService.GetLessonsBySchemeOfWork(schemeID)
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

// GetLessonsByModule handles fetching all lessons for a module
// @Summary Get lessons by module
// @Description Get all lessons for a specific module
// @Tags Lessons
// @Accept json
// @Produce json
// @Param module_id path string true "Module ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/module/{module_id} [get]
func (h *LessonHandler) GetLessonsByModule(c *gin.Context) {
	moduleID := c.Param("module_id")
	if moduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Module ID is required",
		})
		return
	}

	lessons, err := h.lessonService.GetLessonsByModule(moduleID)
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

// GetLessonsByTopic handles fetching all lessons for a topic
// @Summary Get lessons by topic
// @Description Get all lessons for a specific topic
// @Tags Lessons
// @Accept json
// @Produce json
// @Param topic_id path string true "Topic ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/topic/{topic_id} [get]
func (h *LessonHandler) GetLessonsByTopic(c *gin.Context) {
	topicID := c.Param("topic_id")
	if topicID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Topic ID is required",
		})
		return
	}

	lessons, err := h.lessonService.GetLessonsByTopic(topicID)
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

// ReorderLessons handles reordering lessons
// @Summary Reorder lessons
// @Description Reorder lessons within a topic
// @Tags Lessons
// @Accept json
// @Produce json
// @Param request body dto.ReorderLessonsRequest true "Reorder request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lessons/reorder [put]
func (h *LessonHandler) ReorderLessons(c *gin.Context) {
	var req dto.ReorderLessonsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	if err := h.lessonService.ReorderLessons(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lessons reordered successfully",
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
package controllers

import (
	"net/http"
	"strconv"

	"crm-go/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetAllLessons retrieves all lessons with optional filtering
// @Summary Get all lessons
// @Description Get all lessons with optional filtering
// @Tags lessons
// @Accept json
// @Produce json
// @Param course_id query string false "Course ID"
// @Param module_id query string false "Module ID"
// @Param search query string false "Search term"
// @Param sort_by query string false "Sort field (title, order, created_at, updated_at)"
// @Param sort_order query string false "Sort order (asc, desc)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} models.PaginatedLessonsResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /lessons [get]
// @Security BearerAuth
func (ctl *LessonController) GetAllLessons(c *gin.Context) {
	// Parse query parameters
	var filters models.LessonFilters

	// Parse UUID parameters
	if courseIDStr := c.Query("course_id"); courseIDStr != "" {
		if courseID, err := uuid.Parse(courseIDStr); err == nil {
			filters.CourseID = courseID
		}
	}

	if moduleIDStr := c.Query("module_id"); moduleIDStr != "" {
		if moduleID, err := uuid.Parse(moduleIDStr); err == nil {
			filters.ModuleID = moduleID
		}
	}

	// Parse other parameters
	filters.Search = c.Query("search")
	filters.SortBy = c.Query("sort_by")
	filters.SortOrder = c.Query("sort_order")

	// Parse pagination parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filters.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	// Get tutor ID from context (assuming it's set by authentication middleware)
	if tutorID, exists := c.Get("tutor_id"); exists {
		if id, ok := tutorID.(uuid.UUID); ok {
			filters.TutorID = id
		}
	}

	// Get paginated results
	result, err := ctl.getLessonService.GetAllLessonsWithPagination(filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, result)
}

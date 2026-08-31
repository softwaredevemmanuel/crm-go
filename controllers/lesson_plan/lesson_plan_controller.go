// handlers/lesson_plan_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/lesson_plan"
)

type LessonPlanHandler struct {
	lessonPlanService *services.LessonPlanService
}

func NewLessonPlanHandler(lessonPlanService *services.LessonPlanService) *LessonPlanHandler {
	return &LessonPlanHandler{
		lessonPlanService: lessonPlanService,
	}
}

// CreateLessonPlan handles the creation of a new lesson plan
// @Summary Create a lesson plan
// @Description Create a new lesson plan for a lesson
// @Tags Lesson Plans
// @Accept json
// @Produce json
// @Param request body dto.CreateLessonPlanRequest true "Lesson plan request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-plans [post]
func (h *LessonPlanHandler) CreateLessonPlan(c *gin.Context) {
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

	var req dto.CreateLessonPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	lessonPlan, err := h.lessonPlanService.CreateLessonPlan(&req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
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
		"message":     "Lesson plan created successfully",
		"lesson_plan": lessonPlan,
	})
}

// GetAllLessonPlans handles fetching all lesson plans
// @Summary Get all lesson plans
// @Description Get a paginated list of all lesson plans
// @Tags Lesson Plans
// @Accept json
// @Produce json
// @Param lesson_id query string false "Filter by lesson ID"
// @Param search query string false "Search in lesson plan content"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.LessonPlanListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-plans [get]
func (h *LessonPlanHandler) GetAllLessonPlans(c *gin.Context) {
	var params dto.LessonPlanQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.lessonPlanService.GetAllLessonPlans(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Lesson plans retrieved successfully",
		"data":    response,
	})
}

// GetLessonPlanByID handles fetching a single lesson plan by ID
// @Summary Get lesson plan by ID
// @Description Get a single lesson plan by its ID
// @Tags Lesson Plans
// @Accept json
// @Produce json
// @Param id path string true "Lesson plan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-plans/{id} [get]
func (h *LessonPlanHandler) GetLessonPlanByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson plan ID is required",
		})
		return
	}

	lessonPlan, err := h.lessonPlanService.GetLessonPlanByID(id)
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
		"message":     "Lesson plan retrieved successfully",
		"lesson_plan": lessonPlan,
	})
}

// GetLessonPlanByLesson handles fetching a lesson plan by lesson ID
// @Summary Get lesson plan by lesson
// @Description Get a lesson plan by its associated lesson ID
// @Tags Lesson Plans
// @Accept json
// @Produce json
// @Param lesson_id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-plans/lesson/{lesson_id} [get]
func (h *LessonPlanHandler) GetLessonPlanByLesson(c *gin.Context) {
	lessonID := c.Param("lesson_id")
	if lessonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson ID is required",
		})
		return
	}

	lessonPlan, err := h.lessonPlanService.GetLessonPlanByLesson(lessonID)
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
		"message":     "Lesson plan retrieved successfully",
		"lesson_plan": lessonPlan,
	})
}

// UpdateLessonPlan handles updating a lesson plan
// @Summary Update a lesson plan
// @Description Update an existing lesson plan
// @Tags Lesson Plans
// @Accept json
// @Produce json
// @Param id path string true "Lesson plan ID"
// @Param request body dto.UpdateLessonPlanRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-plans/{id} [put]
func (h *LessonPlanHandler) UpdateLessonPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson plan ID is required",
		})
		return
	}

	var req dto.UpdateLessonPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	lessonPlan, err := h.lessonPlanService.UpdateLessonPlan(id, &req)
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
		"message":     "Lesson plan updated successfully",
		"lesson_plan": lessonPlan,
	})
}

// DeleteLessonPlan handles deleting a lesson plan (soft delete)
// @Summary Delete a lesson plan
// @Description Soft delete a lesson plan
// @Tags Lesson Plans
// @Accept json
// @Produce json
// @Param id path string true "Lesson plan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/lesson-plans/{id} [delete]
func (h *LessonPlanHandler) DeleteLessonPlan(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson plan ID is required",
		})
		return
	}

	err := h.lessonPlanService.DeleteLessonPlan(id)
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
		"message": "Lesson plan deleted successfully",
	})
}
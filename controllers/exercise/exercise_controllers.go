// handlers/exercise_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/exercise"
)

type ExerciseHandler struct {
	exerciseService *services.ExerciseService
}

func NewExerciseHandler(exerciseService *services.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{
		exerciseService: exerciseService,
	}
}

// CreateExercise handles the creation of a new exercise
// @Summary Create an exercise
// @Description Create a new exercise for a lesson
// @Tags Exercises
// @Accept json
// @Produce json
// @Param request body dto.CreateExerciseRequest true "Exercise request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exercises [post]
func (h *ExerciseHandler) CreateExercise(c *gin.Context) {
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

	var req dto.CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	exercise, err := h.exerciseService.CreateExercise(&req, userID)
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
		"message":  "Exercise created successfully",
		"exercise": exercise,
	})
}

// BulkCreateExercises handles bulk creation of exercises
// @Summary Bulk create exercises
// @Description Create multiple exercises for a lesson
// @Tags Exercises
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateExercisesRequest true "Bulk exercise request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exercises/bulk [post]
func (h *ExerciseHandler) BulkCreateExercises(c *gin.Context) {
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

	var req dto.BulkCreateExercisesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.exerciseService.BulkCreateExercises(&req, userID)
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
		"message": "Exercises created successfully",
		"data":    result,
	})
}

// GetAllExercises handles fetching all exercises
// @Summary Get all exercises
// @Description Get a paginated list of all exercises
// @Tags Exercises
// @Accept json
// @Produce json
// @Param lesson_id query string false "Filter by lesson ID"
// @Param search query string false "Search by title or content"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.ExerciseListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exercises [get]
func (h *ExerciseHandler) GetAllExercises(c *gin.Context) {
	var params dto.ExerciseQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.exerciseService.GetAllExercises(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Exercises retrieved successfully",
		"data":    response,
	})
}

// GetExerciseByID handles fetching a single exercise by ID
// @Summary Get exercise by ID
// @Description Get a single exercise by its ID
// @Tags Exercises
// @Accept json
// @Produce json
// @Param id path string true "Exercise ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exercises/{id} [get]
func (h *ExerciseHandler) GetExerciseByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exercise ID is required",
		})
		return
	}

	exercise, err := h.exerciseService.GetExerciseByID(id)
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
		"message":  "Exercise retrieved successfully",
		"exercise": exercise,
	})
}

// GetExercisesByLesson handles fetching all exercises for a lesson
// @Summary Get exercises by lesson
// @Description Get all exercises for a specific lesson
// @Tags Exercises
// @Accept json
// @Produce json
// @Param lesson_id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exercises/lesson/{lesson_id} [get]
func (h *ExerciseHandler) GetExercisesByLesson(c *gin.Context) {
	lessonID := c.Param("lesson_id")
	if lessonID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Lesson ID is required",
		})
		return
	}

	exercises, err := h.exerciseService.GetExercisesByLesson(lessonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Exercises retrieved successfully",
		"exercises": exercises,
		"total":    len(exercises),
	})
}

// UpdateExercise handles updating an exercise
// @Summary Update an exercise
// @Description Update an existing exercise
// @Tags Exercises
// @Accept json
// @Produce json
// @Param id path string true "Exercise ID"
// @Param request body dto.UpdateExerciseRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exercises/{id} [put]
func (h *ExerciseHandler) UpdateExercise(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exercise ID is required",
		})
		return
	}

	var req dto.UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	exercise, err := h.exerciseService.UpdateExercise(id, &req)
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
		"message":  "Exercise updated successfully",
		"exercise": exercise,
	})
}

// DeleteExercise handles deleting an exercise (soft delete)
// @Summary Delete an exercise
// @Description Soft delete an exercise
// @Tags Exercises
// @Accept json
// @Produce json
// @Param id path string true "Exercise ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/exercises/{id} [delete]
func (h *ExerciseHandler) DeleteExercise(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Exercise ID is required",
		})
		return
	}

	err := h.exerciseService.DeleteExercise(id)
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
		"message": "Exercise deleted successfully",
	})
}
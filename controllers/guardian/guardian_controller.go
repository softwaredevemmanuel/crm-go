package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/guardian"
)

type GuardianHandler struct {
	guardianService *services.GuardianService
}

func NewGuardianHandler(guardianService *services.GuardianService) *GuardianHandler {
	return &GuardianHandler{
		guardianService: guardianService,
	}
}

// CreateGuardian handles the creation of a new guardian
// @Summary Create a new guardian
// @Description Create a new guardian for a student
// @Tags Guardians
// @Accept json
// @Produce json
// @Param request body dto.CreateGuardianRequest true "Guardian creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/guardians [post]
func (h *GuardianHandler) CreateGuardian(c *gin.Context) {
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

	var req dto.CreateGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	guardian, err := h.guardianService.CreateGuardian(&req, userID)
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
		"message":  "Guardian created successfully",
		"guardian": guardian,
	})
}

// GetAllGuardians handles fetching all guardians with pagination and filters
// @Summary Get all guardians
// @Description Get a paginated list of guardians with optional filtering
// @Tags Guardians
// @Accept json
// @Produce json
// @Param search query string false "Search by name, email, or phone"
// @Param student_id query string false "Filter by student ID"
// @Param relationship query string false "Filter by relationship"
// @Param status query string false "Filter by status"
// @Param is_primary query bool false "Filter by primary status"
// @Param is_emergency query bool false "Filter by emergency status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.GuardianListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/guardians [get]
func (h *GuardianHandler) GetAllGuardians(c *gin.Context) {
	var params dto.GuardianQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.guardianService.GetAllGuardians(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Guardians retrieved successfully",
		"data":    response,
	})
}

// GetGuardianByID handles fetching a single guardian by ID
// @Summary Get guardian by ID
// @Description Get a single guardian by its ID
// @Tags Guardians
// @Accept json
// @Produce json
// @Param id path string true "Guardian ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/guardians/{id} [get]
func (h *GuardianHandler) GetGuardianByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Guardian ID is required",
		})
		return
	}

	guardian, err := h.guardianService.GetGuardianByID(id)
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
		"message":  "Guardian retrieved successfully",
		"guardian": guardian,
	})
}

// GetGuardiansByStudent handles fetching all guardians for a student
// @Summary Get guardians by student
// @Description Get all guardians associated with a specific student
// @Tags Guardians
// @Accept json
// @Produce json
// @Param student_id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/guardians/student/{student_id} [get]
func (h *GuardianHandler) GetGuardiansByStudent(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Student ID is required",
		})
		return
	}

	guardians, err := h.guardianService.GetGuardiansByStudent(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Guardians retrieved successfully",
		"guardians": guardians,
	})
}

// UpdateGuardian handles updating an existing guardian
// @Summary Update a guardian
// @Description Update an existing guardian by ID
// @Tags Guardians
// @Accept json
// @Produce json
// @Param id path string true "Guardian ID"
// @Param request body dto.UpdateGuardianRequest true "Guardian update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/guardians/{id} [put]
func (h *GuardianHandler) UpdateGuardian(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Guardian ID is required",
		})
		return
	}

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

	var req dto.UpdateGuardianRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	guardian, err := h.guardianService.UpdateGuardian(id, &req, userID)
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

	c.JSON(http.StatusOK, gin.H{
		"message":  "Guardian updated successfully",
		"guardian": guardian,
	})
}

// DeleteGuardian handles deleting a guardian (soft delete)
// @Summary Delete a guardian
// @Description Soft delete a guardian by ID
// @Tags Guardians
// @Accept json
// @Produce json
// @Param id path string true "Guardian ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/guardians/{id} [delete]
func (h *GuardianHandler) DeleteGuardian(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Guardian ID is required",
		})
		return
	}

	err := h.guardianService.DeleteGuardian(id)
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
		"message": "Guardian deleted successfully",
	})
}
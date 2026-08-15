package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/arm"
)

type ArmHandler struct {
	armService *services.ArmService
}

func NewArmHandler(armService *services.ArmService) *ArmHandler {
	return &ArmHandler{
		armService: armService,
	}
}

// CreateArm handles the creation of a new arm
// @Summary Create a new arm
// @Description Create a new arm (sub-grade) with the provided details
// @Tags Arms
// @Accept json
// @Produce json
// @Param request body dto.CreateArmRequest true "Arm creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arms [post]
func (h *ArmHandler) CreateArm(c *gin.Context) {
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

	var req dto.CreateArmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	arm, err := h.armService.CreateArm(&req, userID)
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
		"message": "Arm created successfully",
		"arm":     arm,
	})
}

// GetAllArms handles fetching all arms with pagination and filters
// @Summary Get all arms
// @Description Get a paginated list of arms with optional filtering
// @Tags Arms
// @Accept json
// @Produce json
// @Param search query string false "Search by name, code, or description"
// @Param grade_id query string false "Filter by grade ID"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order" default(desc)
// @Success 200 {object} dto.ArmListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arms [get]
func (h *ArmHandler) GetAllArms(c *gin.Context) {
	var params dto.ArmQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.armService.GetAllArms(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Arms retrieved successfully",
		"data":    response,
	})
}

// GetArmByID handles fetching a single arm by ID
// @Summary Get arm by ID
// @Description Get a single arm by its ID
// @Tags Arms
// @Accept json
// @Produce json
// @Param id path string true "Arm ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arms/{id} [get]
func (h *ArmHandler) GetArmByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Arm ID is required",
		})
		return
	}

	arm, err := h.armService.GetArmByID(id)
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
		"message": "Arm retrieved successfully",
		"arm":     arm,
	})
}

// GetArmsByGrade handles fetching all arms for a grade
// @Summary Get arms by grade
// @Description Get all arms associated with a specific grade
// @Tags Arms
// @Accept json
// @Produce json
// @Param grade_id path string true "Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arms/grade/{grade_id} [get]
func (h *ArmHandler) GetArmsByGrade(c *gin.Context) {
	gradeID := c.Param("grade_id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade ID is required",
		})
		return
	}

	arms, err := h.armService.GetArmsByGrade(gradeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Arms retrieved successfully",
		"arms":    arms,
	})
}

// UpdateArm handles updating an existing arm
// @Summary Update an arm
// @Description Update an existing arm by ID
// @Tags Arms
// @Accept json
// @Produce json
// @Param id path string true "Arm ID"
// @Param request body dto.UpdateArmRequest true "Arm update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arms/{id} [put]
func (h *ArmHandler) UpdateArm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Arm ID is required",
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

	var req dto.UpdateArmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	arm, err := h.armService.UpdateArm(id, &req, userID)
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
		"message": "Arm updated successfully",
		"arm":     arm,
	})
}

// DeleteArm handles deleting an arm (soft delete)
// @Summary Delete an arm
// @Description Soft delete an arm by ID
// @Tags Arms
// @Accept json
// @Produce json
// @Param id path string true "Arm ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arms/{id} [delete]
func (h *ArmHandler) DeleteArm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Arm ID is required",
		})
		return
	}

	err := h.armService.DeleteArm(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "has students assigned") {
			c.JSON(http.StatusConflict, gin.H{
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
		"message": "Arm deleted successfully",
	})
}
// DeleteArm handles deleting an arm (permanent delete)
// @Summary Delete an arm
// @Description Permanent delete an arm by ID
// @Tags Arms
// @Accept json
// @Produce json
// @Param id path string true "Arm ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arms/permanent/{id} [delete]
func (h *ArmHandler) DeleteArmPermanently(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Arm ID is required",
		})
		return
	}

	err := h.armService.DeleteArmPermanently(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "has students assigned") {
			c.JSON(http.StatusConflict, gin.H{
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
		"message": "Arm deleted successfully",
	})
}
// handlers/academic_session_handler.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/academic_session"
)

type AcademicSessionHandler struct {
	sessionService *services.AcademicSessionService
}

func NewAcademicSessionHandler(sessionService *services.AcademicSessionService) *AcademicSessionHandler {
	return &AcademicSessionHandler{
		sessionService: sessionService,
	}
}

// CreateAcademicSession handles the creation of a new academic session
// @Summary Create a new academic session
// @Description Create a new academic session with the provided details
// @Tags Academic Sessions
// @Accept json
// @Produce json
// @Param request body dto.CreateAcademicSessionRequest true "Academic session creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/academic-sessions [post]
func (h *AcademicSessionHandler) CreateAcademicSession(c *gin.Context) {
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

	var req dto.CreateAcademicSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	session, err := h.sessionService.CreateAcademicSession(&req, userID)
	if err != nil {
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
		"message": "Academic session created successfully",
		"session": session,
	})
}

// GetAllAcademicSessions handles fetching all academic sessions with pagination and filters
// @Summary Get all academic sessions
// @Description Get a paginated list of academic sessions with optional filtering
// @Tags Academic Sessions
// @Accept json
// @Produce json
// @Param search query string false "Search by year, code, or description"
// @Param status query string false "Filter by status"
// @Param is_current query bool false "Filter by current status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.AcademicSessionListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/academic-sessions [get]
func (h *AcademicSessionHandler) GetAllAcademicSessions(c *gin.Context) {
	var params dto.AcademicSessionQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.sessionService.GetAllAcademicSessions(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Academic sessions retrieved successfully",
		"data":    response,
	})
}

// GetAcademicSessionByID handles fetching a single academic session by ID
// @Summary Get academic session by ID
// @Description Get a single academic session by its ID
// @Tags Academic Sessions
// @Accept json
// @Produce json
// @Param id path string true "Academic Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/academic-sessions/{id} [get]
func (h *AcademicSessionHandler) GetAcademicSessionByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Academic session ID is required",
		})
		return
	}

	session, err := h.sessionService.GetAcademicSessionByID(id)
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
		"message": "Academic session retrieved successfully",
		"session": session,
	})
}

// GetCurrentAcademicSession handles fetching the current academic session
// @Summary Get current academic session
// @Description Get the currently active academic session
// @Tags Academic Sessions
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/academic-sessions/current [get]
func (h *AcademicSessionHandler) GetCurrentAcademicSession(c *gin.Context) {
	session, err := h.sessionService.GetCurrentAcademicSession()
	if err != nil {
		if strings.Contains(err.Error(), "no current") {
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
		"message": "Current academic session retrieved successfully",
		"session": session,
	})
}

// UpdateAcademicSession handles updating an existing academic session
// @Summary Update an academic session
// @Description Update an existing academic session by ID
// @Tags Academic Sessions
// @Accept json
// @Produce json
// @Param id path string true "Academic Session ID"
// @Param request body dto.UpdateAcademicSessionRequest true "Academic session update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/academic-sessions/{id} [put]
func (h *AcademicSessionHandler) UpdateAcademicSession(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Academic session ID is required",
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

	var req dto.UpdateAcademicSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	session, err := h.sessionService.UpdateAcademicSession(id, &req, userID)
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
		"message": "Academic session updated successfully",
		"session": session,
	})
}

// DeleteAcademicSession handles deleting an academic session (soft delete)
// @Summary Delete an academic session
// @Description Soft delete an academic session by ID
// @Tags Academic Sessions
// @Accept json
// @Produce json
// @Param id path string true "Academic Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/academic-sessions/{id} [delete]
func (h *AcademicSessionHandler) DeleteAcademicSession(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Academic session ID is required",
		})
		return
	}

	err := h.sessionService.DeleteAcademicSession(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "cannot delete") {
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
		"message": "Academic session deleted successfully",
	})
}
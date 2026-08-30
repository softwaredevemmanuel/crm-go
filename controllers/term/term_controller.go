// handlers/term_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/term"
)

type TermHandler struct {
	termService *services.TermService
}

func NewTermHandler(termService *services.TermService) *TermHandler {
	return &TermHandler{
		termService: termService,
	}
}

// CreateTerm handles the creation of a new term
// @Summary Create a term
// @Description Create a new term for an academic session
// @Tags Terms
// @Accept json
// @Produce json
// @Param request body dto.CreateTermRequest true "Term request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms [post]
func (h *TermHandler) CreateTerm(c *gin.Context) {
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

	var req dto.CreateTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	term, err := h.termService.CreateTerm(&req, userID)
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
		"message": "Term created successfully",
		"term":    term,
	})
}


// GetAllTerms handles fetching all terms
// @Summary Get all terms
// @Description Get a paginated list of all terms
// @Tags Terms
// @Accept json
// @Produce json
// @Param academic_session_id query string false "Filter by academic session ID"
// @Param status query string false "Filter by status"
// @Param is_current query bool false "Filter by current status"
// @Param term_number query int false "Filter by term number"
// @Param search query string false "Search by name, code, or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.TermListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms [get]
func (h *TermHandler) GetAllTerms(c *gin.Context) {
	var params dto.TermQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.termService.GetAllTerms(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Terms retrieved successfully",
		"data":    response,
	})
}

// GetTermByID handles fetching a single term by ID
// @Summary Get term by ID
// @Description Get a single term by its ID
// @Tags Terms
// @Accept json
// @Produce json
// @Param id path string true "Term ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms/{id} [get]
func (h *TermHandler) GetTermByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Term ID is required",
		})
		return
	}

	term, err := h.termService.GetTermByID(id)
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
		"message": "Term retrieved successfully",
		"term":    term,
	})
}

// GetTermsByAcademicSession handles fetching all terms for an academic session
// @Summary Get terms by academic session
// @Description Get all terms for a specific academic session
// @Tags Terms
// @Accept json
// @Produce json
// @Param session_id path string true "Academic Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms/session/{session_id} [get]
func (h *TermHandler) GetTermsByAcademicSession(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Academic session ID is required",
		})
		return
	}

	terms, err := h.termService.GetTermsByAcademicSession(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Terms retrieved successfully",
		"terms":   terms,
		"total":   len(terms),
	})
}

// GetCurrentTerm handles fetching the current term for an academic session
// @Summary Get current term
// @Description Get the current term for an academic session
// @Tags Terms
// @Accept json
// @Produce json
// @Param session_id path string true "Academic Session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms/session/{session_id}/current [get]
func (h *TermHandler) GetCurrentTerm(c *gin.Context) {
	sessionID := c.Param("session_id")
	if sessionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Academic session ID is required",
		})
		return
	}

	term, err := h.termService.GetCurrentTerm(sessionID)
	if err != nil {
		if strings.Contains(err.Error(), "no current term") {
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
		"message": "Current term retrieved successfully",
		"term":    term,
	})
}

// GetTermStats handles fetching term statistics
// @Summary Get term statistics
// @Description Get statistics for terms
// @Tags Terms
// @Accept json
// @Produce json
// @Param academic_session_id query string false "Filter by academic session ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms/stats [get]
func (h *TermHandler) GetTermStats(c *gin.Context) {
	filter := make(map[string]interface{})
	
	if sessionID := c.Query("academic_session_id"); sessionID != "" {
		filter["academic_session_id"] = sessionID
	}

	stats, err := h.termService.GetTermStats(filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Statistics retrieved successfully",
		"data":    stats,
	})
}

// UpdateTerm handles updating a term
// @Summary Update a term
// @Description Update an existing term
// @Tags Terms
// @Accept json
// @Produce json
// @Param id path string true "Term ID"
// @Param request body dto.UpdateTermRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms/{id} [put]
func (h *TermHandler) UpdateTerm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Term ID is required",
		})
		return
	}

	var req dto.UpdateTermRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	term, err := h.termService.UpdateTerm(id, &req)
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
		"message": "Term updated successfully",
		"term":    term,
	})
}

// DeleteTerm handles deleting a term (soft delete)
// @Summary Delete a term
// @Description Soft delete a term
// @Tags Terms
// @Accept json
// @Produce json
// @Param id path string true "Term ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/terms/{id} [delete]
func (h *TermHandler) DeleteTerm(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Term ID is required",
		})
		return
	}

	err := h.termService.DeleteTerm(id)
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
		"message": "Term deleted successfully",
	})
}
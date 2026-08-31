// handlers/scheme_of_work_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/scheme_of_work"
)

type SchemeOfWorkHandler struct {
	schemeService *services.SchemeOfWorkService
}

func NewSchemeOfWorkHandler(schemeService *services.SchemeOfWorkService) *SchemeOfWorkHandler {
	return &SchemeOfWorkHandler{
		schemeService: schemeService,
	}
}

// CreateSchemeOfWork handles the creation of a new scheme of work
// @Summary Create a scheme of work
// @Description Create a new scheme of work
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param request body dto.CreateSchemeOfWorkRequest true "Scheme request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes [post]
func (h *SchemeOfWorkHandler) CreateSchemeOfWork(c *gin.Context) {
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

	var req dto.CreateSchemeOfWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	scheme, err := h.schemeService.CreateSchemeOfWork(&req, userID)
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
		"message": "Scheme of work created successfully",
		"scheme":  scheme,
	})
}

// BulkCreateSchemes handles bulk creation of schemes
// @Summary Bulk create schemes
// @Description Create multiple schemes of work
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateSchemesRequest true "Bulk scheme request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/bulk [post]
func (h *SchemeOfWorkHandler) BulkCreateSchemes(c *gin.Context) {
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

	var req dto.BulkCreateSchemesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.schemeService.BulkCreateSchemes(&req, userID)
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
		"message": "Schemes created successfully",
		"data":    result,
	})
}

// GetAllSchemes handles fetching all schemes
// @Summary Get all schemes
// @Description Get a paginated list of all schemes of work
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param academic_session_id query string false "Filter by academic session ID"
// @Param term_id query string false "Filter by term ID"
// @Param subject_id query string false "Filter by subject ID"
// @Param class_id query string false "Filter by class ID"
// @Param status query string false "Filter by status"
// @Param search query string false "Search by title or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.SchemeOfWorkListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes [get]
func (h *SchemeOfWorkHandler) GetAllSchemes(c *gin.Context) {
	var params dto.SchemeOfWorkQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.schemeService.GetAllSchemes(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Schemes retrieved successfully",
		"data":    response,
	})
}

// GetSchemeByID handles fetching a single scheme by ID
// @Summary Get scheme by ID
// @Description Get a single scheme of work by its ID
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param id path string true "Scheme ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/{id} [get]
func (h *SchemeOfWorkHandler) GetSchemeByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme ID is required",
		})
		return
	}

	scheme, err := h.schemeService.GetSchemeByID(id)
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
		"message": "Scheme retrieved successfully",
		"scheme":  scheme,
	})
}

// GetSchemesBySubject handles fetching all schemes for a subject
// @Summary Get schemes by subject
// @Description Get all schemes of work for a specific subject
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param subject_id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/subject/{subject_id} [get]
func (h *SchemeOfWorkHandler) GetSchemesBySubject(c *gin.Context) {
	subjectID := c.Param("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	schemes, err := h.schemeService.GetSchemesBySubject(subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Schemes retrieved successfully",
		"schemes": schemes,
		"total":   len(schemes),
	})
}

// GetSchemesByClass handles fetching all schemes for a class
// @Summary Get schemes by class
// @Description Get all schemes of work for a specific class
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param class_id path string true "Class ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/class/{class_id} [get]
func (h *SchemeOfWorkHandler) GetSchemesByClass(c *gin.Context) {
	classID := c.Param("class_id")
	if classID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Class ID is required",
		})
		return
	}

	schemes, err := h.schemeService.GetSchemesByClass(classID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Schemes retrieved successfully",
		"schemes": schemes,
		"total":   len(schemes),
	})
}

// UpdateSchemeOfWork handles updating a scheme
// @Summary Update a scheme
// @Description Update an existing scheme of work
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param id path string true "Scheme ID"
// @Param request body dto.UpdateSchemeOfWorkRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/{id} [put]
func (h *SchemeOfWorkHandler) UpdateSchemeOfWork(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme ID is required",
		})
		return
	}

	var req dto.UpdateSchemeOfWorkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	scheme, err := h.schemeService.UpdateSchemeOfWork(id, &req)
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
		"message": "Scheme updated successfully",
		"scheme":  scheme,
	})
}

// DeleteSchemeOfWork handles deleting a scheme (soft delete)
// @Summary Delete a scheme
// @Description Soft delete a scheme of work
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param id path string true "Scheme ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/{id} [delete]
func (h *SchemeOfWorkHandler) DeleteSchemeOfWork(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Scheme ID is required",
		})
		return
	}

	err := h.schemeService.DeleteSchemeOfWork(id)
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
		"message": "Scheme deleted successfully",
	})
}
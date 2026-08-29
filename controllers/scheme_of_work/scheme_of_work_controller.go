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
// @Description Create a new scheme of work for a subject and grade
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

// BulkCreateSchemes handles bulk creation of schemes of work
// @Summary Bulk create schemes of work
// @Description Create multiple schemes of work for a subject and grade
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateSchemeOfWorkRequest true "Bulk scheme request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
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

	var req dto.BulkCreateSchemeOfWorkRequest
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

// GetAllSchemes handles fetching all schemes of work
// @Summary Get all schemes of work
// @Description Get a paginated list of schemes of work with optional filtering
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param subject_id query string false "Filter by subject ID"
// @Param grade query string false "Filter by grade"
// @Param term query string false "Filter by term"
// @Param status query string false "Filter by status"
// @Param week query int false "Filter by week"
// @Param search query string false "Search by topic or objectives"
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

// GetSchemesByGrade handles fetching all schemes for a grade
// @Summary Get schemes by grade
// @Description Get all schemes of work for a specific grade
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param grade path string true "Grade (e.g., JSS1, SS2, JSS2)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/grade/{grade} [get]
func (h *SchemeOfWorkHandler) GetSchemesByGrade(c *gin.Context) {
	grade := c.Param("grade")
	if grade == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade is required",
		})
		return
	}

	schemes, err := h.schemeService.GetSchemesByGrade(grade)
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

// GetSchemesByGradeAndTerm handles fetching all schemes for a grade and term
// @Summary Get schemes by grade and term
// @Description Get all schemes of work for a specific grade and term
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param grade path string true "Grade (e.g., JSS1, SS2)"
// @Param term path string true "Term (first, second, third)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/grade/{grade}/term/{term} [get]
func (h *SchemeOfWorkHandler) GetSchemesByGradeAndTerm(c *gin.Context) {
	grade := c.Param("grade")
	term := c.Param("term")

	if grade == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade is required",
		})
		return
	}
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Term is required",
		})
		return
	}

	schemes, err := h.schemeService.GetSchemesByGradeAndTerm(grade, term)
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

// GetSchemesByTeacher handles fetching all schemes for subjects taught by a teacher
// @Summary Get schemes by teacher
// @Description Get all schemes of work for subjects taught by a specific teacher
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param teacher_id path string true "Teacher ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/teacher/{teacher_id} [get]
func (h *SchemeOfWorkHandler) GetSchemesByTeacher(c *gin.Context) {
	teacherID := c.Param("teacher_id")
	if teacherID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Teacher ID is required",
		})
		return
	}

	schemes, err := h.schemeService.GetSchemesByTeacher(teacherID)
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

// GetSchemeStats handles fetching scheme statistics
// @Summary Get scheme statistics
// @Description Get statistics for schemes of work
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param subject_id query string false "Filter by subject ID"
// @Param grade query string false "Filter by grade"
// @Param term query string false "Filter by term"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/stats [get]
func (h *SchemeOfWorkHandler) GetSchemeStats(c *gin.Context) {
	filter := make(map[string]interface{})
	
	if subjectID := c.Query("subject_id"); subjectID != "" {
		filter["subject_id"] = subjectID
	}
	if grade := c.Query("grade"); grade != "" {
		filter["grade"] = grade
	}
	if term := c.Query("term"); term != "" {
		filter["term"] = term
	}

	stats, err := h.schemeService.GetSchemeStats(filter)
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

// GetSchemeOverview handles fetching scheme overview
// @Summary Get scheme overview
// @Description Get an overview of schemes for a subject, grade, and term
// @Tags Scheme of Work
// @Accept json
// @Produce json
// @Param subject_id path string true "Subject ID"
// @Param grade path string true "Grade"
// @Param term path string true "Term"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/schemes/overview/{subject_id}/{grade}/{term} [get]
func (h *SchemeOfWorkHandler) GetSchemeOverview(c *gin.Context) {
	subjectID := c.Param("subject_id")
	grade := c.Param("grade")
	term := c.Param("term")

	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}
	if grade == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade is required",
		})
		return
	}
	if term == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Term is required",
		})
		return
	}

	overview, err := h.schemeService.GetSchemeOverview(subjectID, grade, term)
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
		"message": "Overview retrieved successfully",
		"data":    overview,
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
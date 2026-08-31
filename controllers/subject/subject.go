// handlers/subject_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/subjects"
)

type SubjectHandler struct {
	subjectService *services.SubjectService
}

func NewSubjectHandler(subjectService *services.SubjectService) *SubjectHandler {
	return &SubjectHandler{
		subjectService: subjectService,
	}
}

// CreateSubject handles the creation of a new subject
// @Summary Create a subject
// @Description Create a new subject
// @Tags Subjects
// @Accept json
// @Produce json
// @Param request body dto.CreateSubjectRequest true "Subject request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
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

	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	subject, err := h.subjectService.CreateSubject(&req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
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
		"message": "Subject created successfully",
		"subject": subject,
	})
}

// BulkCreateSubjects handles bulk creation of subjects
// @Summary Bulk create subjects
// @Description Create multiple subjects
// @Tags Subjects
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateSubjectsRequest true "Bulk subject request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/bulk [post]
func (h *SubjectHandler) BulkCreateSubjects(c *gin.Context) {
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

	var req dto.BulkCreateSubjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.subjectService.BulkCreateSubjects(&req, userID)
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
		"message": "Subjects created successfully",
		"data":    result,
	})
}

// GetAllSubjects handles fetching all subjects
// @Summary Get all subjects
// @Description Get a paginated list of all subjects
// @Tags Subjects
// @Accept json
// @Produce json
// @Param department_id query string false "Filter by department ID"
// @Param status query string false "Filter by status"
// @Param search query string false "Search by name, code, or description"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.SubjectListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects [get]
func (h *SubjectHandler) GetAllSubjects(c *gin.Context) {
	var params dto.SubjectQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.subjectService.GetAllSubjects(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subjects retrieved successfully",
		"data":    response,
	})
}

// GetSubjectByID handles fetching a single subject by ID
// @Summary Get subject by ID
// @Description Get a single subject by its ID
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/{id} [get]
func (h *SubjectHandler) GetSubjectByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	subject, err := h.subjectService.GetSubjectByID(id)
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
		"message": "Subject retrieved successfully",
		"subject": subject,
	})
}

// GetSubjectsByDepartment handles fetching all subjects for a department
// @Summary Get subjects by department
// @Description Get all subjects for a specific department
// @Tags Subjects
// @Accept json
// @Produce json
// @Param department_id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/department/{department_id} [get]
func (h *SubjectHandler) GetSubjectsByDepartment(c *gin.Context) {
	departmentID := c.Param("department_id")
	if departmentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Department ID is required",
		})
		return
	}

	subjects, err := h.subjectService.GetSubjectsByDepartment(departmentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Subjects retrieved successfully",
		"subjects": subjects,
		"total":    len(subjects),
	})
}

// GetActiveSubjects handles fetching all active subjects
// @Summary Get active subjects
// @Description Get all active subjects
// @Tags Subjects
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/active [get]
func (h *SubjectHandler) GetActiveSubjects(c *gin.Context) {
	subjects, err := h.subjectService.GetActiveSubjects()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Active subjects retrieved successfully",
		"subjects": subjects,
		"total":    len(subjects),
	})
}

// GetSubjectStats handles fetching subject statistics
// @Summary Get subject statistics
// @Description Get statistics for subjects
// @Tags Subjects
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/stats [get]
func (h *SubjectHandler) GetSubjectStats(c *gin.Context) {
	stats, err := h.subjectService.GetSubjectStats()
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

// UpdateSubject handles updating a subject
// @Summary Update a subject
// @Description Update an existing subject
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Param request body dto.UpdateSubjectRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/{id} [put]
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	var req dto.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	subject, err := h.subjectService.UpdateSubject(id, &req)
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
		"message": "Subject updated successfully",
		"subject": subject,
	})
}

// DeleteSubject handles deleting a subject (soft delete)
// @Summary Delete a subject
// @Description Soft delete a subject
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/{id} [delete]
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	err := h.subjectService.DeleteSubject(id)
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
		"message": "Subject deleted successfully",
	})
}
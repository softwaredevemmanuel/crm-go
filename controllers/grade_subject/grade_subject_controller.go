package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/grade_subject"
)

type GradeSubjectHandler struct {
	gradeSubjectService *services.GradeSubjectService
}

func NewGradeSubjectHandler(gradeSubjectService *services.GradeSubjectService) *GradeSubjectHandler {
	return &GradeSubjectHandler{
		gradeSubjectService: gradeSubjectService,
	}
}

// CreateGradeSubject handles the creation of a new grade-subject mapping
// @Summary Create a grade-subject mapping
// @Description Create a new mapping between a grade and a subject
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param request body dto.CreateGradeSubjectRequest true "Grade-subject mapping request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects [post]
func (h *GradeSubjectHandler) CreateGradeSubject(c *gin.Context) {
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

	var req dto.CreateGradeSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	gradeSubject, err := h.gradeSubjectService.CreateGradeSubject(&req, userID)
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
		"message":      "Grade-subject mapping created successfully",
		"grade_subject": gradeSubject,
	})
}

// BulkCreateGradeSubjects handles bulk creation of grade-subject mappings
// @Summary Bulk create grade-subject mappings
// @Description Create multiple subject mappings for a grade at once
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateGradeSubjectRequest true "Bulk creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects/bulk [post]
func (h *GradeSubjectHandler) BulkCreateGradeSubjects(c *gin.Context) {
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

	var req dto.BulkCreateGradeSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.gradeSubjectService.BulkCreateGradeSubjects(&req, userID)
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
		"message": "Grade-subject mappings created successfully",
		"data":    result,
	})
}

// GetAllGradeSubjects handles fetching all grade-subject mappings
// @Summary Get all grade-subject mappings
// @Description Get a paginated list of grade-subject mappings with optional filtering
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param grade_id query string false "Filter by grade ID"
// @Param subject_id query string false "Filter by subject ID"
// @Param status query string false "Filter by status"
// @Param is_compulsory query bool false "Filter by compulsory status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.GradeSubjectListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects [get]
func (h *GradeSubjectHandler) GetAllGradeSubjects(c *gin.Context) {
	var params dto.GradeSubjectQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.gradeSubjectService.GetAllGradeSubjects(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Grade-subject mappings retrieved successfully",
		"data":    response,
	})
}

// GetGradeSubjectByID handles fetching a single grade-subject mapping by ID
// @Summary Get grade-subject mapping by ID
// @Description Get a single grade-subject mapping by its ID
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param id path string true "Grade-Subject Mapping ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects/{id} [get]
func (h *GradeSubjectHandler) GetGradeSubjectByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade-subject mapping ID is required",
		})
		return
	}

	gradeSubject, err := h.gradeSubjectService.GetGradeSubjectByID(id)
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
		"message":      "Grade-subject mapping retrieved successfully",
		"grade_subject": gradeSubject,
	})
}

// GetSubjectsByGrade handles fetching all subjects for a grade
// @Summary Get subjects by grade
// @Description Get all subjects mapped to a specific grade
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param grade_id path string true "Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects/grade/{grade_id} [get]
func (h *GradeSubjectHandler) GetSubjectsByGrade(c *gin.Context) {
	gradeID := c.Param("grade_id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade ID is required",
		})
		return
	}

	gradeSubjects, err := h.gradeSubjectService.GetSubjectsByGrade(gradeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Subjects retrieved successfully",
		"grade_subjects": gradeSubjects,
	})
}

// GetGradesBySubject handles fetching all grades for a subject
// @Summary Get grades by subject
// @Description Get all grades mapped to a specific subject
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param subject_id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects/subject/{subject_id} [get]
func (h *GradeSubjectHandler) GetGradesBySubject(c *gin.Context) {
	subjectID := c.Param("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	gradeSubjects, err := h.gradeSubjectService.GetGradesBySubject(subjectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Grades retrieved successfully",
		"grade_subjects": gradeSubjects,
	})
}

// UpdateGradeSubject handles updating an existing grade-subject mapping
// @Summary Update a grade-subject mapping
// @Description Update an existing grade-subject mapping by ID
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param id path string true "Grade-Subject Mapping ID"
// @Param request body dto.UpdateGradeSubjectRequest true "Grade-subject update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects/{id} [put]
func (h *GradeSubjectHandler) UpdateGradeSubject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade-subject mapping ID is required",
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

	var req dto.UpdateGradeSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	gradeSubject, err := h.gradeSubjectService.UpdateGradeSubject(id, &req, userID)
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
		"message":      "Grade-subject mapping updated successfully",
		"grade_subject": gradeSubject,
	})
}

// DeleteGradeSubject handles deleting a grade-subject mapping (soft delete)
// @Summary Delete a grade-subject mapping
// @Description Soft delete a grade-subject mapping by ID
// @Tags Grade Subjects
// @Accept json
// @Produce json
// @Param id path string true "Grade-Subject Mapping ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/grade-subjects/{id} [delete]
func (h *GradeSubjectHandler) DeleteGradeSubject(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade-subject mapping ID is required",
		})
		return
	}

	err := h.gradeSubjectService.DeleteGradeSubject(id)
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
		"message": "Grade-subject mapping deleted successfully",
	})
}
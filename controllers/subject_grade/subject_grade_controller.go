package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/subject_grade"
	"gorm.io/gorm"

)

type SubjectGradeController struct {
		db  *gorm.DB
	subjectGradeService *services.SubjectGradeService
}

func NewSubjectGradeController(db *gorm.DB, subjectGradeService *services.SubjectGradeService) *SubjectGradeController {
	return &SubjectGradeController{
				db: 			db,
		subjectGradeService: subjectGradeService,
	}
}

// CreateSubjectGrade handles the creation of a new subject-grade relationship
// @Summary Create a new subject-grade relationship
// @Description Create a relationship between a subject and a class grade
// @Tags Subject Grades
// @Accept json
// @Produce json
// @Param request body dto.CreateSubjectGradeRequest true "Subject-grade creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subject-grades [post]
func (h *SubjectGradeController) CreateSubjectGrade(c *gin.Context) {
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

	var req dto.CreateSubjectGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	subjectGrade, err := h.subjectGradeService.CreateSubjectGrade(&req, userID)
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
		"message":       "Subject-grade relationship created successfully",
		"subject_grade": subjectGrade,
	})
}

// BulkCreateSubjectGrades handles bulk creation of subject-grade relationships
// @Summary Bulk create subject-grade relationships
// @Description Create multiple subject-grade relationships for a grade
// @Tags Subject Grades
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateSubjectGradeRequest true "Bulk creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subject-grades/bulk [post]
func (h *SubjectGradeController) BulkCreateSubjectGrades(c *gin.Context) {
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

	var req dto.BulkCreateSubjectGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	response, err := h.subjectGradeService.BulkCreateSubjectGrades(&req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Subject-grade relationships created successfully",
		"data":    response,
	})
}

// GetAllSubjectGrades handles fetching all subject-grade relationships
// @Summary Get all subject-grade relationships
// @Description Get a paginated list of subject-grade relationships
// @Tags Subject Grades
// @Accept json
// @Produce json
// @Param subject_id query string false "Filter by subject ID"
// @Param grade_id query string false "Filter by grade ID"
// @Param academic_year query string false "Filter by academic year"
// @Param status query string false "Filter by status"
// @Param is_required query bool false "Filter by required status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.SubjectGradeListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subject-grades/relationship [get]
func (h *SubjectGradeController) GetAllSubjectGrades(c *gin.Context) {
	var params dto.SubjectGradeQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.subjectGradeService.GetAllSubjectGrades(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subject-grade relationships retrieved successfully",
		"data":    response,
	})
}

// GetSubjectGradeByID handles fetching a single subject-grade relationship by ID
// @Summary Get subject-grade relationship by ID
// @Description Get a single subject-grade relationship by its ID
// @Tags Subject Grades
// @Accept json
// @Produce json
// @Param id path string true "Subject-Grade Relationship ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subject-grades/{id} [get]
func (h *SubjectGradeController) GetSubjectGradeByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject-grade relationship ID is required",
		})
		return
	}

	subjectGrade, err := h.subjectGradeService.GetSubjectGradeByID(id)
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
		"message":       "Subject-grade relationship retrieved successfully",
		"subject_grade": subjectGrade,
	})
}

// GetSubjectsByGrade handles fetching all subjects for a grade
// @Summary Get subjects by grade
// @Description Get all subjects associated with a specific grade
// @Tags Subject Grades
// @Accept json
// @Produce json
// @Param grade_id path string true "Grade ID"
// @Param academic_year query string false "Filter by academic year"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subject-grades/grade/{grade_id} [get]
func (h *SubjectGradeController) GetSubjectsByGrade(c *gin.Context) {
	gradeID := c.Param("grade_id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade ID is required",
		})
		return
	}

	academicYear := c.Query("academic_year")

	subjects, err := h.subjectGradeService.GetSubjectsByGrade(gradeID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Subjects retrieved successfully",
		"subjects": subjects,
	})
}

// GetGradesBySubject handles fetching all grades for a subject
// @Summary Get grades by subject
// @Description Get all grades associated with a specific subject
// @Tags Subject Grades
// @Accept json
// @Produce json
// @Param subject_id path string true "Subject ID"
// @Param academic_year query string false "Filter by academic year"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subject-grades/subject/{subject_id} [get]
func (h *SubjectGradeController) GetGradesBySubject(c *gin.Context) {
	subjectID := c.Param("subject_id")
	if subjectID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	academicYear := c.Query("academic_year")

	grades, err := h.subjectGradeService.GetGradesBySubject(subjectID, academicYear)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Grades retrieved successfully",
		"grades":  grades,
	})
}

// DeleteSubjectGrade handles deleting a subject-grade relationship
// @Summary Delete a subject-grade relationship
// @Description Soft delete a subject-grade relationship by ID
// @Tags Subject Grades
// @Accept json
// @Produce json
// @Param id path string true "Subject-Grade Relationship ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subject-grades/{id} [delete]
func (h *SubjectGradeController) DeleteSubjectGrade(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject-grade relationship ID is required",
		})
		return
	}

	err := h.subjectGradeService.DeleteSubjectGrade(id)
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
		"message": "Subject-grade relationship deleted successfully",
	})
}
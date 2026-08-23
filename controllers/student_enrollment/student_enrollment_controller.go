// handlers/student_enrollment_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/student_enrollment"
)

type StudentEnrollmentHandler struct {
	enrollmentService *services.StudentEnrollmentService
}

func NewStudentEnrollmentHandler(enrollmentService *services.StudentEnrollmentService) *StudentEnrollmentHandler {
	return &StudentEnrollmentHandler{
		enrollmentService: enrollmentService,
	}
}

// CreateStudentEnrollment handles the creation of a new student enrollment
// @Summary Create a student enrollment
// @Description Enroll a student in a grade
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param request body dto.CreateStudentEnrollmentRequest true "Enrollment request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments [post]
func (h *StudentEnrollmentHandler) CreateStudentEnrollment(c *gin.Context) {
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

	var req dto.CreateStudentEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	enrollment, err := h.enrollmentService.CreateStudentEnrollment(&req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already enrolled") {
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
		"message":    "Student enrolled successfully",
		"enrollment": enrollment,
	})
}

// BulkCreateStudentEnrollments handles bulk creation of student enrollments
// @Summary Bulk create student enrollments
// @Description Enroll multiple students in a grade
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param request body dto.BulkCreateStudentEnrollmentRequest true "Bulk enrollment request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments/bulk [post]
func (h *StudentEnrollmentHandler) BulkCreateStudentEnrollments(c *gin.Context) {
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

	var req dto.BulkCreateStudentEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.enrollmentService.BulkCreateStudentEnrollments(&req, userID)
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
		"message": "Students enrolled successfully",
		"data":    result,
	})
}

// GetAllStudentEnrollments handles fetching all student enrollments
// @Summary Get all student enrollments
// @Description Get a paginated list of student enrollments with optional filtering
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param student_id query string false "Filter by student ID"
// @Param grade_id query string false "Filter by grade ID"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.StudentEnrollmentListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/fetch/student-enrollments [get]
func (h *StudentEnrollmentHandler) GetAllStudentEnrollments(c *gin.Context) {
	var params dto.StudentEnrollmentQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.enrollmentService.GetAllStudentEnrollments(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Enrollments retrieved successfully",
		"data":    response,
	})
}

// GetStudentEnrollmentByID handles fetching a single enrollment by ID
// @Summary Get enrollment by ID
// @Description Get a single student enrollment by its ID
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments/{id} [get]
func (h *StudentEnrollmentHandler) GetStudentEnrollmentByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Enrollment ID is required",
		})
		return
	}

	enrollment, err := h.enrollmentService.GetStudentEnrollmentByID(id)
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
		"message":    "Enrollment retrieved successfully",
		"enrollment": enrollment,
	})
}

// GetEnrollmentsByStudent handles fetching all enrollments for a student
// @Summary Get enrollments by student
// @Description Get all enrollments for a specific student
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param student_id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments/student/{student_id} [get]
func (h *StudentEnrollmentHandler) GetEnrollmentsByStudent(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Student ID is required",
		})
		return
	}

	enrollments, err := h.enrollmentService.GetEnrollmentsByStudent(studentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Enrollments retrieved successfully",
		"enrollments": enrollments,
	})
}

// GetEnrollmentsByGrade handles fetching all enrollments for a grade
// @Summary Get enrollments by grade
// @Description Get all enrollments for a specific grade
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param grade_id path string true "Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments/grade/{grade_id} [get]
func (h *StudentEnrollmentHandler) GetEnrollmentsByGrade(c *gin.Context) {
	gradeID := c.Param("grade_id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade ID is required",
		})
		return
	}

	enrollments, err := h.enrollmentService.GetEnrollmentsByGrade(gradeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Enrollments retrieved successfully",
		"enrollments": enrollments,
	})
}

// GetCurrentEnrollmentByStudent handles fetching the current enrollment for a student
// @Summary Get current enrollment by student
// @Description Get the current active enrollment for a student
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param student_id path string true "Student ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments/student/{student_id}/current [get]
func (h *StudentEnrollmentHandler) GetCurrentEnrollmentByStudent(c *gin.Context) {
	studentID := c.Param("student_id")
	if studentID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Student ID is required",
		})
		return
	}

	enrollment, err := h.enrollmentService.GetCurrentEnrollmentByStudent(studentID)
	if err != nil {
		if strings.Contains(err.Error(), "no active enrollment") {
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
		"message":    "Current enrollment retrieved successfully",
		"enrollment": enrollment,
	})
}

// UpdateStudentEnrollment handles updating an existing enrollment
// @Summary Update an enrollment
// @Description Update an existing student enrollment by ID
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Param request body dto.UpdateStudentEnrollmentRequest true "Enrollment update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments/{id} [put]
func (h *StudentEnrollmentHandler) UpdateStudentEnrollment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Enrollment ID is required",
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

	var req dto.UpdateStudentEnrollmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	enrollment, err := h.enrollmentService.UpdateStudentEnrollment(id, &req, userID)
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
		"message":    "Enrollment updated successfully",
		"enrollment": enrollment,
	})
}

// DeleteStudentEnrollment handles deleting an enrollment (soft delete)
// @Summary Delete an enrollment
// @Description Soft delete a student enrollment by ID
// @Tags Student Enrollments
// @Accept json
// @Produce json
// @Param id path string true "Enrollment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-enrollments/{id} [delete]
func (h *StudentEnrollmentHandler) DeleteStudentEnrollment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Enrollment ID is required",
		})
		return
	}

	err := h.enrollmentService.DeleteStudentEnrollment(id)
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
		"message": "Enrollment deleted successfully",
	})
}
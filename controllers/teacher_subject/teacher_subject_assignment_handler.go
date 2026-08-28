// handlers/teacher_subject_assignment_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"crm-go/dto"
	"crm-go/services/teacher_subject"
)

type TeacherSubjectAssignmentHandler struct {
	assignmentService *services.TeacherSubjectAssignmentService
}

func NewTeacherSubjectAssignmentHandler(assignmentService *services.TeacherSubjectAssignmentService) *TeacherSubjectAssignmentHandler {
	return &TeacherSubjectAssignmentHandler{
		assignmentService: assignmentService,
	}
}

// CreateAssignment handles the creation of a new subject assignment
// @Summary Create a subject assignment
// @Description Assign a subject to a teacher for a specific grade
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param request body dto.CreateTeacherSubjectAssignmentRequest true "Assignment request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments [post]
func (h *TeacherSubjectAssignmentHandler) CreateAssignment(c *gin.Context) {
	var req dto.CreateTeacherSubjectAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	assignment, err := h.assignmentService.CreateAssignment(&req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already assigned") {
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
		"message":    "Subject assigned successfully",
		"assignment": assignment,
	})
}

// BulkAssignSubjects handles bulk assignment of subjects
// @Summary Bulk assign subjects
// @Description Assign multiple subjects to a teacher for a specific grade
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param request body dto.BulkAssignSubjectsRequest true "Bulk assignment request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments/bulk [post]
func (h *TeacherSubjectAssignmentHandler) BulkAssignSubjects(c *gin.Context) {
	var req dto.BulkAssignSubjectsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.assignmentService.BulkAssignSubjects(&req)
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
		"message": "Subjects assigned successfully",
		"data":    result,
	})
}

// GetAllAssignments handles fetching all assignments
// @Summary Get all assignments
// @Description Get a paginated list of all subject-teacher assignments
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param grade_id query string false "Filter by grade ID"
// @Param subject_id query string false "Filter by subject ID"
// @Param teacher_id query string false "Filter by teacher ID"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.TeacherSubjectAssignmentListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments [get]
func (h *TeacherSubjectAssignmentHandler) GetAllAssignments(c *gin.Context) {
	var params dto.TeacherSubjectAssignmentQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	response, err := h.assignmentService.GetAllAssignments(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignments retrieved successfully",
		"data":    response,
	})
}

// GetAssignmentByID handles fetching a single assignment by ID
// @Summary Get assignment by ID
// @Description Get a single subject-teacher assignment by its ID
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments/{id} [get]
func (h *TeacherSubjectAssignmentHandler) GetAssignmentByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Assignment ID is required",
		})
		return
	}

	assignment, err := h.assignmentService.GetAssignmentByID(id)
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
		"message":    "Assignment retrieved successfully",
		"assignment": assignment,
	})
}

// GetAssignmentsByTeacher handles fetching all assignments for a teacher
// @Summary Get assignments by teacher
// @Description Get all subject assignments for a specific teacher
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param teacher_id path string true "Teacher ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments/teacher/{teacher_id} [get]
func (h *TeacherSubjectAssignmentHandler) GetAssignmentsByTeacher(c *gin.Context) {
	teacherID := c.Param("teacher_id")
	if teacherID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Teacher ID is required",
		})
		return
	}

	assignments, err := h.assignmentService.GetAssignmentsByTeacher(teacherID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Assignments retrieved successfully",
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// GetAssignmentsByGrade handles fetching all assignments for a grade
// @Summary Get assignments by grade
// @Description Get all subject-teacher assignments for a specific grade
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param grade_id path string true "Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments/grade/{grade_id} [get]
func (h *TeacherSubjectAssignmentHandler) GetAssignmentsByGrade(c *gin.Context) {
	gradeID := c.Param("grade_id")
	if gradeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Grade ID is required",
		})
		return
	}

	assignments, err := h.assignmentService.GetAssignmentsByGrade(gradeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Assignments retrieved successfully",
		"assignments": assignments,
		"total":       len(assignments),
	})
}

// UpdateAssignment handles updating an assignment
// @Summary Update an assignment
// @Description Update an existing subject-teacher assignment
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Param request body dto.UpdateTeacherSubjectAssignmentRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments/{id} [put]
func (h *TeacherSubjectAssignmentHandler) UpdateAssignment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Assignment ID is required",
		})
		return
	}

	var req dto.UpdateTeacherSubjectAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	assignment, err := h.assignmentService.UpdateAssignment(id, &req)
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
		"message":    "Assignment updated successfully",
		"assignment": assignment,
	})
}

// DeleteAssignment handles deleting an assignment (soft delete)
// @Summary Delete an assignment
// @Description Soft delete a subject-teacher assignment
// @Tags Teacher Subject Assignments
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/teacher-subject-assignments/{id} [delete]
func (h *TeacherSubjectAssignmentHandler) DeleteAssignment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Assignment ID is required",
		})
		return
	}

	err := h.assignmentService.DeleteAssignment(id)
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
		"message": "Assignment deleted successfully",
	})
}
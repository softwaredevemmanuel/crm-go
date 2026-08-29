// handlers/arm_class_teacher_handler.go
package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/arm_class_teacher"
)

type ArmClassTeacherHandler struct {
	assignmentService *services.ArmClassTeacherService
}

func NewArmClassTeacherHandler(assignmentService *services.ArmClassTeacherService) *ArmClassTeacherHandler {
	return &ArmClassTeacherHandler{
		assignmentService: assignmentService,
	}
}

// CreateAssignment handles the creation of a new arm class teacher assignment
// @Summary Create an arm class teacher assignment
// @Description Assign a class teacher to an arm
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param request body dto.CreateArmClassTeacherRequest true "Assignment request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers [post]
func (h *ArmClassTeacherHandler) CreateAssignment(c *gin.Context) {
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

	var req dto.CreateArmClassTeacherRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	assignment, err := h.assignmentService.CreateAssignment(&req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "already") {
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
		"message":    "Class teacher assigned successfully",
		"assignment": assignment,
	})
}

// BulkAssignClassTeachers handles bulk assignment of class teachers
// @Summary Bulk assign class teachers
// @Description Assign multiple class teachers to arms
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param request body dto.BulkAssignClassTeachersRequest true "Bulk assignment request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers/bulk [post]
func (h *ArmClassTeacherHandler) BulkAssignClassTeachers(c *gin.Context) {
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

	var req dto.BulkAssignClassTeachersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	result, err := h.assignmentService.BulkAssignClassTeachers(&req, userID)
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
		"message": "Class teachers assigned successfully",
		"data":    result,
	})
}

// GetAllAssignments handles fetching all assignments
// @Summary Get all arm class teacher assignments
// @Description Get a paginated list of all arm class teacher assignments
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param arm_id query string false "Filter by arm ID"
// @Param teacher_id query string false "Filter by teacher ID"
// @Param status query string false "Filter by status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} dto.ArmClassTeacherListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers [get]
func (h *ArmClassTeacherHandler) GetAllAssignments(c *gin.Context) {
	var params dto.ArmClassTeacherQueryParams
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
// @Description Get a single arm class teacher assignment by its ID
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers/{id} [get]
func (h *ArmClassTeacherHandler) GetAssignmentByID(c *gin.Context) {
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
// @Description Get all arm class teacher assignments for a specific teacher
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param teacher_id path string true "Teacher ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers/teacher/{teacher_id} [get]
func (h *ArmClassTeacherHandler) GetAssignmentsByTeacher(c *gin.Context) {
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

// GetAssignmentsByArm handles fetching all assignments for an arm
// @Summary Get assignments by arm
// @Description Get all class teacher assignments for a specific arm
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param arm_id path string true "Arm ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers/arm/{arm_id} [get]
func (h *ArmClassTeacherHandler) GetAssignmentsByArm(c *gin.Context) {
	armID := c.Param("arm_id")
	if armID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Arm ID is required",
		})
		return
	}

	assignments, err := h.assignmentService.GetAssignmentsByArm(armID)
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



// GetArmsWithClassTeachers handles fetching all arms with their class teachers
// @Summary Get arms with class teachers
// @Description Get all arms with their assigned class teachers
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers/arms-with-teachers [get]
func (h *ArmClassTeacherHandler) GetArmsWithClassTeachers(c *gin.Context) {
	arms, err := h.assignmentService.GetArmsWithClassTeachers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Arms retrieved successfully",
		"data":    arms,
		"total":   len(arms),
	})
}

// UpdateAssignment handles updating an assignment
// @Summary Update an assignment
// @Description Update an existing arm class teacher assignment
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Param request body dto.UpdateArmClassTeacherRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers/{id} [put]
func (h *ArmClassTeacherHandler) UpdateAssignment(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Assignment ID is required",
		})
		return
	}

	var req dto.UpdateArmClassTeacherRequest
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
// @Description Soft delete an arm class teacher assignment
// @Tags Arm Class Teacher
// @Accept json
// @Produce json
// @Param id path string true "Assignment ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/arm-class-teachers/{id} [delete]
func (h *ArmClassTeacherHandler) DeleteAssignment(c *gin.Context) {
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
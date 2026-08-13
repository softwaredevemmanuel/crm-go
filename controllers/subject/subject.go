// handlers/subject_handler.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/subjects"
)

type SubjectHandler struct {
	service *services.SubjectService
}

func NewSubjectController(service *services.SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

// ============ 1. CREATE SUBJECT ============
// @Summary Create a new subject
// @Tags Subjects
// @Accept json
// @Produce json
// @Param request body dto.CreateSubjectRequest true "Subject creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects [post]
func (h *SubjectHandler) CreateSubject(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.service.CreateSubject(&req, uuid.MustParse(userID.(string)))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "subject with this name already exists" ||
			err.Error() == "subject with this code already exists" {
			status = http.StatusConflict
		}
		if err.Error() == "department not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Subject created successfully",
		"subject": subject,
	})
}

// ============ 2. GET ALL SUBJECTS ============
// @Summary Get all subjects with pagination
// @Tags Subjects
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search by name, code, or description"
// @Param status query string false "Filter by status"
// @Param department_id query string false "Filter by department ID"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order" default(desc)
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects [get]
func (h *SubjectHandler) GetAllSubjects(c *gin.Context) {
	var params dto.SubjectQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetAllSubjects(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subjects retrieved successfully",
		"data":    result,
	})
}


// ============ 5. GET SUBJECT BY ID WITH DEPARTMENT AND HEAD ============
// @Summary Get subject by ID with department and head of department details
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/department/head-of-department/{id} [get]
func (h *SubjectHandler) GetSubjectWithDepartmentAndHead(c *gin.Context) {
	id := c.Param("id")

	subject, err := h.service.GetSubjectWithDepartmentAndHead(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "subject not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subject retrieved successfully",
		"subject": subject,
	})
}

// ============ 6. UPDATE SUBJECT ============
// @Summary Update a subject
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Param request body dto.UpdateSubjectRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/{id} [put]
func (h *SubjectHandler) UpdateSubject(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject, err := h.service.UpdateSubject(id, &req, uuid.MustParse(userID.(string)))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "subject not found" {
			status = http.StatusNotFound
		}
		if err.Error() == "subject with this name already exists" ||
			err.Error() == "subject with this code already exists" {
			status = http.StatusConflict
		}
		if err.Error() == "department not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subject updated successfully",
		"subject": subject,
	})
}

// ============ 7. DELETE SUBJECT ============
// @Summary Delete a subject (soft delete)
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/{id} [delete]
func (h *SubjectHandler) DeleteSubject(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.DeleteSubject(id, uuid.MustParse(userID.(string))); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "subject not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Subject deleted successfully",
	})
}
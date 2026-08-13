// handlers/department_handler.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/department"
)

type DepartmentHandler struct {
	service *services.DepartmentService
}

func NewDepartmentHandler(service *services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

// ============ 1. CREATE DEPARTMENT ============
// @Summary Create a new department
// @Tags Departments
// @Accept json
// @Produce json
// @Param request body dto.CreateDepartmentRequest true "Department creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/departments [post]
func (h *DepartmentHandler) CreateDepartment(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	department, err := h.service.CreateDepartment(&req, uuid.MustParse(userID.(string)))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "department with this name already exists" ||
			err.Error() == "department with this code already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "Department created successfully",
		"department": department,
	})
}

// ============ 2. GET ALL DEPARTMENTS ============
// @Summary Get all departments with pagination
// @Tags Departments
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param search query string false "Search by name, code, or description"
// @Param status query string false "Filter by status"
// @Param sort_by query string false "Sort by field" default(created_at)
// @Param sort_order query string false "Sort order" default(desc)
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/departments [get]
func (h *DepartmentHandler) GetAllDepartments(c *gin.Context) {
	var params dto.DepartmentQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.service.GetAllDepartments(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Departments retrieved successfully",
		"data":    result,
	})
}

// ============ 3. GET DEPARTMENT WITH SUBJECTS ============
// @Summary Get department with its subjects
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/departments/{id}/subjects [get]
func (h *DepartmentHandler) GetDepartmentWithSubjects(c *gin.Context) {
	id := c.Param("id")

	department, err := h.service.GetDepartmentWithSubjects(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "department not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Department retrieved successfully",
		"department": department,
	})
}

// ============ 4. GET DEPARTMENT WITH HEAD AND SUBJECTS ============
// @Summary Get department with head of department and subjects
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/departments/{id}/head-subjects [get]
func (h *DepartmentHandler) GetDepartmentWithHeadAndSubjects(c *gin.Context) {
	id := c.Param("id")

	department, err := h.service.GetDepartmentWithHeadAndSubjects(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "department not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Department retrieved successfully",
		"department": department,
	})
}

// ============ 5. GET DEPARTMENT BY ID ============
// @Summary Get department by ID with all details
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/departments/{id} [get]
func (h *DepartmentHandler) GetDepartmentByID(c *gin.Context) {
	id := c.Param("id")

	department, err := h.service.GetDepartmentByID(id)
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "department not found" {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Department retrieved successfully",
		"department": department,
	})
}

// ============ 6. UPDATE DEPARTMENT ============
// @Summary Update a department
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Param request body dto.UpdateDepartmentRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/departments/{id} [put]
func (h *DepartmentHandler) UpdateDepartment(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req dto.UpdateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	department, err := h.service.UpdateDepartment(id, &req, uuid.MustParse(userID.(string)))
	if err != nil {
		status := http.StatusBadRequest
		if err.Error() == "department not found" {
			status = http.StatusNotFound
		}
		if err.Error() == "department with this name already exists" ||
			err.Error() == "department with this code already exists" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "Department updated successfully",
		"department": department,
	})
}

// ============ 7. DELETE DEPARTMENT ============
// @Summary Delete a department (soft delete)
// @Tags Departments
// @Accept json
// @Produce json
// @Param id path string true "Department ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/departments/{id} [delete]
func (h *DepartmentHandler) DeleteDepartment(c *gin.Context) {
	id := c.Param("id")

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := h.service.DeleteDepartment(id, uuid.MustParse(userID.(string))); err != nil {
		status := http.StatusBadRequest
		if err.Error() == "department not found" {
			status = http.StatusNotFound
		}
		if err.Error() == "cannot delete department: it has associated subjects" {
			status = http.StatusConflict
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Department deleted successfully",
	})
}
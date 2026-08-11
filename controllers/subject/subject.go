package controllers
import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"crm-go/dto"
	"crm-go/services"
	"gorm.io/gorm"

)

type SubjectController struct {
	db          *gorm.DB
	subjectService *services.SubjectService
}

func NewSubjectController(db *gorm.DB, subjectService *services.SubjectService) *SubjectController {
	return &SubjectController{
		db: 			db,
		subjectService: subjectService,
	}
}

// CreateSubject handles the creation of a new subject
// @Summary Create a new subject
// @Description Create a new subject with the provided details
// @Tags Subjects
// @Accept json
// @Produce json
// @Param request body dto.CreateSubjectRequest true "Subject creation request"
// @Success 201 {object} dto.SubjectResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects [post]
func (ctl *SubjectController) CreateSubject(c *gin.Context) {
	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user ID not found",
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Parse request body
	var req dto.CreateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Trim whitespace from fields
	req.Name = strings.TrimSpace(req.Name)
	req.Code = strings.TrimSpace(req.Code)
	req.Description = strings.TrimSpace(req.Description)
	req.Department = strings.TrimSpace(req.Department)

	// Create subject using service
	subject, err := ctl.subjectService.CreateSubject(&req, userID)
	if err != nil {
		// Handle specific errors
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

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "Subject created successfully",
		"subject": subject,
	})
}


// GetAllSubjects handles fetching all subjects with pagination and filters
// @Summary Get all subjects
// @Description Get a paginated list of subjects with optional filtering
// @Tags Subjects
// @Accept json
// @Produce json
// @Param search query string false "Search by name, code, or description"
// @Param department query string false "Filter by department"
// @Param status query string false "Filter by status (active, inactive)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort by field (name, code, department, credits, created_at, updated_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} dto.SubjectListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subjects [get]
func (h *SubjectController) GetAllSubjects(c *gin.Context) {
	// Parse query parameters
	var params dto.SubjectQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Get subjects from service
	response, err := h.subjectService.GetAllSubjects(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response
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
// @Router /subjects/{id} [get]
func (h *SubjectController) GetSubjectByID(c *gin.Context) {
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

// GetSubjectDepartments handles fetching all unique departments
// @Summary Get all subject departments
// @Description Get a list of all unique departments
// @Tags Subjects
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /subjects/departments [get]
func (h *SubjectController) GetSubjectDepartments(c *gin.Context) {
	departments, err := h.subjectService.GetSubjectDepartments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     "Departments retrieved successfully",
		"departments": departments,
	})
}


// UpdateSubject handles updating an existing subject
// @Summary Update a subject
// @Description Update an existing subject by ID
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Param request body dto.UpdateSubjectRequest true "Subject update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/{id} [put]
func (h *SubjectController) UpdateSubject(c *gin.Context) {
	// Get subject ID from URL
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user ID not found",
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Parse request body
	var req dto.UpdateSubjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Check if at least one field is being updated
	if req.Name == "" && req.Code == "" && req.Description == "" &&
		req.Department == "" && req.Credits == 0 && req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one field must be provided for update",
		})
		return
	}

	// Update subject using service
	subject, err := h.subjectService.UpdateSubject(id, &req, userID)
	if err != nil {
		// Handle specific errors
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

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Subject updated successfully",
		"subject": subject,
	})
}

// DeleteSubject handles deleting a subject (soft delete)
// @Summary Delete a subject
// @Description Soft delete a subject by ID
// @Tags Subjects
// @Accept json
// @Produce json
// @Param id path string true "Subject ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/subjects/{id} [delete]
func (h *SubjectController) DeleteSubject(c *gin.Context) {
	// Get subject ID from URL
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Subject ID is required",
		})
		return
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized: user ID not found",
		})
		return
	}

	// Parse user ID
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	// Delete subject using service
	err = h.subjectService.DeleteSubject(id, userID)
	if err != nil {
		// Handle specific errors
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "being used") {
			c.JSON(http.StatusConflict, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Subject deleted successfully",
	})
}

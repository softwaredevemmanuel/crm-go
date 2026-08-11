package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/class_grade"
	"gorm.io/gorm"

)

type ClassGradeController struct {
	db  *gorm.DB
	classGradeService *services.ClassGradeService
}

func NewClassGradeController(db *gorm.DB, classGradeService *services.ClassGradeService) *ClassGradeController {
	return &ClassGradeController{
		db: 			db,

		classGradeService: classGradeService,
	}
}

// CreateClassGrade handles the creation of a new class grade
// @Summary Create a new class grade
// @Description Create a new class grade with the provided details (Grade 1-6)
// @Tags Class Grades
// @Accept json
// @Produce json
// @Param request body dto.CreateClassGradeRequest true "Class grade creation request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/class-grades [post]
func (h *ClassGradeController) CreateClassGrade(c *gin.Context) {
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
	var req dto.CreateClassGradeRequest
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
	req.AcademicYear = strings.TrimSpace(req.AcademicYear)

	// Create class grade using service
	classGrade, err := h.classGradeService.CreateClassGrade(&req, userID)
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
		"message":    "Class grade created successfully",
		"class_grade": classGrade,
	})
}


// GetAllClassGrades handles fetching all class grades with pagination and filters
// @Summary Get all class grades
// @Description Get a paginated list of class grades with optional filtering
// @Tags Class Grades
// @Accept json
// @Produce json
// @Param search query string false "Search by name, code, or description"
// @Param level query int false "Filter by level (1-6)"
// @Param academic_year query string false "Filter by academic year"
// @Param status query string false "Filter by status (active, inactive, archived)"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort by field (name, code, level, academic_year, capacity, status, created_at)" default(created_at)
// @Param sort_order query string false "Sort order (asc, desc)" default(desc)
// @Success 200 {object} dto.ClassGradeListResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/class-grades [get]
func (h *ClassGradeController) GetAllClassGrades(c *gin.Context) {
	// Parse query parameters
	var params dto.ClassGradeQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	// Get class grades from service
	response, err := h.classGradeService.GetAllClassGrades(&params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"message": "Class grades retrieved successfully",
		"data":    response,
	})
}

// GetClassGradeByID handles fetching a single class grade by ID
// @Summary Get class grade by ID
// @Description Get a single class grade by its ID
// @Tags Class Grades
// @Accept json
// @Produce json
// @Param id path string true "Class Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/class-grades/{id} [get]
func (h *ClassGradeController) GetClassGradeByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Class grade ID is required",
		})
		return
	}

	classGrade, err := h.classGradeService.GetClassGradeByID(id)
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
		"message":     "Class grade retrieved successfully",
		"class_grade": classGrade,
	})
}

// GetAcademicYears handles fetching all unique academic years
// @Summary Get all academic years
// @Description Get a list of all unique academic years
// @Tags Class Grades
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/class-grades/academic-years [get]
func (h *ClassGradeController) GetAcademicYears(c *gin.Context) {
	academicYears, err := h.classGradeService.GetAcademicYears()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Academic years retrieved successfully",
		"academic_years": academicYears,
	})
}

// GetLevels handles fetching all unique levels
// @Summary Get all levels
// @Description Get a list of all unique levels (1-6)
// @Tags Class Grades
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/class-grades/levels [get]
func (h *ClassGradeController) GetLevels(c *gin.Context) {
	levels, err := h.classGradeService.GetLevels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Levels retrieved successfully",
		"levels":  levels,
	})
}

// UpdateClassGrade handles updating an existing class grade
// @Summary Update a class grade
// @Description Update an existing class grade by ID
// @Tags Class Grades
// @Accept json
// @Produce json
// @Param id path string true "Class Grade ID"
// @Param request body dto.UpdateClassGradeRequest true "Class grade update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/class-grades/{id} [put]
func (h *ClassGradeController) UpdateClassGrade(c *gin.Context) {
	// Get class grade ID from URL
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Class grade ID is required",
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
	var req dto.UpdateClassGradeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Check if at least one field is being updated
	if req.Name == "" && req.Code == "" && req.Level == 0 && 
		req.Description == "" && req.AcademicYear == "" && req.Capacity == 0 && req.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "At least one field must be provided for update",
		})
		return
	}

	// Update class grade using service
	classGrade, err := h.classGradeService.UpdateClassGrade(id, &req, userID)
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
		"message":     "Class grade updated successfully",
		"class_grade": classGrade,
	})
}

// DeleteClassGrade handles deleting a class grade (soft delete)
// @Summary Delete a class grade
// @Description Soft delete a class grade by ID
// @Tags Class Grades
// @Accept json
// @Produce json
// @Param id path string true "Class Grade ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/class-grades/{id} [delete]
func (h *ClassGradeController) DeleteClassGrade(c *gin.Context) {
	// Get class grade ID from URL
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Class grade ID is required",
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

	// Delete class grade using service
	err = h.classGradeService.DeleteClassGrade(id, userID)
	if err != nil {
		// Handle specific errors
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "has students enrolled") {
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
		"message": "Class grade deleted successfully",
	})
}
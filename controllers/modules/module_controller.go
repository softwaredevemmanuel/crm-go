package modules

import (
	"crm-go/config"
	"crm-go/models"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateModules handles the creation of a new module
// @Summary Create a new modules
// @Description Create a new module
// @Tags Modules
// @Accept json
// @Produce json
// @Param module body models.ModuleInput true "Module"
// @Success 201 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 409 {object} models.ConflictResponse
// @Failure 500 {object} models.FailureResponse
// @Router /api/modules [post]
// @Security BearerAuth
func CreateModule(c *gin.Context) {
	var input models.ModuleInput
	db := config.DB

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var existingModule models.Module

	// ✅ Check for duplicate by module number within the same course

	err := db.
		Where("course_id = ? AND module_number = ?", input.CourseID, input.ModuleNumber).
		First(&existingModule).Error

	if err == nil {
		c.JSON(http.StatusConflict, models.ConflictResponse{Error: "Module number already exists for this course"})
		return
	}
	// ✅ Check for duplicate by title
	var existingCourse models.Module
	if err := db.Where("title = ?", input.Title).First(&existingCourse).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Module with the same title already exists"})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Real DB error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to validate module number",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Real DB error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to validate module number",
		})
		return
	}

	module := models.Module{
		ID:            uuid.New(),
		CourseID:      input.CourseID,
		Title:         input.Title,
		Slug:          input.Slug,
		Description:   input.Description,
		ModuleNumber:  input.ModuleNumber,
		IsFree:        input.IsFree,
		Status:        input.Status,
		EstimatedTime: input.EstimatedTime,
	}

	if err := config.DB.Create(&module).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Module created successfully",
	})
}

// GetAllModules handles the retrieval of all modules
// @Summary Get all modules
// @Description Get all modules
// @Tags Modules
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Failure 500 {object} models.FailureResponse
// @Router /modules [get]
func GetAllModules(c *gin.Context) {
	var modules []models.Module

	if err := config.DB.
		Preload("Course", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "title")
		}).
		Order("module_number ASC").
		Find(&modules).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch modules"})
		return
	}

	var response []models.ModuleResponse

	for _, ch := range modules {
		var r models.ModuleResponse
		r.ID = ch.ID
		r.CourseID = ch.CourseID
		r.Title = ch.Title
		r.Slug = ch.Slug
		r.Description = ch.Description
		r.ModuleNumber = ch.ModuleNumber
		r.IsFree = ch.IsFree
		r.Status = ch.Status
		r.EstimatedTime = ch.EstimatedTime
		r.TotalTopics = ch.TotalTopics
		r.TotalDuration = ch.TotalDuration
		r.CreatedAt = ch.CreatedAt
		r.UpdatedAt = ch.UpdatedAt

		response = append(response, r)
	}

	c.JSON(http.StatusOK, response)
}

// GetModuleByID handles the retrieval of a single module by ID
// @Summary Get a single module by ID
// @Description Get a single module by ID
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.FailureResponse
// @Router /modules/{id} [get]
func GetModuleByID(c *gin.Context) {
	id := c.Param("id")

	moduleID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var module models.Module

	if err := config.DB.
		Preload("Course").
		Preload("Topics").
		First(&module, "id = ?", moduleID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	// Map Course
	course := models.CourseMiniResponse{
		ID:    module.Course.ID,
		Title: module.Course.Title,
	}

	// Optional: pick one topic (or latest)
	var topic *models.TopicMiniResponse
	if module.Topics != nil && len(*module.Topics) > 0 {
		topic = &models.TopicMiniResponse{
			ID:          (*module.Topics)[0].ID,
			Title:       (*module.Topics)[0].Title,
			ContentType: (*module.Topics)[0].ContentType,
			ContentURL:  (*module.Topics)[0].ContentURL,
		}
	}

	response := models.ModuleViewResponse{
		ID:            module.ID,
		CourseID:      module.CourseID,
		Title:         module.Title,
		Slug:          module.Slug,
		Description:   module.Description,
		ModuleNumber:  module.ModuleNumber,
		IsFree:        module.IsFree,
		Status:        module.Status,
		EstimatedTime: module.EstimatedTime,
		TotalDuration: module.TotalDuration,
		CreatedAt:     module.CreatedAt,
		UpdatedAt:     module.UpdatedAt,
		Course:        course,
		Topics:        topic,
	}

	c.JSON(http.StatusOK, response)
}

// UpdateModule handles the updating of an existing module
// @Summary Update an existing module
// @Description Update an existing module
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Param module body models.ModuleInput true "Module"
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.FailureResponse
// @Router /api/modules/{id} [put]
func UpdateModule(c *gin.Context) {
	id := c.Param("id")

	moduleID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	var module models.Module
	if err := config.DB.First(&module, "id = ?", moduleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	var input models.ModuleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Prevent duplicate module numbers per course
	var count int64
	config.DB.Model(&models.Module{}).
		Where("course_id = ? AND module_number = ? AND id != ?",
			input.CourseID,
			input.ModuleNumber,
			moduleID,
		).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Module number already exists for this course",
		})
		return
	}

	updates := map[string]interface{}{
		"course_id":      input.CourseID,
		"title":          input.Title,
		"slug":           input.Slug,
		"description":    input.Description,
		"module_number":  input.ModuleNumber,
		"is_free":        input.IsFree,
		"status":         input.Status,
		"estimated_time": input.EstimatedTime,
	}

	if err := config.DB.Model(&module).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Reload with relations
	if err := config.DB.
		Preload("Course").
		Preload("Topics").
		First(&module, "id = ?", moduleID).Error; err != nil {

		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module updated successfully", "details": module})
}

// DeleteModule handles the deletion of an existing module
// @Summary Delete an existing module
// @Description Delete an existing module
// @Tags Modules
// @Accept json
// @Produce json
// @Param id path string true "Module ID"
// @Success 200 {object} models.DeleteSuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.NotFoundResponse
// @Failure 500 {object} models.FailureResponse
// @Router /api/modules/{id} [delete]
func DeleteModule(c *gin.Context) {
	id := c.Param("id")

	moduleID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid module ID"})
		return
	}

	result := config.DB.Delete(&models.Module{}, "id = ?", moduleID)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete module"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Module deleted successfully"})
}

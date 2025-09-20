package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



// CreateCategory godoc
// @Summary Create a new category
// @Description Admin can create a new course category
// @Tags Categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 201 {object} CategoryResponse
// @Failure 400 {object} models.ErrorResponse "Invalid request payload"
// @Failure 409 {object} models.ErrorResponse "Category with this name already exists"
// @Failure 500 {object} models.ErrorResponse "Failed to create category"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Param   category body models.CategoryInput true "Category data"
// @Router /api/categories [post]
func CreateCategory(c *gin.Context) {
	var category models.Category

	// Bind JSON request
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	// ✅ Check if category with the same name already exists
	var existing models.Category
	if err := db.Where("name = ?", category.Name).First(&existing).Error; err == nil {
		// Found a duplicate
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error:   "Duplicate Error",
			Message: "Category with this name already exists",
		})
		return
	}

	// ✅ Create new category
	if err := db.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, models.Category{
			ID:   category.ID,
			Name: category.Name,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
		return
	}

	c.JSON(http.StatusCreated, category)
}


// GetCategories - List all categories
func GetCategories(c *gin.Context) {
	var categories []models.Category
	db := config.DB
	db.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

// GetCourse - Get one course by ID
func GetCourseByID(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	db := config.GetDB()
	var course models.Course
	if err := db.First(&course, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	c.JSON(http.StatusOK, course)
}

// UpdateCategory - Update a category
func UpdateCategory(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	db := config.GetDB()
	var category models.Category
	if err := db.First(&category, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&category)
	c.JSON(http.StatusOK, category)
}

// DeleteCategory - Delete a category
func DeleteCategory(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	db := config.GetDB()
	var category models.Category

	// Check if category exists before deleting
	if err := db.First(&category, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Delete category
	if err := db.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted successfully"})
}


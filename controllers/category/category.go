package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)



// Create A Category
// @Summary Create a new category
// @Description Admin can create a new course category
// @Tags Categories
// @Accept  json
// @Produce  json
// @Security BearerAuth
// @Success 201 {object} models.Category
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
			Description: category.Description,
			CreatedAt: category.CreatedAt,
			UpdatedAt: category.UpdatedAt,
		})
		return
	}

	c.JSON(http.StatusCreated, category)
}


// Get All Categories 
// @Summary Get all categories
// @Description Retrieve a list of all available categories
// @Tags Categories
// @Produce  json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string "Internal Server Error"
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category
	db := config.DB
	db.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

// UpdateCategory godoc
// @Summary Update a category
// @Description Update an existing category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Category ID"
// @Param category body models.CategoryInput true "Updated Category Data"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string "Invalid category ID or bad request body"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 409 {object} models.ErrorResponse "Category with this name already exists"
// @Router /api/categories/{id} [put]
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

	// Check if category exists
	if err := db.First(&category, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	// Bind new data
	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ✅ Check for duplicate name (excluding current category)
	var existing models.Category
	if err := db.Where("name = ? AND id <> ?", input.Name, uid).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, models.ErrorResponse{
			Error:   "Duplicate Error",
			Message: "Category with this name already exists",
		})
		return
	}

	// Update fields
	category.Name = input.Name
	category.Description = input.Description

	// Save
	if err := db.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(http.StatusOK, category)
}



// DeleteCategory godoc
// @Summary Delete a category
// @Description Delete an existing category by ID
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Security BearerAuth
// @Success 200 {object} map[string]string "Category deleted successfully"
// @Failure 400 {object} map[string]string "Invalid category ID"
// @Failure 404 {object} map[string]string "Category not found"
// @Failure 500 {object} map[string]string "Failed to delete category"
// @Router /api/categories/{id} [delete]
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



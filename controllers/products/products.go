package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateProduct - Admin creates a product
// @Summary      Create a new product
// @Description  Admin creates a new product. Prevents duplicates by name.
// @Tags         Products
// @Accept       json
// @Produce      json
// @Param request body models.ProductInput true "Product Payload"
// @Success      201     {object}  models.Product
// @Failure      400     {object}  map[string]string "Invalid request body"
// @Failure      409     {object}  map[string]string "Product with this name already exists"
// @Failure      500     {object}  map[string]string "Failed to create product"
// @Router       /api/products [post]
// @Security BearerAuth
func CreateProduct(c *gin.Context) {
	var product models.Product

	// Bind JSON request
	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	// ✅ Check if product with the same name already exists
	var existing models.Product
	if err := db.Where("name = ?", product.Name).First(&existing).Error; err == nil {
		// Found a duplicate
		c.JSON(http.StatusConflict, gin.H{"error": "Product with this name already exists"})
		return
	}

	// ✅ Create new product
	if err := db.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, product)
}


// GetProducts - List all products
// @Summary      List all products
// @Description  Retrieve a list of all products
// @Tags         Products
// @Produce      json
// @Success      200  {array}   models.Product
// @Failure      500  {object}  map[string]string "Failed to fetch products"
// @Router       /products [get]	
func GetProducts(c *gin.Context) {
	var products []models.Product
	db := config.DB
	db.Find(&products)
	c.JSON(http.StatusOK, products)
}

// GetProductByID - Get one product by ID
// @Summary      Get product by ID
// @Description  Retrieve a single product by its ID
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {object}  models.Product
// @Failure      400  {object}  map[string]string "Invalid product ID"
// @Failure      404  {object}  map[string]string "Product not found"
// @Router       /products/{id} [get]	
func GetProductByID(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := config.GetDB()
	var product models.Product
	if err := db.First(&product, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

type ListCoursesByProduct struct {
	ID          string `json:"course_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type ListCategoryCourses struct {
	ID          string `json:"course_id" binding:"required,uuid4"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

// GetCoursesByProduct - fetch all courses linked to a specific product
// @Summary      Get related courses with product ID
// @Description  Retrieve all courses associated with a specific product via the pivot table
// @Tags         Products
// @Produce      json
// @Param        id   path      string  true  "Product ID"
// @Success      200  {array}   ListCoursesByProduct
// @Failure      400  {object}  map[string]string "Invalid product ID"
// @Failure      500  {object}  map[string]string "Failed to fetch courses"
// @Router       /products/{id}/products [get]
// GetProductDetailsWithCourses - Get product details with its related courses
func GetProductDetailsWithCourses(c *gin.Context) {
	productID := c.Param("id")
	db := config.DB

	// Fetch product details
	var product models.Product
	if err := db.First(&product, "id = ?", productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Fetch courses related to this product
	var courses []models.Course
	if err := db.Joins("JOIN course_product_tables ON courses.id = course_product_tables.course_id").
		Where("course_product_tables.product_id = ?", productID).
		Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	// Map courses to DTO
	var relatedCourses []ListCategoryCourses
	for _, course := range courses {
		relatedCourses = append(relatedCourses, ListCategoryCourses{
			ID:          course.ID.String(),
			Title:       course.Title,
			Description: course.Description,
			Image:       course.Image,
		})
	}

	// Final response
	response := gin.H{
		"id":          product.ID.String(),
		"name":        product.Name,
		"description": product.Description,
		"price":       product.Price,
		"courses":     relatedCourses,
	}

	c.JSON(http.StatusOK, response)
}


// UpdateProduct - Update a product
// @Summary      Update a product
// @Description  Update product details by its ID
// @Tags         Products	
// @Accept       json
// @Produce      json
// @Param        id      path      string         true  "Product ID"
// @Param        request body      models.ProductInput true  "Product Payload"
// @Success      200     {object}  models.Product
// @Failure      400     {object}  map[string]string "Invalid product ID or request body"
// @Failure      404     {object}  map[string]string "Product not found"
// @Failure      500     {object}  map[string]string "Failed to update product"
// @Router       /api/products/{id} [put]
// @Security BearerAuth		
func UpdateProduct(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := config.GetDB()
	var product models.Product
	if err := db.First(&product, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	if err := c.ShouldBindJSON(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&product)
	c.JSON(http.StatusOK, product)
}

	

// DeleteProduct - Delete a product
// @Summary      Delete a product
// @Description  Delete a product by its ID
// @Tags         Products	
// @Param        id   path      string  true  "Product ID"	
// @Success      200  {object}  map[string]string "Product deleted successfully"
// @Failure      400  {object}  map[string]string "Invalid product ID"
// @Failure      404  {object}  map[string]string "Product not found"
// @Failure      500  {object}  map[string]string "Failed to delete product"
// @Router       /api/products/{id} [delete]
// @Security BearerAuth	
func DeleteProduct(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := config.GetDB()
	var product models.Product

	// Check if product exists before deleting
	if err := db.First(&product, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Delete product
	if err := db.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

package products

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateProduct - Admin creates a product
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
func GetProducts(c *gin.Context) {
	var products []models.Product
	db := config.DB
	db.Find(&products)
	c.JSON(http.StatusOK, products)
}

// GetProductByID - Get one product by ID
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

// UpdateProduct - Update a product
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



func GetProductWithCourseMates(c *gin.Context) {
	productID := c.Param("id")
	db := config.DB

	// Get the product itself
	var product models.Product
	if err := db.First(&product, "id = ?", productID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	// Find course this product belongs to
	var courseProducts []models.CourseProduct
	db.Where("product_id = ?", productID).Find(&courseProducts)

	if len(courseProducts) == 0 {
		c.JSON(http.StatusOK, gin.H{"product": product, "related_products": []models.Product{}})
		return
	}

	// Get related products in the same course
	var relatedProducts []models.Product
	db.Joins("JOIN course_products ON products.id = course_products.product_id").
		Where("course_products.course_id IN (?) AND products.id != ?", 
			db.Select("course_id").Where("product_id = ?", productID).Table("course_products"), 
			productID).
		Find(&relatedProducts)

	c.JSON(http.StatusOK, gin.H{
		"product":          product,
		"related_products": relatedProducts,
	})
}
package controllers

import (
	"crm-go/config"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"errors"
	"log"
)




// CreateCourseProduct - Admin creates a course product relationship
func CreateCourseProduct(c *gin.Context) {
    var request struct {
        CourseID   string `json:"course_id" binding:"required,uuid4"`
        ProductID string `json:"product_id" binding:"required,uuid4"`
    }

    // Bind JSON request
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Parse UUIDs
    courseUUID, err := uuid.Parse(request.CourseID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID format"})
        return
    }

    productUUID, err := uuid.Parse(request.ProductID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID format"})
        return
    }

    db := config.DB

    // Debug log
    log.Printf("✅ value of courseId: %v", courseUUID)
    log.Printf("✅ value of productId: %v", productUUID)

    // ✅ Check if course exists
    var course models.Course
    if err := db.Where("id = ?", courseUUID).First(&course).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Course ID does not exist"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check course existence"})
        }
        return
    }

    // ✅ Check if product exists
    var product models.Product
    if err := db.Where("id = ?", productUUID).First(&product).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Product ID does not exist"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check product existence"})
        }
        return
    }

    // ✅ Check if course-product relationship already exists
    var existing models.CourseProduct
    if err := db.Where("course_id = ? AND product_id = ?", courseUUID, productUUID).First(&existing).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Course already exists in this product"})
        return
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing relationship"})
        return
    }

    // ✅ Create new course product relationship
    courseProduct := models.CourseProduct{
        CourseID:   courseUUID,
        ProductID: productUUID,
    }

    if err := db.Create(&courseProduct).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course product relationship"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":          courseProduct.ID,
        "course_id":   courseProduct.CourseID,
        "product_id": courseProduct.ProductID,
        "message":     "Course product created successfully",
    })
}


func GetCoursesByProduct(c *gin.Context) {
	productID := c.Param("id")

	var courses []models.Course
	db := config.DB

	err := db.Joins("JOIN course_products ON courses.id = course_products.course_id").
		Where("course_products.product_id = ?", productID).
		Find(&courses).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
		return
	}

	c.JSON(http.StatusOK, courses)
}


// DeleteCategory - Delete a course category
func DeleteCourseProduct(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := config.GetDB()
	var course_product models.CourseProduct

	// Check if course product exists before deleting
	if err := db.First(&course_product, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course Product not found"})
		return
	}

	// Delete course product
	if err := db.Delete(&course_product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course product"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course Product deleted successfully"})
}

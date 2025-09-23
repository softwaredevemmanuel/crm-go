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
// @Summary      Create a course-product relationship
// @Description  Admin creates a relationship between a course and a product. Prevents duplicates.
// @Tags         Course Products
// @Accept       json
// @Produce      json
// @Param request body models.CreateCourseProductRequest true "Course-Category Payload"
// @Success      201     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]string "Invalid request body or IDs"
// @Failure      404     {object}  map[string]string "Course ID does not exist" or "Product ID does not exist"
// @Failure      409     {object}  map[string]string "Product already exists for this course"
// @Failure      500     {object}  map[string]string "Failed to create course product relationship"
// @Router       /api/course-products [post]
// @Security BearerAuth
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
    var existing models.CourseProductTable
    if err := db.Where("course_id = ? AND product_id = ?", courseUUID, productUUID).First(&existing).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Product already exists for this course"})
        return
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing relationship"})
        return
    }

    // ✅ Create new course product relationship
    courseProduct := models.CourseProductTable{
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




// DeleteCategory - Delete a course category
// @Summary      Delete a course-product relationship
// @Description  Remove a relationship between a course and a product by its ID
// @Tags         Course Products
// @Param        id   path      string  true  "Course Product ID"
// @Success      200  {object}  map[string]string "Course Product deleted successfully"
// @Failure      400  {object}  map[string]string "Invalid product ID"
// @Failure      404  {object}  map[string]string "Course Product not found"
// @Failure      500  {object}  map[string]string "Failed to delete course product"
// @Router       /api/course-products/{id} [delete]
// @Security BearerAuth
func DeleteCourseProduct(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	db := config.GetDB()
	var course_product models.CourseProductTable

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

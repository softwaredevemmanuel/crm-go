package controllers

import (
	"net/http"
	"crm-go/config"
	"crm-go/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
)

// CreateCourse - Admin creates a course
func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB
	if err := db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}

// GetCourses - List all courses
func GetCourses(c *gin.Context) {
	var courses []models.Course
	db := config.DB
	db.Find(&courses)
	c.JSON(http.StatusOK, courses)
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


func GetCourseWithCategoryMates(c *gin.Context) {
    courseID := c.Param("id")
    db := config.DB

    // Get the course itself
    var course models.Course
    if err := db.First(&course, "id = ?", courseID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
        return
    }

    // Find categories this course belongs to
    var courseCategories []models.CourseCategory
    db.Where("course_id = ?", courseID).Find(&courseCategories)

    if len(courseCategories) == 0 {
        c.JSON(http.StatusOK, gin.H{"course": course, "related_courses": []models.Course{}})
        return
    }

    // Get related courses in the same categories
    var relatedCourses []models.Course
    db.Joins("JOIN course_categories ON courses.id = course_categories.course_id").
        Where("course_categories.category_id IN (?) AND courses.id != ?", 
            db.Select("category_id").Where("course_id = ?", courseID).Table("course_categories"), 
            courseID).
        Find(&relatedCourses)

    c.JSON(http.StatusOK, gin.H{
        "course":          course,
        "related_courses": relatedCourses,
    })
}



// GetCourseWithProducts - fetch a course and its related products
func GetCourseWithProducts(c *gin.Context) {
	courseID := c.Param("id")
	db := config.DB

	// Parse UUID from param
	uid, err := uuid.Parse(courseID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	// Get the course itself
	var course models.Course
	if err := db.First(&course, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Get all products linked to this course (from the pivot table course_products)
	var relatedProducts []models.Product
	err = db.Joins("JOIN course_products cp ON cp.product_id = products.id").
		Where("cp.course_id = ?", uid).
		Find(&relatedProducts).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch related products"})
		return
	}

	// Return course with related products
	c.JSON(http.StatusOK, gin.H{
		"course":           course,
		"related_products": relatedProducts,
	})
}



// UpdateCourse - Update a course
func UpdateCourse(c *gin.Context) {
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

	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Save(&course)
	c.JSON(http.StatusOK, course)
}

// DeleteCourse - Delete a course
func DeleteCourse(c *gin.Context) {
	id := c.Param("id")

	// Parse UUID
	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}

	db := config.GetDB()
	var course models.Course

	// Check if course exists before deleting
	if err := db.First(&course, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}

	// Delete course
	if err := db.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}


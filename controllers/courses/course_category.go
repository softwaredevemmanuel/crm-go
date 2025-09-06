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




// CreateCourseCategory - Admin creates a course category relationship
func CreateCourseCategory(c *gin.Context) {
    var request struct {
        CourseID   string `json:"course_id" binding:"required,uuid4"`
        CategoryID string `json:"category_id" binding:"required,uuid4"`
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

    categoryUUID, err := uuid.Parse(request.CategoryID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
        return
    }

    db := config.DB

    // Debug log
    log.Printf("✅ value of courseId: %v", courseUUID)
    log.Printf("✅ value of categoryId: %v", categoryUUID)

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

    // ✅ Check if category exists
    var category models.Category
    if err := db.Where("id = ?", categoryUUID).First(&category).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID does not exist"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check category existence"})
        }
        return
    }

    // ✅ Check if course-category relationship already exists
    var existing models.CourseCategoryTable
    if err := db.Where("course_id = ? AND category_id = ?", courseUUID, categoryUUID).First(&existing).Error; err == nil {
        c.JSON(http.StatusConflict, gin.H{"error": "Course already exists in this category"})
        return
    } else if !errors.Is(err, gorm.ErrRecordNotFound) {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing relationship"})
        return
    }

    // ✅ Create new course category relationship
    courseCategory := models.CourseCategoryTable{
        CourseID:   courseUUID,
        CategoryID: categoryUUID,
    }

    if err := db.Create(&courseCategory).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course category relationship"})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "id":          courseCategory.ID,
        "course_id":   courseCategory.CourseID,
        "category_id": courseCategory.CategoryID,
        "message":     "Course category created successfully",
    })
}



func GetCoursesByCategory(c *gin.Context) {
    categoryID := c.Param("id")

    var courses []models.Course
    db := config.DB

    err := db.Joins("JOIN course_category_tables ON courses.id = course_category_tables.course_id").
        Where("course_category_tables.category_id = ?", categoryID).
        Find(&courses).Error

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch courses"})
        return
    }

    c.JSON(http.StatusOK, courses)
}


// DeleteCategory - Delete a course category
func DeleteCourseCategory(c *gin.Context) {
    id := c.Param("id")

    // Parse UUID
    uid, err := uuid.Parse(id)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
        return
    }

    db := config.GetDB()
    var course_category models.CourseCategoryTable

    // Check if course category exists before deleting
    if err := db.First(&course_category, "id = ?", uid).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Course Category not found"})
        return
    }

    // Delete course category
    if err := db.Delete(&course_category).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course category"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Course Category deleted successfully"})
}


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




// CreateCourseCategory godoc
// @Summary Create Course-Category relationships
// @Description Admin creates relationships between multiple courses and a category. Accepts an array of course IDs.
// @Tags Course Categories
// @Accept json
// @Produce json
// @Param request body models.CreateCourseCategoryRequest true "Course-Category Payload"
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/category-courses [post]
// @Security BearerAuth
func CreateCourseCategory(c *gin.Context) {
    var request struct {
        CourseIDs  []string `json:"course_id" binding:"required"`
        CategoryID string   `json:"category_id" binding:"required,uuid4"`
    }
    log.Printf("✅ value of courseId: 1")

    // Bind JSON request
    if err := c.ShouldBindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    log.Printf("✅ value of courseId: 2")

    // Validate that course_ids array is not empty
    if len(request.CourseIDs) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "At least one course ID is required"})
        return
    }
    log.Printf("✅ value of courseId: 3")

    // Parse category UUID
    categoryUUID, err := uuid.Parse(request.CategoryID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
        return
    }

    db := config.DB

    // ✅ Start a transaction
    tx := db.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start transaction"})
        return
    }
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()

    // ✅ Check if category exists
    var category models.Category
    if err := tx.Where("id = ?", categoryUUID).First(&category).Error; err != nil {
        tx.Rollback()
        if errors.Is(err, gorm.ErrRecordNotFound) {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Category ID does not exist"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check category existence"})
        }
        return
    }

    // Parse all course UUIDs and validate
    var courseUUIDs []uuid.UUID
    var invalidCourseIDs []string
    
    for _, courseID := range request.CourseIDs {
        parsedUUID, err := uuid.Parse(courseID)
        if err != nil {
            invalidCourseIDs = append(invalidCourseIDs, courseID)
            continue
        }
        courseUUIDs = append(courseUUIDs, parsedUUID)
    }

    if len(invalidCourseIDs) > 0 {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "error":              "Invalid course ID format",
            "invalid_course_ids": invalidCourseIDs,
        })
        return
    }

    // ✅ Check if all courses exist and get their details
    var courseDetails []models.Course
    var nonExistentCourses []string
    
    for _, courseUUID := range courseUUIDs {
        var course models.Course
        if err := tx.Where("id = ?", courseUUID).First(&course).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                nonExistentCourses = append(nonExistentCourses, courseUUID.String())
            }
        } else {
            courseDetails = append(courseDetails, course)
        }
    }

    if len(nonExistentCourses) > 0 {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "error":                 "Some courses do not exist",
            "non_existent_courses": nonExistentCourses,
        })
        return
    }

    // ✅ Check existing relationships and filter
    var courseCategories []models.CourseCategoryTable
    var alreadyExists []string
    var createdCourses []models.Course

    for _, course := range courseDetails {
        var existing models.CourseCategoryTable
        if err := tx.Where("course_id = ? AND category_id = ?", course.ID, categoryUUID).First(&existing).Error; err == nil {
            alreadyExists = append(alreadyExists, course.ID.String())
        } else if errors.Is(err, gorm.ErrRecordNotFound) {
            // Create new relationship
            courseCategories = append(courseCategories, models.CourseCategoryTable{
                CourseID:   course.ID,
                CategoryID: categoryUUID,
            })
            createdCourses = append(createdCourses, course)
        } else {
            tx.Rollback()
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to check existing relationships",
            })
            return
        }
    }

    // If all courses already exist, return conflict
    if len(courseCategories) == 0 {
        tx.Rollback()
        c.JSON(http.StatusConflict, gin.H{
            "error":            "All courses already exist in this category",
            "existing_courses": alreadyExists,
        })
        return
    }

    // ✅ Batch insert new relationships
    if err := tx.Create(&courseCategories).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to create course category relationships",
        })
        return
    }

    // ✅ Commit transaction
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to commit transaction",
        })
        return
    }

    // Return success response
    var createdIDs []string
    var createdCourseNames []string
    
    for i, cc := range courseCategories {
        createdIDs = append(createdIDs, cc.ID.String())
        if i < len(createdCourses) {
            createdCourseNames = append(createdCourseNames, createdCourses[i].Title)
        }
    }

    // Prepare response message
    message := "Course(s) added to category successfully"
    if len(alreadyExists) > 0 {
        message = "Some courses were added successfully. Some already existed."
    }

    c.JSON(http.StatusCreated, gin.H{
        "message":             message,
        "created_relationships": createdIDs,
        "added_courses":       len(courseCategories),
        "added_course_names":  createdCourseNames,
        "existing_courses":    alreadyExists,
        "total_courses":       len(request.CourseIDs),
    })
}




// DeleteCourseCategory godoc
// @Summary      Delete a course category relationship
// @Description  Remove a course-category relationship by its ID
// @Tags         Course Categories
// @Param        id   path      string  true  "Course Category ID (UUID)"
// @Produce      json
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/category-courses/{id} [delete]
// @Security BearerAuth
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


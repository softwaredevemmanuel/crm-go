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


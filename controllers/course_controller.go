package controllers

import (
	"net/http"
	"crm-go/config"
	"crm-go/models"
	"github.com/gin-gonic/gin"
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
func GetCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course

	db := config.DB
	if err := db.First(&course, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Course not found"})
		return
	}
	c.JSON(http.StatusOK, course)
}

// UpdateCourse - Update a course
func UpdateCourse(c *gin.Context) {
	id := c.Param("id")
	var course models.Course
	db := config.DB

	if err := db.First(&course, id).Error; err != nil {
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
	db := config.DB

	if err := db.Delete(&models.Course{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete course"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

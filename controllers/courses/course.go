package controllers

import (
	"net/http"
	"crm-go/config"
	"crm-go/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	
)
type CreateCourseInput struct {
    Title            string   `json:"title" binding:"required"`
    Description      string   `json:"description" binding:"required"`
    Image            string   `json:"image"`
    VideoURL         string   `json:"video_url"`
    TutorID          string   `json:"tutor_id" binding:"required,uuid4"`
    LearningOutcomes []string `json:"learning_outcomes"`
    Requirements     []string `json:"requirements"`
}

type CourseResponse struct {
    ID              string   `json:"id"`
    Title           string   `json:"title"`
    Description     string   `json:"description"`
    Image           string   `json:"image"`
    VideoURL        string   `json:"video_url"`
    TutorID         string   `json:"tutor_id"`
    LearningOutcomes []string `json:"learning_outcomes"`
    Requirements    []string `json:"requirements"`
}

// CreateCourse godoc
// @Summary      Create a new course
// @Description  Admin creates a new course. Prevents duplicate titles.
// @Tags         Courses
// @Accept       json
// @Produce      json
// @Param        course  body      CreateCourseInput  true  "Course details"
// @Success      201     {object}  CourseResponse
// @Failure      400     {object}  map[string]string "Invalid request body"
// @Failure      409     {object}  map[string]string "Course with the same title already exists"
// @Failure      500     {object}  map[string]string "Failed to create course"
// @Router       /api/courses [post]
// @Security BearerAuth
func CreateCourse(c *gin.Context) {
	var course models.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	// ✅ Check for duplicate by title
	var existingCourse models.Course
	if err := db.Where("title = ?", course.Title).First(&existingCourse).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Course with the same title already exists"})
		return
	}

	// Create course
	if err := db.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create course"})
		return
	}

	c.JSON(http.StatusCreated, course)
}


// GetCourses godoc
// @Summary      List all courses
// @Description  Retrieve all courses available in the system
// @Tags         Courses
// @Produce      json
// @Success      200  {array}   CourseResponse
// @Failure      500  {object}  map[string]string "Failed to fetch courses"
// @Router       /courses [get]
func GetCourses(c *gin.Context) {
	var courses []models.Course
	db := config.DB
	db.Find(&courses)
	c.JSON(http.StatusOK, courses)
}

// GetCourseByID godoc
// @Summary      Get a course by ID
// @Description  Retrieve details of a specific course using its ID
// @Tags         Courses
// @Produce      json
// @Param        id   path      string  true  "Course ID"
// @Success      200  {object}  CourseResponse
// @Failure      400  {object}  map[string]string "Invalid course ID"
// @Failure      404  {object}  map[string]string "Course not found"
// @Router       /courses/{id} [get]
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

type CourseWithCategoryMatesResponse struct {
    Course         CourseResponse   `json:"course"`
    RelatedCourses []CourseResponse `json:"related_courses"`
}





// GetCourseWithProducts - fetch a course and its related products
//@Summary      Get Products related to a course ID
//@Description  Retrieve a course along with products linked to it via the pivot table
//@Tags         Courses
//@Produce      json
//@Param        id   path      string  true  "Course ID"
//@Success      200  {object}  map[string]interface{}
//@Failure      400  {object}  map[string]string "Invalid course ID"
//@Failure      404  {object}  map[string]string "Course not found"
//@Failure      500  {object}  map[string]string "Failed to fetch related products"
//@Router       /courses/{id}/products [get]
func GetProductsWithRalatedCourseID(c *gin.Context) {
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

	// Get all products linked to this course (from the pivot table course_product_tables)
	var relatedProducts []models.Product
	err = db.Joins("JOIN course_product_tables cp ON cp.product_id = products.id").
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
//@Summary      Update an existing course
//@Description  Admin updates course details by ID
//@Tags         Courses
//@Accept       json
//@Produce      json
//@Param        id      path      string            true  "Course ID"
//@Param        course  body      CreateCourseInput true  "Updated course details"
//@Success      200     {object}  CourseResponse
//@Failure      400     {object}  map[string]string "Invalid course ID or request body"
//@Failure      404     {object}  map[string]string "Course not found"
//@Failure      500     {object}  map[string]string "Failed to update course"
//@Router       /api/courses/{id} [put]
//@Security BearerAuth	
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
//@Summary      Delete a course
//@Description  Admin deletes a course by ID
//@Tags         Courses
//@Produce      json
//@Param        id   path      string  true  "Course ID"
//@Success      200  {object}  map[string]string "Course deleted successfully"
//@Failure      400  {object}  map[string]string "Invalid course ID"
//@Failure      404  {object}  map[string]string "Course not found"
//@Failure      500  {object}  map[string]string "Failed to delete course"
//@Router       /api/courses/{id} [delete]
//@Security BearerAuth
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


package enrollments

import (
	"net/http"
	"crm-go/config"
	"crm-go/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"errors"
	"gorm.io/gorm"
	
)
type EnrollmentResponse struct {
    ID              string   `json:"id"`
    StudentID       string   `json:"student_id"`
    CourseID        string   `json:"course_id"`
}
// GetEnrollments godoc
// @Summary      List enrollments
// @Description  Get all enrollments
// @Tags         Enrollments
// @Accept       json
// @Produce      json
// @Success      200     {object}  EnrollmentResponse
// @Failure      500 {object} map[string]string
// @Router       /enrollments [get]
func GetEnrollments(c *gin.Context) {
	var enrollments []models.Enrollment
	if err := config.DB.Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollments)
}


// GetEnrollmentByID retrieves an enrollment by its ID
// @Summary      Get enrollment by ID
// @Description  Get an enrollment by its ID
// @Tags         Enrollments
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Enrollment ID"
// @Success      200  {object}  EnrollmentResponse
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /enrollments/{id} [get]
func GetEnrollmentByID(c *gin.Context) {
	id := c.Param("id")
	var enrollment models.Enrollment
	if err := config.DB.First(&enrollment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}
	c.JSON(http.StatusOK, enrollment)
}

type CreateEnrollmentInput struct {
	StudentID string `json:"student_id" binding:"required,uuid4"`
	CourseID  string `json:"course_id" binding:"required,uuid4"`
}
// CreateEnrollment creates a new enrollment
// @Summary      Create a new enrollment
// @Description  Create a new enrollment with the provided details
// @Tags         Enrollments
// @Accept       json
// @Produce      json
// @Param        enrollment  body      CreateEnrollmentInput  true  "Course details"
// @Success      201  {object}  EnrollmentResponse
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /enrollments [post]
func CreateEnrollment(c *gin.Context) {
	var enrollment models.Enrollment

	if err := c.ShouldBindJSON(&enrollment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db := config.DB

	// 🔒 Check for existing enrollment (student_id + course_id)
	var existing models.Enrollment
	err := db.
		Where("student_id = ? AND course_id = ?", enrollment.StudentID, enrollment.CourseID).
		First(&existing).Error

	if err == nil {
		// Found → duplicate
		c.JSON(http.StatusConflict, gin.H{
			"error": "Student is already enrolled in this course",
		})
		return
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		// Unexpected DB error
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create enrollment
	enrollment.ID = uuid.New()
	if err := db.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, enrollment)
}


// UpdateEnrollment updates an existing enrollment
func UpdateEnrollment(c *gin.Context) {
	id := c.Param("id")
	var enrollment models.Enrollment
	if err := config.DB.First(&enrollment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}
	if err := c.ShouldBindJSON(&enrollment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Save(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, enrollment)
}


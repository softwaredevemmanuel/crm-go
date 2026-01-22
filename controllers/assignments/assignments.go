package assignments

import (
	"github.com/gin-gonic/gin"
	"crm-go/models"
	"crm-go/config"
	"github.com/google/uuid"
	"net/http"
	"gorm.io/gorm"
	"errors"
	"time"
)

// AssignmentInput is used for creating/updating assignments via API
type AssignmentInput struct {
	CourseID           uuid.UUID  `json:"course_id" binding:"required"`
	ChapterID          *uuid.UUID `json:"chapter_id,omitempty"`
	LessonID           *uuid.UUID `json:"lesson_id,omitempty"`
	Title              string     `json:"title" binding:"required"`
	Slug               string     `json:"slug,omitempty"` // optional, can be auto-generated
	Description        string     `json:"description,omitempty"`
	LearningObjectives []string   `json:"learning_objectives,omitempty"`
	Type               string     `json:"type,omitempty" binding:"oneof=homework project essay quiz exam lab presentation discussion peer_review group research creative"`
	DueDate            time.Time  `json:"due_date" binding:"required"`
	PublishedBy        uuid.UUID  `json:"published_by" binding:"required"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"Invalid assignment ID"`
}
type FailureResponse struct {
	Error string `json:"error" example:"Failed to create assignment"`
}
type NotFoundResponse struct {
	Error string `json:"error" example:"Assignment not found"`
}
type SuccessResponse struct {
	Message string `json:"message" example:"Assignment created successfully"`
}
type DeleteSuccessResponse struct {
	Message string `json:"message" example:"Assignment deleted successfully"`
}

// CreateAssignment handles the creation of a new assignment
//@Summary Create a new assignment
//@Description Create a new assignment
//@Tags Assignments
//@Accept json
//@Produce json
//@Param assignment body AssignmentInput true "Assignment"
//@Success 201 {object} SuccessResponse
//@Failure 400 {object} ErrorResponse
//@Failure 500 {object} FailureResponse
//@Router /assignments [post]
func CreateAssignment(c *gin.Context) {
	var assignment models.Assignment

	if err := c.ShouldBindJSON(&assignment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	assignment.ID = uuid.New()

	if err := config.DB.Create(&assignment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create assignment",
		})
		return
	}

	c.JSON(http.StatusCreated, assignment)
}


// GetAssignments handles the retrieval of all assignments
//@Summary Get all assignments
//@Description Get all assignments
//@Tags Assignments
//@Accept json
//@Produce json
//@Success 200 {object} SuccessResponse
//@Failure 500 {object} FailureResponse
//@Router /assignments [get]
func GetAssignments(c *gin.Context) {
	var assignments []models.Assignment

	if err := config.DB.
		Preload("Course").
		Preload("Chapter").
		Preload("Lesson").
		Preload("Publisher").
		Order("due_date ASC").
		Find(&assignments).Error; err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch assignments",
		})
		return
	}

	c.JSON(http.StatusOK, assignments)
}


// GetAssignmentByID handles the retrieval of a single assignment by ID
//@Summary Get a single assignment by ID
//@Description Get a single assignment by ID
//@Tags Assignments
//@Accept json
//@Produce json
//@Param id path string true "Assignment ID"
//@Failure 400 {object} ErrorResponse
//@Failure 404 {object} NotFoundResponse
//@Failure 500 {object} FailureResponse
//@Router /assignments/{id} [get]
func GetAssignmentByID(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid assignment ID",
		})
		return
	}

	var assignment models.Assignment

	if err := config.DB.
		Preload("Course").
		Preload("Chapter").
		Preload("Lesson").
		Preload("Publisher").
		Preload("Submissions").
		First(&assignment, "id = ?", uid).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Assignment not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch assignment",
		})
		return
	}

	c.JSON(http.StatusOK, assignment)
}


// UpdateAssignment handles the updating of an existing assignment
//@Summary Update an existing assignment
//@Description Update an existing assignment
//@Tags Assignments
//@Accept json
//@Produce json
//@Param id path string true "Assignment ID"
//@Param assignment body models.Assignment true "Assignment"
//@Failure 400 {object} ErrorResponse
//@Failure 404 {object} NotFoundResponse
//@Failure 500 {object} FailureResponse
//@Router /assignments/{id} [put]
func UpdateAssignment(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid assignment ID",
		})
		return
	}

	var assignment models.Assignment

	if err := config.DB.First(&assignment, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Assignment not found",
		})
		return
	}

	var input models.Assignment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := config.DB.Model(&assignment).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update assignment",
		})
		return
	}

	c.JSON(http.StatusOK, assignment)
}

// DeleteAssignment handles the deletion of an existing assignment
//@Summary Delete an existing assignment
//@Description Delete an existing assignment
//@Tags Assignments
//@Accept json
//@Produce json
//@Param id path string true "Assignment ID"
//@Success 200 {object} DeleteSuccessResponse
//@Failure 400 {object} ErrorResponse
//@Failure 404 {object} NotFoundResponse
//@Failure 500 {object} FailureResponse
//@Router /assignments/{id} [delete]
func DeleteAssignment(c *gin.Context) {
	id := c.Param("id")

	uid, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid assignment ID",
		})
		return
	}

	if err := config.DB.Delete(&models.Assignment{}, "id = ?", uid).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete assignment",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Assignment deleted successfully",
	})
}

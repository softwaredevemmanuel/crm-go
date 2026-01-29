package assignmentsubmission

import (
	"crm-go/models"
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

    "crm-go/services/activity"
    "gorm.io/gorm"
)

// Create DTO
type CreateAssignmentSubmissionRequest struct {
    AssignmentID   uuid.UUID       `json:"assignment_id" binding:"required"`
    StudentID      uuid.UUID       `json:"student_id" binding:"required"`
    SubmissionType string          `json:"submission_type" binding:"required,oneof=text file url code video audio image document presentation multiple"`
    TextContent    string          `json:"text_content,omitempty"`
    FileURL        string          `json:"file_url,omitempty"`
    ExternalURL    string          `json:"external_url,omitempty"`
    CodeRepoURL    string          `json:"code_repo_url,omitempty"`
    Metadata       []string        `json:"metadata,omitempty"` // Raw JSON
    Status         string          `json:"status,omitempty"`   // Defaults to "submitted"
}

// Response DTO
type AssignmentSubmissionResponse struct {
    ID             uuid.UUID       `json:"id"`
    AssignmentID   uuid.UUID       `json:"assignment_id"`
    StudentID      uuid.UUID       `json:"student_id"`
    SubmissionType string          `json:"submission_type"`
    TextContent    string          `json:"text_content,omitempty"`
    FileURL        string          `json:"file_url,omitempty"`
    ExternalURL    string          `json:"external_url,omitempty"`
    CodeRepoURL    string          `json:"code_repo_url,omitempty"`
    Metadata       []string        `json:"metadata,omitempty"`
    Status         string          `json:"status"`
    SubmittedAt    time.Time       `json:"submitted_at"`
    CreatedAt      time.Time       `json:"created_at"`
    UpdatedAt      time.Time       `json:"updated_at"`
}




type AssignmentController struct {
	db       *gorm.DB
	activity *activity.Service
}

func NewAssignmentController(db *gorm.DB, activitySvc *activity.Service) *AssignmentController {
	return &AssignmentController{
		db:       db,
		activity: activitySvc,
	}
}

// CreateAssignmentSubmission creates a new assignment submission
//@Summary Create a new assignment submission
//@Description Create a new assignment submission
//@Tags Assignment Submissions
//@Accept json
//@Produce json
//@Param request body CreateAssignmentSubmissionRequest true "Create Assignment Submission Request"
//@Success 201 {object} models.SuccessResponse
//@Failure 400 {object} models.ErrorResponse
//@Failure 404 {object} models.ErrorResponse
//@Failure 403 {object} models.ErrorResponse
//@Router /api/assignment_submissions [post]
// @Security BearerAuth
func (ctl *AssignmentController) CreateAssignmentSubmission(c *gin.Context) {
	var req models.CreateAssignmentSubmissionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := ctl.db.Begin()

    // Validate assignment exists
	var assignment models.Assignment
	if err := tx.First(&assignment, "id = ?", req.AssignmentID).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusNotFound, gin.H{"error": "Assignment not found"})
		return
	}

	// Validate student is enrolled in the course
	var enrollment models.Enrollment
	if err := tx.Where("student_id = ? AND course_id = ?", req.StudentID, assignment.CourseID).First(&enrollment).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Student is not enrolled in this course",
		})
		return
    }

    // Check if submission already exists (prevent duplicates)
	var existingSubmission models.AssignmentSubmission
	if err := tx.Where("assignment_id = ? AND student_id = ?", req.AssignmentID, req.StudentID).First(&existingSubmission).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":        "Submission already exists for this assignment",
			"submission_id": existingSubmission.ID,
		})
		return
	}


	var status string
	var submittedAt time.Time

	// Check if assignment is past due date
	if !assignment.DueDate.IsZero() && time.Now().After(assignment.DueDate) {
		status = "late"
		submittedAt = time.Now()
	}

	submission := models.AssignmentSubmission{
		ID:             uuid.New(),
		AssignmentID:   req.AssignmentID,
		StudentID:      req.StudentID,
		SubmissionType: req.SubmissionType,
		TextContent:    req.TextContent,
		FileURL:       req.FileURL,
		ExternalURL:   req.ExternalURL,
		CodeRepoURL:   req.CodeRepoURL,
		Status:        status,
		SubmittedAt:   submittedAt,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := tx.Create(&submission).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 🔥 ACTIVITY LOG — CLEAN & REUSABLE
	_ = ctl.activity.Assignments.Submitted(
		tx,
		req.StudentID,
		assignment,
		submission,
		map[string]interface{}{
			"submission_type": req.SubmissionType,
		},
	)

	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Assignment submitted successfully",
		"id":      submission.ID,
	})
}

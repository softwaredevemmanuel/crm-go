package assignmentsubmission

import (
	"crm-go/config"
	controllers "crm-go/controllers/activitylogs"
	"crm-go/models"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/datatypes"
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

var activityLogger = controllers.NewActivityLogger(config.DB)


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
func CreateAssignmentSubmission(c *gin.Context) {
    var req models.CreateAssignmentSubmissionRequest
    
    // Bind and validate request
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid request data",
            "details": err.Error(),
        })
        return
    }

    // ✅ TRANSACTION STARTS HERE
    tx := config.DB.Begin()

    // Validate assignment exists
    var assignment models.Assignment
    if err := tx.First(&assignment, "id = ?", req.AssignmentID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Assignment not found",
        })
        return
    }

    // Validate student exists
    var student models.User
    if err := tx.First(&student, "id = ?", req.StudentID).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusNotFound, gin.H{
            "error": "Student not found",
        })
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
            "error": "Submission already exists for this assignment",
            "submission_id": existingSubmission.ID,
        })
        return
    }



    // Validate submission type vs provided content
    if err := validateSubmissionContent(req); err != nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid submission content",
            "details": err.Error(),
        })
        return
    }

    // Calculate if submission is late
    var submittedAt time.Time
    var status string = req.Status
    if status == "" {
        status = "submitted"
    }
    
    // Check if assignment is past due date
    if !assignment.DueDate.IsZero() && time.Now().After(assignment.DueDate) {
        status = "late"
    }
    submittedAt = time.Now()

    // Prepare submission data
    submission := models.AssignmentSubmission{
        ID:             uuid.New(),
        AssignmentID:   req.AssignmentID,
        StudentID:      req.StudentID,
        SubmissionType: req.SubmissionType,
        TextContent:    req.TextContent,
        FileURL:        req.FileURL,
        ExternalURL:    req.ExternalURL,
        CodeRepoURL:    req.CodeRepoURL,
        Status:         status,
        SubmittedAt:    submittedAt,
        CreatedAt:      time.Now(),
        UpdatedAt:      time.Now(),
    }

    // Handle metadata
    if len(req.Metadata) > 0 {
        // Validate JSON
        if !json.Valid(req.Metadata) {
            tx.Rollback()
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid JSON in metadata field",
            })
            return
        }
        submission.Metadata = datatypes.JSON(req.Metadata)
    }

    // Create the submission
    if err := tx.Create(&submission).Error; err != nil {
        tx.Rollback()
        
        // Handle specific errors
        if strings.Contains(err.Error(), "check constraint") {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": "Invalid submission type or status",
                "details": err.Error(),
            })
        } else if strings.Contains(err.Error(), "duplicate key") {
            c.JSON(http.StatusConflict, gin.H{
                "error": "Submission already exists",
            })
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{
                "error": "Failed to create submission: " + err.Error(),
            })
        }
        return
    }

    // Update assignment statistics (if needed)
    // if err := updateAssignmentStats(tx, req.AssignmentID); err != nil {
    //     // Log error but don't fail the submission
    //     log.Printf("Warning: Failed to update assignment stats: %v", err)
    // }

    // Create activity log
    // Create activity log with proper metadata
    metadata := map[string]interface{}{
        "assignment_id":   assignment.ID,
        "assignment_title": assignment.Title,
        "submission_type": req.SubmissionType,
        "course_id":       assignment.CourseID,
        "file_size":       len(req.FileURL) > 0, // Simplified, you might want actual file size
        "has_metadata":    len(req.Metadata) > 0,
    }
    
    // Log the activity
    if err := activityLogger.LogAssignmentSubmission(
        tx,
        req.StudentID,
        assignment.ID,
        submission.ID,
        assignment.Title,
        metadata,
    ); err != nil {
        // Log the error but don't fail the submission
        log.Printf("Warning: Failed to create activity log: %v", err)
    }



    tx.Commit()

    // Prepare response
    response := models.AssignmentSubmissionResponse{
        ID:             submission.ID,
        AssignmentID:   submission.AssignmentID,
        StudentID:      submission.StudentID,
        SubmissionType: submission.SubmissionType,
        TextContent:    submission.TextContent,
        FileURL:        submission.FileURL,
        ExternalURL:    submission.ExternalURL,
        CodeRepoURL:    submission.CodeRepoURL,
        Metadata:       json.RawMessage(submission.Metadata),
        Status:         submission.Status,
        SubmittedAt:    submission.SubmittedAt,
        CreatedAt:      submission.CreatedAt,
        UpdatedAt:      submission.UpdatedAt,
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Assignment submitted successfully",
        "data":    response,
    })
}

// validateSubmissionContent validates submission content based on type
func validateSubmissionContent(req models.CreateAssignmentSubmissionRequest) error {
    switch req.SubmissionType {
    case "text":
        if req.TextContent == "" {
            return fmt.Errorf("text content is required for text submissions")
        }
        if len(req.TextContent) > 10000 {
            return fmt.Errorf("text content exceeds maximum length of 10,000 characters")
        }
    case "file":
        if req.FileURL == "" {
            return fmt.Errorf("file URL is required for file submissions")
        }
    case "url":
        if req.ExternalURL == "" {
            return fmt.Errorf("external URL is required for URL submissions")
        }
        if _, err := url.ParseRequestURI(req.ExternalURL); err != nil {
            return fmt.Errorf("invalid URL format")
        }
    case "code":
        if req.CodeRepoURL == "" && req.FileURL == "" {
            return fmt.Errorf("either code repository URL or file URL is required for code submissions")
        }
    case "video", "audio", "image", "document", "presentation":
        if req.FileURL == "" && req.ExternalURL == "" {
            return fmt.Errorf("either file URL or external URL is required for %s submissions", req.SubmissionType)
        }
    case "multiple":
        // For multiple submissions, at least one content field should be present
        if req.FileURL == "" && req.ExternalURL == "" && req.TextContent == "" && req.CodeRepoURL == "" {
            return fmt.Errorf("at least one content field is required for multiple submissions")
        }
    }
    
    return nil
}

// updateAssignmentStats updates assignment statistics
// func updateAssignmentStats(tx *gorm.DB, assignmentID uuid.UUID) error {
//     var assignment models.Assignment
//     if err := tx.First(&assignment, "id = ?", assignmentID).Error; err != nil {
//         return err
//     }

//     // Count submissions
//     var submissionCount int64
//     if err := tx.Model(&models.AssignmentSubmission{}).
//         Where("assignment_id = ?", assignmentID).
//         Count(&submissionCount).Error; err != nil {
//         return err
//     }

//     // Update assignment with new submission count
//     return tx.Model(&assignment).
//         Update("submission_count", submissionCount).Error
// }
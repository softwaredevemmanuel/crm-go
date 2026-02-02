// controllers/grade_controller.go
package controllers

import (
    "net/http"
    
    "crm-go/models"
    "github.com/gin-gonic/gin"
    "fmt"
    "github.com/google/uuid"
)

// UpdateGrade handler
// @Summary Update a grade
// @Description Update an existing grade
// @Tags grades
// @Accept json
// @Produce json
// @Param id path string true "Grade ID"
// @Param grade body models.GradeUpdateInput true "Grade update data"
// @Success 200 {object} models.GradeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Router /api/grades/{id} [put]
// @Security BearerAuth
func (ctl *GradeController) UpdateGrade(c *gin.Context) {
    gradeID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid grade ID",
        })
        return
    }
    
    var req models.GradeUpdateInput
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input: " + err.Error(),
        })
        return
    }
    
    // Validate at least one field is being updated
    if req.Score == nil && req.Remarks == nil && req.AssignmentID == nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "no update data provided",
        })
        return
    }
    
    
    // Start transaction
    tx := ctl.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Update grade
    updatedGrade, err := ctl.gradeService.UpdateGradeWithTx(tx, gradeID, req)
    if err != nil {
        tx.Rollback()
        
        // Handle specific error types
        if err.Error() == "grade not found" {
            c.JSON(http.StatusNotFound, gin.H{
                "error": "grade not found",
            })
        } else {
            c.JSON(http.StatusBadRequest, gin.H{
                "error": err.Error(),
            })
        }
        return
    }
    
    // Log activity
        gradeModel := models.Grade{
            ID:           updatedGrade.ID,
            StudentID:    updatedGrade.StudentID,
            CourseID:     updatedGrade.CourseID,
            AssignmentID: updatedGrade.AssignmentID,
            Score:        updatedGrade.Score,
            Grade:        updatedGrade.Grade,
            Remarks:      updatedGrade.Remarks,
        }
        
        // Determine what was updated for activity details
        details := "Updated grade"
        if req.Score != nil {
            details += fmt.Sprintf(" - Score: %.2f", *req.Score)
        }
        if req.Remarks != nil {
            details += " - Remarks updated"
        }
        if req.AssignmentID != nil {
            details += " - Assignment changed"
        }
        
        _ = ctl.activity.Grades.Updated(tx, req.TutorID, gradeModel)
    
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save changes: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Grade updated successfully",
        "data":    updatedGrade,
    })
}

// BulkUpdateGrades handler
// @Summary Bulk update grades
// @Description Update multiple grades at once
// @Tags grades
// @Accept json
// @Produce json
// @Param grades body []models.BulkGradeUpdate true "Array of grade updates"
// @Success 200 {array} models.GradeResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/grades/bulk/update [put]
// @Security BearerAuth
func (ctl *GradeController) BulkUpdateGrades(c *gin.Context) {
    var updates []models.BulkGradeUpdate
    
    if err := c.ShouldBindJSON(&updates); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input: " + err.Error(),
        })
        return
    }
    
    if len(updates) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "no updates provided",
        })
        return
    }
    
    // Limit bulk operations
    if len(updates) > 100 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "cannot update more than 100 grades at once",
        })
        return
    }
    
    results, err := ctl.gradeService.BulkUpdateGrades(updates)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Grades updated successfully",
        "data":    results,
        "count":   len(results),
    })
}

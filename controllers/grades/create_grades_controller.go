// controllers/grade_controller.go
package controllers

import (
    "net/http"
    
    "crm-go/models"
    "crm-go/services/grades"
    "crm-go/services/activity"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type GradeController struct {
    db          *gorm.DB
    gradeService *services.GradeService
    activity    *activity.Service
}

func NewGradeController(db *gorm.DB, gradeService *services.GradeService, activitySvc *activity.Service) *GradeController {
    return &GradeController{
        db:          db,
        gradeService: gradeService,
        activity:    activitySvc,
    }
}

// CreateGrade handler
// @Summary Create a new grade
// @Description Create a grade for a student
// @Tags grades
// @Accept json
// @Produce json
// @Param grade body models.GradeInput true "Grade data"
// @Success 201 {object} models.GradeResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/grades [post]
// @Security BearerAuth
func (ctl *GradeController) CreateGrade(c *gin.Context) {
    var req models.GradeInput
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input: " + err.Error(),
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
    
    // Create grade
    grade, err := ctl.gradeService.CreateGradeWithTx(tx, req)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
  
        // Create grade model for activity logging
        gradeModel := models.Grade{
            ID:           grade.ID,
            StudentID:    grade.StudentID,
            CourseID:     grade.CourseID,
            TutorID:      grade.TutorID,
            AssignmentID: grade.AssignmentID,
            Score:        grade.Score,
            Grade:        grade.Grade,
            Remarks:      grade.Remarks,
        }

        _ = ctl.activity.Grades.Created(tx, req.TutorID, gradeModel)
    
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save grade: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Grade created successfully",
        "data":    grade,
    })
}


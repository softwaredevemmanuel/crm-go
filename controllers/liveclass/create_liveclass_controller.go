// controllers/live_class_controller.go
package controllers

import (
	"fmt"
	"net/http"

	"crm-go/models"
	"crm-go/services/activity"
	"crm-go/services/liveclass"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LiveClassController struct {
    db              *gorm.DB
    liveClassService *services.LiveClassService
    activity        *activity.Service
}

func NewLiveClassController(db *gorm.DB, liveClassService *services.LiveClassService, activitySvc *activity.Service) *LiveClassController {
    return &LiveClassController{
        db:              db,
        liveClassService: liveClassService,
        activity:        activitySvc,
    }
}

// CreateLiveClass handler
// @Summary Create a new live class
// @Description Create a scheduled live class with meeting setup
// @Tags live-classes
// @Accept json
// @Produce json
// @Param live_class body models.LiveClassInput true "Live class data"
// @Success 201 {object} models.LiveClassResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/live-classes [post]
// @Security BearerAuth
func (ctl *LiveClassController) CreateLiveClass(c *gin.Context) {

    var req models.LiveClassInput
    
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

    
    // Create live class
    liveClass, err := ctl.liveClassService.CreateLiveClassWithTx(tx, req)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    fmt.Println("here 2")

    // Log activity
        liveClassModel := models.LiveClass{
            ID:        liveClass.ID,
            CourseID:  liveClass.CourseID,
            Title:     liveClass.Title,
            StartTime: liveClass.StartTime,
            EndTime:   liveClass.EndTime,
            TutorID:   liveClass.TutorID,
        }
        
        _ = ctl.activity.LiveClasses.Created(tx, req.TutorID, liveClassModel)
    
    
    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save live class: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Live class created successfully",
        "data":    liveClass,
    })
}
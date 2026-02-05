// controllers/live_class_controller.go
package controllers

import (
	"crm-go/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UpdateLiveClass handler
// @Summary Update a live class
// @Description Update live class details. Some fields can only be updated before class starts.
// @Tags live-classes
// @Accept json
// @Produce json
// @Param id path string true "Live Class ID"
// @Param update body models.LiveClassUpdateInput true "Update data"
// @Success 200 {object} models.LiveClassResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 403 {object} models.ErrorResponse
// @Router /api/live-classes/{id} [put]
// @Security BearerAuth
// controllers/live_class_controller.go
// controllers/live_class_controller.go
func (ctl *LiveClassController) UpdateLiveClass(c *gin.Context) {
    liveClassID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid live class ID"})
        return
    }
    
    var req models.LiveClassUpdateInput
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
    
    // Update live class
    updatedClass, err := ctl.liveClassService.UpdateLiveClass(liveClassID, req, *req.TutorID)
    if err != nil {
        tx.Rollback()
        
        status := http.StatusBadRequest
        if strings.Contains(err.Error(), "cannot update") || 
           strings.Contains(err.Error(), "cannot cancel") ||
           strings.Contains(err.Error(), "cannot reduce") {
            status = http.StatusForbidden
        }
        
        c.JSON(status, gin.H{"error": err.Error()})
        return
    }
    
    // Log activity
        // Your activity logging here
        liveClassModel := models.LiveClass{
            ID:        updatedClass.ID,
            CourseID:  updatedClass.CourseID,
            Title:     updatedClass.Title,
            TutorID:   updatedClass.TutorID,
        }
        _ = ctl.activity.LiveClasses.Updated(tx, *req.TutorID, liveClassModel)

    // Commit transaction
    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save changes: " + err.Error(),
        })
        return
    }
    
    message := "Live class updated successfully"
    if req.IsCancelled != nil && *req.IsCancelled {
        message = "Live class cancelled successfully"
    } else if req.IsCancelled != nil && !*req.IsCancelled {
        message = "Live class uncancelled successfully"
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": message,
        "data":    updatedClass,
    })
}

// CancelLiveClass handler (convenience endpoint)
// @Summary Cancel a live class
// @Description Cancel a live class that hasn't started yet
// @Tags live-classes
// @Accept json
// @Produce json
// @Param id path string true "Live Class ID"
// @Param with_details query boolean false "Include course details"
// @Success 200 {object} models.LiveClassResponse
// @Router /api/live-classes/{id}/cancel [post]
// @Security BearerAuth
func (ctl *LiveClassController) CancelLiveClass(c *gin.Context) {
    liveClassID, err := uuid.Parse(c.Param("id"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid live class ID",
        })
        return
    }
    withDetails := false
    if details := c.Query("with_details"); details != "" {
        if val, err := strconv.ParseBool(details); err == nil {
            withDetails = val
        }
    }
    
    req := models.LiveClassUpdateInput{
        IsCancelled: &[]bool{true}[0],
    }
    
    updatedBy := uuid.Nil
    if userID, exists := c.Get("tutor_id"); exists {
        if id, ok := userID.(uuid.UUID); ok {
            updatedBy = id
        }
    }
    
    tx := ctl.db.Begin()
    
    updatedClass, err := ctl.liveClassService.UpdateLiveClassWithTx(tx, liveClassID, req, updatedBy, withDetails)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Log activity
        liveClassModel := models.LiveClass{
            ID:        updatedClass.ID,
            CourseID:  updatedClass.CourseID,
            Title:     updatedClass.Title,
            TutorID:   updatedClass.TutorID,
        }
        _ = ctl.activity.LiveClasses.Updated(tx, liveClassModel.TutorID, liveClassModel)

        fmt.Println("Live class updated:", liveClassModel.TutorID)

    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to cancel class: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusOK, gin.H{
        "message": "Live class cancelled successfully",
        "data":    updatedClass,
    })
}
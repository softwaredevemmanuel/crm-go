// controllers/objective_question_controller.go
package controllers

import (
    "net/http"
    
    "crm-go/models"
    "crm-go/services/objective_questions"
    "crm-go/services/activity"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "github.com/google/uuid"
)

type ObjectiveQuestionController struct {
    db               *gorm.DB
    questionService  *services.ObjectiveQuestionService
    activity        *activity.Service

}

func NewObjectiveQuestionController(db *gorm.DB, questionService *services.ObjectiveQuestionService, activitySvc *activity.Service) *ObjectiveQuestionController {
    return &ObjectiveQuestionController{
        db:              db,
        questionService: questionService,
        activity:    activitySvc,
    }
}

// CreateObjectiveQuestion handler
// @Summary Create a new objective question
// @Description Create objective question with options
// @Tags questions
// @Accept json
// @Produce json
// @Param question body models.ObjectiveQuestionInput true "Question data"
// @Success 201 {object} models.ObjectiveQuestionResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/questions/objective [post]
// @Security BearerAuth
func (ctl *ObjectiveQuestionController) CreateObjectiveQuestion(c *gin.Context) {
    var req models.ObjectiveQuestionInput
    
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input: " + err.Error(),
        })
        return
    }
    
    // Get user ID from context (creator)
    if userID, exists := c.Get("user_id"); exists {
        if id, ok := userID.(uuid.UUID); ok {
            req.TutorID = id
        }
    }
    
    // Set default approval based on user role
    // (You can implement role-based logic here)
    if req.IsApproved == false {
        // Default to false for non-admin users
        req.IsApproved = false
    }
    
    tx := ctl.db.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
        }
    }()
    
    // Create question
    question, err := ctl.questionService.CreateObjectiveQuestionWithTx(tx, req)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    // Log activity
    objectiveQuestionModel := models.ObjectiveQuestion{
        ID:        question.ID,
        CourseID:  question.CourseID,
    }
    _ = ctl.activity.ObjectiveQuestions.Created(tx, req.TutorID, objectiveQuestionModel)

    if err := tx.Commit().Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to save question: " + err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Question created successfully",
        "data":    question,
    })
}

// CreateBulkQuestions handler
// @Summary Create multiple questions
// @Description Create multiple questions in bulk
// @Tags questions
// @Accept json
// @Produce json
// @Param questions body []models.ObjectiveQuestionInput true "Array of question data"
// @Success 201 {array} models.ObjectiveQuestionResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/questions/objective/bulk [post]
// @Security BearerAuth
func (ctl *ObjectiveQuestionController) CreateBulkQuestions(c *gin.Context) {
    var requests []models.ObjectiveQuestionInput
    
    if err := c.ShouldBindJSON(&requests); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid input: " + err.Error(),
        })
        return
    }
    
    if len(requests) == 0 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "No questions provided",
        })
        return
    }
    
    // Limit bulk operations
    if len(requests) > 100 {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "Cannot create more than 100 questions at once",
        })
        return
    }
    
    // Set creator and defaults for each request
    if userID, exists := c.Get("user_id"); exists {
        if id, ok := userID.(uuid.UUID); ok {
            for i := range requests {
                requests[i].TutorID = id
                if !requests[i].IsApproved {
                    requests[i].IsApproved = false
                }
            }
        }
    }
    
    questions, err := ctl.questionService.CreateBulkQuestions(requests)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }
    
    c.JSON(http.StatusCreated, gin.H{
        "message": "Questions created successfully",
        "data":    questions,
        "count":   len(questions),
    })
}
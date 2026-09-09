// controllers/student_answer_controller.go
package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"crm-go/dto"
	"crm-go/services/objective_question_answers"
)

type StudentAnswerController struct {
	answerService *services.StudentAnswerService
}

func NewStudentAnswerController(answerService *services.StudentAnswerService) *StudentAnswerController {
	return &StudentAnswerController{
		answerService: answerService,
	}
}

// SubmitAnswer handles submitting an answer to a question
// @Summary Submit answer
// @Description Submit a student's answer to an objective question
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param request body dto.SubmitAnswerRequest true "Answer data"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-answers/submit [post]
func (c *StudentAnswerController) SubmitAnswer(ctx *gin.Context) {
	// Get student ID from context
	studentIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	studentID, err := uuid.Parse(studentIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var req dto.SubmitAnswerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"details": err.Error(),
		})
		return
	}

	response, err := c.answerService.SubmitAnswer(studentID, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not active") {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Answer submitted successfully",
		"data":    response,
	})
}

// GetStudentAnswers retrieves all answers for a student
// @Summary Get student answers
// @Description Get all answers submitted by a student
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param question_id query string false "Question ID"
// @Param is_correct query bool false "Filter by correctness"
// @Param status query string false "Status"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-answers [get]
func (c *StudentAnswerController) GetStudentAnswers(ctx *gin.Context) {
	// Get student ID from context
	studentIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	studentID, err := uuid.Parse(studentIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	var params dto.StudentAnswerQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	answers, err := c.answerService.GetStudentAnswers(studentID, &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Answers retrieved successfully",
		"data":    answers,
	})
}

// GetStudentStats retrieves statistics for a student
// @Summary Get student statistics
// @Description Get statistics for a student's answers
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param lesson_id query string false "Lesson ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-answers/stats [get]
func (c *StudentAnswerController) GetStudentStats(ctx *gin.Context) {
	// Get student ID from context
	studentIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	studentID, err := uuid.Parse(studentIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	lessonID := ctx.Query("lesson_id")

	stats, err := c.answerService.GetStudentStats(studentID, lessonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Statistics retrieved successfully",
		"data":    stats,
	})
}

// GetQuestionAttempts retrieves all attempts for a specific question
// @Summary Get question attempts
// @Description Get all attempts for a specific question by a student
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param question_id path string true "Question ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /api/student-answers/question/{question_id} [get]
func (c *StudentAnswerController) GetQuestionAttempts(ctx *gin.Context) {
	// Get student ID from context
	studentIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		return
	}

	studentID, err := uuid.Parse(studentIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	questionID := ctx.Param("question_id")
	if questionID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Question ID is required",
		})
		return
	}

	answers, err := c.answerService.GetQuestionAttempts(studentID, questionID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Attempts retrieved successfully",
		"data":    answers,
	})
}
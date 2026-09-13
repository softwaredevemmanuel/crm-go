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
// @Description Submit a student's answer to an objective question. If no option is selected, the answer is marked as wrong.
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param request body dto.SubmitAnswerRequest true "Answer data"
// @Success 200 {object} map[string]interface{} "Answer submitted successfully with result"
// @Failure 400 {object} map[string]interface{} "Invalid request or question not active"
// @Failure 401 {object} map[string]interface{} "Unauthorized - Invalid or missing token"
// @Failure 404 {object} map[string]interface{} "Question or grade not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/student-answers/submit [post]
func (c *StudentAnswerController) SubmitAnswer(ctx *gin.Context) {
	studentIDStr, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	studentID, err := uuid.Parse(studentIDStr.(string))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
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
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if strings.Contains(err.Error(), "not active") {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Answer submitted successfully",
		"data":    response,
	})
}

// GetStudentAnswers retrieves all answers with filters
// @Summary Get student answers
// @Description Get all answers with optional filters. Teachers can view all students' answers, students can view their own.
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param student_id query string false "Filter by Student ID"
// @Param lesson_id query string false "Filter by Lesson ID"
// @Param grade_id query string false "Filter by Grade ID"
// @Param question_id query string false "Filter by Question ID"
// @Param is_correct query bool false "Filter by correctness"
// @Param status query string false "Filter by status" Enums(active, inactive)
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20)
// @Param sort_by query string false "Sort field" default(created_at)
// @Param sort_order query string false "Sort direction" default(desc) Enums(asc, desc)
// @Success 200 {object} map[string]interface{} "Answers retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid query parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/fetch/student-answers [get]
func (c *StudentAnswerController) GetStudentAnswers(ctx *gin.Context) {


	var params dto.StudentAnswerQueryParams
	if err := ctx.ShouldBindQuery(&params); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid query parameters",
			"details": err.Error(),
		})
		return
	}

	answers, err := c.answerService.GetStudentAnswers( &params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Answers retrieved successfully",
		"data":    answers,
	})
}

// GetStudentStats retrieves statistics for a student
// @Summary Get student statistics
// @Description Get statistics for a student's answers, optionally filtered by lesson and grade
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param lesson_id query string false "Filter by Lesson ID"
// @Param grade_id query string false "Filter by Grade ID"
// @Success 200 {object} map[string]interface{} "Statistics retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Invalid parameters"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/student-answers/stats [get]
func (c *StudentAnswerController) GetStudentStats(ctx *gin.Context) {

	lessonID := ctx.Query("lesson_id")
	gradeID := ctx.Query("grade_id")

	stats, err := c.answerService.GetStudentStats(lessonID, gradeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Statistics retrieved successfully",
		"data":    stats,
	})
}

// GetLessonTestResults handles fetching test results for a lesson
// @Summary Get lesson test results
// @Description Get aggregated test results for all students in a specific lesson. Only accessible by teachers and admins.
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param lesson_id path string true "Lesson ID"
// @Success 200 {object} map[string]interface{} "Lesson results retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Lesson ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Not authorized to view all test results"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/student-answers/lesson/{lesson_id}/results [get]
func (c *StudentAnswerController) GetLessonTestResults(ctx *gin.Context) {
	lessonID := ctx.Param("lesson_id")
	if lessonID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Lesson ID is required"})
		return
	}

	// Check if user is a teacher or admin
	userRole, exists := ctx.Get("user_role")
	if exists && userRole != "admin" && userRole != "staff" && userRole != "teacher" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to view all test results",
		})
		return
	}

	response, err := c.answerService.GetLessonTestResults(lessonID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Lesson results retrieved successfully",
		"data":    response,
	})
}

// GetGradeTestResults handles fetching test results for a grade
// @Summary Get grade test results
// @Description Get aggregated test results for all students in a specific grade. Only accessible by teachers and admins.
// @Tags Objective Question Answers
// @Accept json
// @Produce json
// @Param grade_id path string true "Grade ID"
// @Success 200 {object} map[string]interface{} "Grade results retrieved successfully"
// @Failure 400 {object} map[string]interface{} "Grade ID is required"
// @Failure 401 {object} map[string]interface{} "Unauthorized"
// @Failure 403 {object} map[string]interface{} "Forbidden - Not authorized to view all test results"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Security BearerAuth
// @Router /api/student-answers/grade/{grade_id}/results [get]
func (c *StudentAnswerController) GetGradeTestResults(ctx *gin.Context) {
	gradeID := ctx.Param("grade_id")
	if gradeID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Grade ID is required"})
		return
	}

	// Check if user is a teacher or admin
	userRole, exists := ctx.Get("user_role")
	if exists && userRole != "admin" && userRole != "staff" && userRole != "teacher" {
		ctx.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to view all test results",
		})
		return
	}

	response, err := c.answerService.GetGradeTestResults(gradeID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Grade results retrieved successfully",
		"data":    response,
	})
}
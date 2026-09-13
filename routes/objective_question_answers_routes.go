// routes/student_answer_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/objective_question_answers"
	"crm-go/middleware"
	"crm-go/services/objective_question_answers"
)

func ObjectiveQuestionAnswersRoutes(router *gin.RouterGroup, db *gorm.DB) {
	answerService := services.NewStudentAnswerService(db)
	answerController := controllers.NewStudentAnswerController(answerService)

	answerGroup := router.Group("/api")
	answerGroup.Use(middleware.AuthMiddleware())
	{
		// Submit answer
		answerGroup.POST("/student-answers/submit", answerController.SubmitAnswer)

		// Get answers (with filters)
		answerGroup.GET("/fetch/student-answers", answerController.GetStudentAnswers)
		
		// Get statistics
		answerGroup.GET("/student-answers/stats", answerController.GetStudentStats)
		
		// Get lesson test results (for teachers/admins)
		answerGroup.GET("/student-answers/lesson/:lesson_id/results", answerController.GetLessonTestResults)
		
		// Get grade test results (for teachers/admins)
		answerGroup.GET("/student-answers/grade/:grade_id/results", answerController.GetGradeTestResults)
	}
}
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

		// Get answers
		answerGroup.GET("/student-answers", answerController.GetStudentAnswers)
		answerGroup.GET("/student-answers/stats", answerController.GetStudentStats)
		answerGroup.GET("/student-answers/question/:question_id", answerController.GetQuestionAttempts)
	}
}
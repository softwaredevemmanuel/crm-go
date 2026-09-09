// routes/objective_question_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/objective_question"
	"crm-go/middleware"
	"crm-go/services/objective_question"
)

func ObjectiveQuestionRoutes(router *gin.RouterGroup, db *gorm.DB) {
	questionService := services.NewObjectiveQuestionService(db)
	questionController := controllers.NewObjectiveQuestionController(questionService)

	questionGroup := router.Group("/api")
	questionGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		questionGroup.POST("/objective-questions", questionController.CreateQuestion)

		// READ
		questionGroup.GET("/objective-questions", questionController.GetQuestions)
		questionGroup.GET("/objective-questions/lesson/:lesson_id", questionController.GetQuestionsByLesson)
		questionGroup.GET("/objective-questions/topic/:topic_id", questionController.GetQuestionsByTopic)
		questionGroup.GET("/objective-questions/module/:module_id", questionController.GetQuestionsByModule)
		questionGroup.GET("/objective-questions/:id", questionController.GetQuestionByID)

		// UPDATE
		questionGroup.PUT("/objective-questions/:id", questionController.UpdateQuestion)

		// DELETE
		questionGroup.DELETE("/objective-questions/:id", questionController.DeleteQuestion)
	}
}
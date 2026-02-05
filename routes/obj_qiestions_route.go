// routes/live_class_routes.go
package routes

import (
    "crm-go/controllers/objective_questions"
    "crm-go/services/objective_questions"
    "crm-go/services/activity"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func ObjectiveQuestionRoutes(r *gin.Engine, db *gorm.DB) {
	questionService := services.NewObjectiveQuestionService(db)
	activityService := activity.NewService(db)
    questionController := controllers.NewObjectiveQuestionController(db, questionService, activityService)

    questionRoutes := r.Group("/api/questions/objective")
    {
        questionRoutes.POST("", questionController.CreateObjectiveQuestion)
        questionRoutes.POST("bulk", questionController.CreateBulkQuestions)

    }
    
}
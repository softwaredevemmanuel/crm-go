// routes/Lesson_report_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/lesson_report"
	"crm-go/middleware"
	"crm-go/services/lesson_report"
)

func LessonReportRoutes(router *gin.RouterGroup, db *gorm.DB) {
	reportService := services.NewLessonReportService(db)
	reportController := controllers.NewLessonReportController(reportService)

	reportGroup := router.Group("/api")
	reportGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		reportGroup.POST("/lesson-reports", reportController.CreateReport)

		// READ
		reportGroup.GET("/fetch-lesson-reports", reportController.GetReports)
		reportGroup.GET("/lesson-reports/stats", reportController.GetReportStats)
		reportGroup.GET("/lesson-reports/lesson/:lesson_id", reportController.GetReportsByLesson)
		reportGroup.GET("/lesson-reports/:id", reportController.GetReportByID)

		// UPDATE
		reportGroup.PUT("/lesson-reports/:id", reportController.UpdateReport)

		// DELETE
		reportGroup.DELETE("/lesson-reports/:id", reportController.DeleteReport)
	}
}
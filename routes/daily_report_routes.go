// routes/daily_report_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/daily_report"
	"crm-go/middleware"
	"crm-go/services/daily_report"
)

func DailyReportRoutes(router *gin.RouterGroup, db *gorm.DB) {
	reportService := services.NewDailyReportService(db)
	reportController := controllers.NewDailyReportController(reportService)

	reportGroup := router.Group("/api")
	reportGroup.Use(middleware.AuthMiddleware())
	{
		// CREATE
		reportGroup.POST("/daily-reports", reportController.CreateReport)

		// READ
		reportGroup.GET("/daily-reports", reportController.GetReports)
		reportGroup.GET("/daily-reports/stats", reportController.GetReportStats)
		reportGroup.GET("/daily-reports/lesson/:lesson_id", reportController.GetReportsByLesson)
		reportGroup.GET("/daily-reports/:id", reportController.GetReportByID)

		// UPDATE
		reportGroup.PUT("/daily-reports/:id", reportController.UpdateReport)

		// DELETE
		reportGroup.DELETE("/daily-reports/:id", reportController.DeleteReport)
	}
}
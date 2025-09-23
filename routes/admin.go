package routes

import (
	"crm-go/controllers/admin"
	"github.com/gin-gonic/gin"
	"crm-go/middleware"
)

func AdminDangerRoutes(r *gin.Engine) {
		adminGroup := r.Group("/admin")
	
	{
		// 🚨 Dangerous endpoint
		adminGroup.DELETE("/clear-db", admin.ClearDatabaseHandler)

		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.Use(middleware.RoleMiddleware("admin"))	
		protected.GET("/export/excel", admin.ExportExcelHandler)

	}
}

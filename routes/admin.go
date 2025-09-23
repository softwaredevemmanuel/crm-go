package routes

import (
	"crm-go/controllers/admin"
	"github.com/gin-gonic/gin"
	"crm-go/middleware"
)

func AdminDangerRoutes(r *gin.Engine) {
	{
		// 🚨 Dangerous endpoint
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.Use(middleware.RoleMiddleware("admin"))	
		protected.DELETE("/clear-db", admin.ClearDatabaseHandler)
		protected.GET("/export/excel", admin.ExportExcelHandler)

	}
}

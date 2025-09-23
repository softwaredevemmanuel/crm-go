package routes

import (
	"crm-go/controllers/admin"
	"github.com/gin-gonic/gin"
)

func AdminDangerRoutes(r *gin.Engine) {
	adminGroup := r.Group("/admin")
	{
		// 🚨 Dangerous endpoint
		adminGroup.DELETE("/clear-db", admin.ClearDatabaseHandler)
	}
}

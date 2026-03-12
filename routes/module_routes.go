package routes

import (
	modules "crm-go/controllers/modules"
	"crm-go/middleware"

	"github.com/gin-gonic/gin"
)

func ModuleRoutes(r *gin.Engine) {
	r.GET("/modules", modules.GetAllModules)
	r.GET("/modules/:id", modules.GetModuleByID)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	protected.POST("/modules", middleware.RoleMiddleware("admin"), modules.CreateModule)
	protected.PUT("/modules/:id", middleware.RoleMiddleware("admin"), modules.UpdateModule)
	protected.DELETE("/modules/:id", middleware.RoleMiddleware("admin"), modules.DeleteModule)
}

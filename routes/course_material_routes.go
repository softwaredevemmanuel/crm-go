package routes

import (
	courseMaterialController "crm-go/controllers/coursematerial" // alias for course materials
	"crm-go/middleware"
	"github.com/gin-gonic/gin"
)

func CourseMaterialRoutes(r *gin.Engine) {
	{
		courseMaterials := r.Group("/course-materials")
		courseMaterials.GET("/", courseMaterialController.GetCourseMaterials)
		courseMaterials.GET("/:id", courseMaterialController.GetCourseMaterialByID)

		// Protected routes
		protected := r.Group("/api")
		protected.Use(middleware.AuthMiddleware())
		protected.POST("/course-materials", middleware.RoleMiddleware("admin"), courseMaterialController.CreateCourseMaterial)
		protected.PUT("/course-materials/:id", middleware.RoleMiddleware("admin"), courseMaterialController.UpdateCourseMaterial)
		protected.DELETE("/course-materials/:id", middleware.RoleMiddleware("admin"), courseMaterialController.DeleteCourseMaterial)

	}
}

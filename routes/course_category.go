package routes

import (
    "crm-go/controllers/courses"
    "crm-go/middleware"
    "github.com/gin-gonic/gin"
)

func CourseCategoryRoutes(r *gin.Engine) {
    {
        // Protected routes
        protected := r.Group("/api")
        protected.Use(middleware.AuthMiddleware())
        protected.POST("/category-courses", middleware.RoleMiddleware("admin"), controllers.CreateCourseCategory)
        protected.DELETE("/category-courses/:id", middleware.RoleMiddleware("admin"), controllers.DeleteCourseCategory)

    }
}
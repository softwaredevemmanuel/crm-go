package routes

import (
	"crm-go/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterCourseRoutes(r *gin.Engine) {
	courses := r.Group("/courses")
	{
		courses.GET("/", controllers.GetCourses)
		courses.GET("/:id", controllers.GetCourse)

	}
}

// routes/department_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/department"
	"crm-go/services/department"
	"crm-go/middleware"
)

func DepartmentRoutes(router *gin.RouterGroup, db *gorm.DB) {
	service := services.NewDepartmentService(db)
	handler := controllers.NewDepartmentHandler(service)

	departmentGroup := router.Group("/api/departments")
	departmentGroup.Use(middleware.AuthMiddleware()) // Your auth middleware
	{
		// 1. Create department
		departmentGroup.POST("", handler.CreateDepartment)

		// 2. Get all departments with pagination
		departmentGroup.GET("", handler.GetAllDepartments)

		// 3. Get department with subjects
		departmentGroup.GET("/:id/subjects", handler.GetDepartmentWithSubjects)

		// 4. Get department with head and subjects
		departmentGroup.GET("/:id/head-subjects", handler.GetDepartmentWithHeadAndSubjects)

		// 5. Get department by ID with all details
		departmentGroup.GET("/:id", handler.GetDepartmentByID)

		// 6. Update department
		departmentGroup.PUT("/:id", handler.UpdateDepartment)

		// 7. Delete department
		departmentGroup.DELETE("/:id", handler.DeleteDepartment)
	}
}
// routes/upload_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers"
	"crm-go/middleware"
	"crm-go/services"
)

func UploadRoutes(router *gin.RouterGroup, db *gorm.DB) {
	uploadService := services.NewUploadService()
	uploadController := controllers.NewUploadController(uploadService)

	uploadGroup := router.Group("/api")
	uploadGroup.Use(middleware.AuthMiddleware())
	{
		// Single image upload (form-data: image)
		uploadGroup.POST("/upload/image", uploadController.UploadImage)

		// Multiple image upload (form-data: images[])
		uploadGroup.POST("/upload/images", uploadController.UploadMultipleImages)

		// Upload from remote URL (JSON: { url, caption })
		uploadGroup.POST("/upload/image/url", uploadController.UploadFromURL)

		// Delete image by public ID
		uploadGroup.DELETE("/upload/image/:public_id", uploadController.DeleteImage)
	}
}
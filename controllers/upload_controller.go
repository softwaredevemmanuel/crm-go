// controllers/upload_controller.go
package controllers

import (
	"log"
	"net/http"

	"crm-go/services"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploadService *services.UploadService
}

func NewUploadController(uploadService *services.UploadService) *UploadController {
	return &UploadController{uploadService: uploadService}
}

// UploadImage handles single image upload
// POST /api/upload/image
func (ctrl *UploadController) UploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No image file provided",
		})
		return
	}

	result, err := ctrl.uploadService.UploadImage(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	log.Println("image uploaded successfully")
	log.Println("URL image uploaded successfully", result.URL)
	log.Println("result image uploaded successfully", result)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image uploaded successfully",
		"url":     result.URL, // convenience field for frontend
		"data":    result,
	})
}

// UploadMultipleImages handles multiple image uploads (max 10)
// POST /api/upload/images
func (ctrl *UploadController) UploadMultipleImages(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Failed to parse multipart form",
		})
		return
	}

	files := form.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "No image files provided",
		})
		return
	}

	results, err := ctrl.uploadService.UploadMultipleImages(files)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Images uploaded successfully",
		"data": gin.H{
			"images": results,
			"count":  len(results),
		},
	})
}

// UploadFromURL handles uploading an image from a remote URL
// POST /api/upload/image/url
func (ctrl *UploadController) UploadFromURL(c *gin.Context) {
	var req struct {
		URL     string `json:"url" binding:"required"`
		Caption string `json:"caption"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Image URL is required",
		})
		return
	}

	result, err := ctrl.uploadService.UploadFromURL(req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image uploaded from URL successfully",
		"url":     result.URL,
		"data":    result,
	})
}

// DeleteImage deletes an image from Cloudinary
// DELETE /api/upload/image/:public_id
func (ctrl *UploadController) DeleteImage(c *gin.Context) {
	publicID := c.Param("public_id")
	if publicID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Public ID is required",
		})
		return
	}

	if err := ctrl.uploadService.DeleteImage(publicID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Image deleted successfully",
	})
}
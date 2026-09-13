// services/upload_service.go
package services

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"crm-go/config"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type UploadService struct{}

func NewUploadService() *UploadService {
	return &UploadService{}
}

// UploadResult holds the result of an image upload
type UploadResult struct {
	URL          string `json:"url"`
	PublicID     string `json:"public_id"`
	Format       string `json:"format"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Bytes        int64  `json:"bytes"`
	OriginalName string `json:"original_name"`
	CreatedAt    string `json:"created_at"`
}

const (
	maxFileSize = 5 * 1024 * 1024 // 5MB
	folder      = "inventory_items"
)

var allowedExtensions = map[string]bool{
	".jpg":  true,
	".jpeg": true,
	".png":  true,
	".gif":  true,
	".webp": true,
	".svg":  true,
}

// UploadImage uploads a single image file to Cloudinary
func (s *UploadService) UploadImage(file *multipart.FileHeader) (*UploadResult, error) {
	// Validate file size
	if file.Size > maxFileSize {
		return nil, fmt.Errorf("file too large. Maximum size is 5MB")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return nil, fmt.Errorf("invalid file type. Only JPEG, PNG, GIF, WEBP, and SVG are allowed")
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Determine folder (env override or default)
	uploadFolder := os.Getenv("CLOUDINARY_FOLDER")
	if uploadFolder == "" {
		uploadFolder = folder
	}

	// Generate unique public ID
	publicID := fmt.Sprintf("%d_%d", time.Now().Unix(), time.Now().UnixNano()%1e9)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// Upload to Cloudinary
	uploadParams := uploader.UploadParams{
		Folder:       uploadFolder,
		PublicID:     publicID,
		ResourceType: "image",
		Transformation: "c_limit,w_1200,h_1200/q_auto:good/f_auto",
	}

	result, err := config.CLD.Upload.Upload(ctx, src, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload to Cloudinary: %w", err)
	}
	log.Printf("Cloudinary raw response: %+v", result)


	return &UploadResult{
		URL:          result.SecureURL,
		PublicID:     result.PublicID,
		Format:       result.Format,
		Width:        result.Width,
		Height:       result.Height,
		Bytes:        int64(result.Bytes),
		OriginalName: file.Filename,
		CreatedAt:    time.Now().Format(time.RFC3339),
	}, nil
}

// UploadMultipleImages uploads multiple image files to Cloudinary
func (s *UploadService) UploadMultipleImages(files []*multipart.FileHeader) ([]*UploadResult, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files provided")
	}

	if len(files) > 10 {
		return nil, fmt.Errorf("too many files. Maximum is 10 files per request")
	}

	results := make([]*UploadResult, 0, len(files))
	for _, file := range files {
		result, err := s.UploadImage(file)
		if err != nil {
			return nil, fmt.Errorf("failed to upload %s: %w", file.Filename, err)
		}
		results = append(results, result)
	}

	return results, nil
}

// UploadFromURL uploads an image from a remote URL to Cloudinary
func (s *UploadService) UploadFromURL(imageURL string) (*UploadResult, error) {
	if imageURL == "" {
		return nil, fmt.Errorf("image URL is required")
	}

	uploadFolder := os.Getenv("CLOUDINARY_FOLDER")
	if uploadFolder == "" {
		uploadFolder = folder
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	result, err := config.CLD.Upload.Upload(ctx, imageURL, uploader.UploadParams{
		Folder:         uploadFolder,
		ResourceType:   "image",
		Transformation: "c_limit,w_1200,h_1200/q_auto:good/f_auto",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload from URL: %w", err)
	}

	return &UploadResult{
		URL:       result.SecureURL,
		PublicID:  result.PublicID,
		Format:    result.Format,
		Width:     result.Width,
		Height:    result.Height,
		Bytes:     int64(result.Bytes),
		CreatedAt: time.Now().Format(time.RFC3339),
	}, nil
}

// DeleteImage deletes an image from Cloudinary by its public ID
func (s *UploadService) DeleteImage(publicID string) error {
	if publicID == "" {
		return fmt.Errorf("public ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result, err := config.CLD.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID:     publicID,
		ResourceType: "image",
		Invalidate:   func(b bool) *bool { return &b }(true),
	})
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	if result.Result != "ok" && result.Result != "not found" {
		return fmt.Errorf("unexpected delete result: %s", result.Result)
	}

	return nil
}
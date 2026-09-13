// config/cloudinary.go
package config

import (
	"log"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CLD *cloudinary.Cloudinary

// InitCloudinary initializes the Cloudinary client
func InitCloudinary() {
	cloudName := os.Getenv("CLOUDINARY_CLOUD_NAME")
	apiKey := os.Getenv("CLOUDINARY_API_KEY")
	apiSecret := os.Getenv("CLOUDINARY_API_SECRET")

	if cloudName == "" || apiKey == "" || apiSecret == "" {
		log.Fatal("Cloudinary credentials are not set in environment variables")
	}

	cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	if err != nil {
		log.Fatalf("Failed to initialize Cloudinary: %v", err)
	}

	// Enable secure HTTPS URLs
	cld.Config.URL.Secure = true

	CLD = cld
	log.Println("✅ Cloudinary initialized successfully")
}
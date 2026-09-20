// services/wallet_service.go
package services

import (
	"net/http"
	"os"
	"strings"
	"time"

	"gorm.io/gorm"
)

type WalletService struct {
	db         *gorm.DB
	httpClient *http.Client
	secretKey  string
	baseURL    string
}

func NewWalletService(db *gorm.DB, secretKey string) *WalletService {
	baseURL := "https://api.flutterwave.com/v3"

	// Use sandbox in dev/test environments
	env := strings.ToLower(os.Getenv("FLUTTERWAVE_ENV"))
	if env == "sandbox" || env == "test" || env == "dev" {
		baseURL = "https://api.flutterwave.com/v3"
	}

	return &WalletService{
		db: db,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		secretKey: secretKey,
		baseURL:   baseURL,
	}
}
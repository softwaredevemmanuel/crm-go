// services/push_service.go
package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/models"
)

type PushSubscription struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID       uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Endpoint     string    `gorm:"type:text;not null" json:"endpoint"`
	P256dh       string    `gorm:"type:text;not null" json:"p256dh"`
	Auth         string    `gorm:"type:text;not null" json:"auth"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	User models.User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

func (PushSubscription) TableName() string {
	return "push_subscriptions"
}

// SubscriptionRequest represents the request to subscribe
type SubscriptionRequest struct {
	UserID       string `json:"user_id"`
	Subscription struct {
		Endpoint string `json:"endpoint"`
		Keys     struct {
			P256dh string `json:"p256dh"`
			Auth   string `json:"auth"`
		} `json:"keys"`
	} `json:"subscription"`
}

// PushPayload represents the notification payload
type PushPayload struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	Icon     string `json:"icon,omitempty"`
	Badge    string `json:"badge,omitempty"`
	URL      string `json:"url,omitempty"`
	Tag      string `json:"tag,omitempty"`
	Vibrate  []int  `json:"vibrate,omitempty"`
}

type PushService struct {
	db *gorm.DB
}

func NewPushService(db *gorm.DB) *PushService {
	return &PushService{db: db}
}

// Subscribe saves a push subscription
func (s *PushService) Subscribe(userID uuid.UUID, endpoint, p256dh, auth string) error {
	// Check if subscription already exists
	var existing PushSubscription
	if err := s.db.Where("user_id = ? AND endpoint = ?", userID, endpoint).First(&existing).Error; err == nil {
		// Update existing
		existing.P256dh = p256dh
		existing.Auth = auth
		existing.UpdatedAt = time.Now()
		return s.db.Save(&existing).Error
	}

	// Create new subscription
	subscription := &PushSubscription{
		ID:        uuid.New(),
		UserID:    userID,
		Endpoint:  endpoint,
		P256dh:    p256dh,
		Auth:      auth,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.db.Create(subscription).Error
}

// Unsubscribe removes a push subscription
func (s *PushService) Unsubscribe(userID uuid.UUID) error {
	return s.db.Where("user_id = ?", userID).Delete(&PushSubscription{}).Error
}

// GetUserSubscriptions gets all subscriptions for a user
func (s *PushService) GetUserSubscriptions(userID uuid.UUID) ([]PushSubscription, error) {
	var subscriptions []PushSubscription
	if err := s.db.Where("user_id = ? AND deleted_at IS NULL", userID).Find(&subscriptions).Error; err != nil {
		return nil, err
	}
	return subscriptions, nil
}

// SendPushNotification sends a push notification to a user
func (s *PushService) SendPushNotification(userID uuid.UUID, title, body, url string) error {
	// Get user's subscriptions
	subscriptions, err := s.GetUserSubscriptions(userID)
	if err != nil {
		return err
	}

	if len(subscriptions) == 0 {
		return errors.New("no push subscriptions found for user")
	}

	// Get VAPID keys from environment
	vapidPrivateKey := os.Getenv("VAPID_PRIVATE_KEY")
	vapidPublicKey := os.Getenv("VAPID_PUBLIC_KEY")
	if vapidPrivateKey == "" || vapidPublicKey == "" {
		return errors.New("VAPID keys not configured")
	}

	// Create payload
	payload := PushPayload{
		Title: title,
		Body:  body,
		Icon:  "/logo192.png",
		Badge: "/logo192.png",
		URL:   url,
		Tag:   "notification-" + uuid.New().String(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Send to each subscription
	for _, subscription := range subscriptions {
		if err := s.sendToEndpoint(&subscription, payloadBytes); err != nil {
			// Log error but continue with other subscriptions
			fmt.Printf("Failed to send to endpoint %s: %v\n", subscription.Endpoint, err)
		}
	}

	return nil
}

// sendToEndpoint sends a push notification to a specific endpoint
func (s *PushService) sendToEndpoint(subscription *PushSubscription, payload []byte) error {
	// Create request
	req, err := http.NewRequest("POST", subscription.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TTL", "86400") // 24 hours
	req.Header.Set("Urgency", "high")

	// Add VAPID headers
	// Note: This requires a proper VAPID implementation
	// You can use libraries like "github.com/SherClockHolmes/webpush-go"
	vapidHeader, err := s.generateVAPIDHeader(subscription)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", vapidHeader)

	// Send request
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("push service returned status: %s", resp.Status)
	}

	return nil
}

// generateVAPIDHeader generates the VAPID authorization header
func (s *PushService) generateVAPIDHeader(subscription *PushSubscription) (string, error) {
	// This is a placeholder - you'll need a proper VAPID implementation
	// Consider using: github.com/SherClockHolmes/webpush-go
	
	// For now, return a placeholder
	return "WebPush ...", nil
}

// SendBulkPushNotification sends push notifications to multiple users
func (s *PushService) SendBulkPushNotification(userIDs []uuid.UUID, title, body, url string) error {
	for _, userID := range userIDs {
		if err := s.SendPushNotification(userID, title, body, url); err != nil {
			fmt.Printf("Failed to send to user %s: %v\n", userID, err)
		}
	}
	return nil
}
// utils/session.go
package utils

import (
	"crm-go/config"
	"crm-go/models"
	"time"
)

func CleanExpiredSessions() {
	// Run this as a cron job to clean up expired sessions
	config.DB.Where("expires_at < ? OR is_active = false", time.Now()).
		Delete(&models.UserSession{})
}

func GetUserActiveSessions(userID string) ([]models.UserSession, error) {
	var sessions []models.UserSession
	err := config.DB.Where("user_id = ? AND is_active = true AND expires_at > ?", 
		userID, time.Now()).Find(&sessions).Error
	return sessions, err
}
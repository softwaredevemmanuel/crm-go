// services/notification_service.go
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/models"
	"crm-go/dto"
)

type NotificationService struct {
	db *gorm.DB
}

func NewNotificationService(db *gorm.DB) *NotificationService {
	return &NotificationService{db: db}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(req *dto.CreateNotificationRequest, userID uuid.UUID) (*dto.NotificationResponse, error) {
	// Validate input
	if err := s.validateNotificationRequest(req); err != nil {
		return nil, err
	}

	// Parse User ID
	userUUID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, errors.New("invalid user ID format")
	}

	// Check if user exists
	var user models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", userUUID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, errors.New("failed to verify user: " + err.Error())
	}

	// Parse expires_at if provided
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		exp, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return nil, errors.New("invalid expires_at format. Use RFC3339")
		}
		expiresAt = &exp
	}

	// Set default priority if not provided
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	// Create notification
	notification := &models.Notification{
		ID:          uuid.New(),
		UserID:      userUUID,
		Type:        req.Type,
		Title:       strings.TrimSpace(req.Title),
		Message:     strings.TrimSpace(req.Message),
		Icon:        req.Icon,
		Color:       req.Color,
		Link:        req.Link,
		Priority:    priority,
		SentAt:      time.Now(),
		ExpiresAt:   expiresAt,
		Metadata:    req.Metadata,
		CreatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.db.Create(notification).Error; err != nil {
		return nil, errors.New("failed to create notification: " + err.Error())
	}

	// Preload user for response
	if err := s.db.Preload("User").First(notification, notification.ID).Error; err != nil {
		return nil, errors.New("failed to load notification details: " + err.Error())
	}

	return s.toNotificationResponse(notification), nil
}

// BulkCreateNotifications creates multiple notifications
func (s *NotificationService) BulkCreateNotifications(req *dto.BulkCreateNotificationRequest, userID uuid.UUID) ([]dto.NotificationResponse, error) {
	// Validate input
	if len(req.UserIDs) == 0 {
		return nil, errors.New("at least one user ID is required")
	}

	// Parse expires_at if provided
	var expiresAt *time.Time
	if req.ExpiresAt != "" {
		exp, err := time.Parse(time.RFC3339, req.ExpiresAt)
		if err != nil {
			return nil, errors.New("invalid expires_at format. Use RFC3339")
		}
		expiresAt = &exp
	}

	// Set default priority if not provided
	priority := req.Priority
	if priority == "" {
		priority = "normal"
	}

	var created []dto.NotificationResponse
	var failed []string

	for _, userIDStr := range req.UserIDs {
		userUUID, err := uuid.Parse(userIDStr)
		if err != nil {
			failed = append(failed, userIDStr+" (invalid ID format)")
			continue
		}

		// Check if user exists
		var user models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", userUUID).First(&user).Error; err != nil {
			failed = append(failed, userIDStr+" (user not found)")
			continue
		}

		// Create notification
		notification := &models.Notification{
			ID:          uuid.New(),
			UserID:      userUUID,
			Type:        req.Type,
			Title:       strings.TrimSpace(req.Title),
			Message:     strings.TrimSpace(req.Message),
			Icon:        req.Icon,
			Color:       req.Color,
			Link:        req.Link,
			Priority:    priority,
			SentAt:      time.Now(),
			ExpiresAt:   expiresAt,
			Metadata:    req.Metadata,
			CreatedBy:   userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := s.db.Create(notification).Error; err != nil {
			failed = append(failed, userIDStr+" (failed to create)")
			continue
		}

		// Preload user for response
		s.db.Preload("User").First(notification, notification.ID)
		created = append(created, *s.toNotificationResponse(notification))
	}

	if len(failed) > 0 {
		return created, fmt.Errorf("bulk create completed with failures: %v", failed)
	}

	return created, nil
}

// GetAllNotifications retrieves all notifications with pagination and filters
func (s *NotificationService) GetAllNotifications(params *dto.NotificationQueryParams, userID uuid.UUID) (*dto.NotificationListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "sent_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND deleted_at IS NULL", userID)

	// Apply filters
	if params.Type != "" {
		query = query.Where("type = ?", params.Type)
	}

	if params.Priority != "" {
		query = query.Where("priority = ?", params.Priority)
	}

	if params.IsRead != nil {
		query = query.Where("is_read = ?", *params.IsRead)
	}

	if params.IsDismissed != nil {
		query = query.Where("is_dismissed = ?", *params.IsDismissed)
	}

	// Get unread count
	var unreadCount int64
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND deleted_at IS NULL", userID, false).
		Count(&unreadCount).Error; err != nil {
		return nil, errors.New("failed to count unread notifications: " + err.Error())
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count notifications: " + err.Error())
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order(params.SortBy + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query
	var notifications []models.Notification
	if err := query.Preload("User").Find(&notifications).Error; err != nil {
		return nil, errors.New("failed to fetch notifications: " + err.Error())
	}

	// Convert to response
	responses := make([]dto.NotificationResponse, len(notifications))
	for i, notification := range notifications {
		responses[i] = *s.toNotificationResponse(&notification)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.NotificationListResponse{
		Notifications: responses,
		UnreadCount:   unreadCount,
		Total:         total,
		Page:          params.Page,
		Limit:         params.Limit,
		TotalPages:    totalPages,
	}, nil
}

// GetNotificationByID retrieves a single notification by ID
func (s *NotificationService) GetNotificationByID(id string, userID uuid.UUID) (*dto.NotificationResponse, error) {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid notification ID")
	}

	var notification models.Notification
	if err := s.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", notificationID, userID).
		Preload("User").
		First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, errors.New("failed to fetch notification: " + err.Error())
	}

	return s.toNotificationResponse(&notification), nil
}

// UpdateNotification updates an existing notification
func (s *NotificationService) UpdateNotification(id string, req *dto.UpdateNotificationRequest, userID uuid.UUID) (*dto.NotificationResponse, error) {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid notification ID")
	}

	// Find existing notification
	var notification models.Notification
	if err := s.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", notificationID, userID).
		First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("notification not found")
		}
		return nil, errors.New("failed to fetch notification: " + err.Error())
	}

	// Update fields
	if req.IsRead != nil {
		notification.IsRead = *req.IsRead
		if *req.IsRead && notification.ReadAt == nil {
			now := time.Now()
			notification.ReadAt = &now
		}
	}

	if req.IsDismissed != nil {
		notification.IsDismissed = *req.IsDismissed
		if *req.IsDismissed && notification.DismissedAt == nil {
			now := time.Now()
			notification.DismissedAt = &now
		}
	}

	if req.Priority != "" {
		notification.Priority = req.Priority
	}

	notification.UpdatedAt = time.Now()

	if err := s.db.Save(&notification).Error; err != nil {
		return nil, errors.New("failed to update notification: " + err.Error())
	}

	// Preload user for response
	if err := s.db.Preload("User").First(&notification, notification.ID).Error; err != nil {
		return nil, errors.New("failed to load notification details: " + err.Error())
	}

	return s.toNotificationResponse(&notification), nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(id string, userID uuid.UUID) error {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid notification ID")
	}

	now := time.Now()
	result := s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", notificationID, userID).
		Updates(map[string]interface{}{
			"is_read":  true,
			"read_at":  now,
			"updated_at": now,
		})

	if result.Error != nil {
		return errors.New("failed to mark notification as read: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.New("notification not found")
	}

	return nil
}

// MarkAllAsRead marks all notifications as read for a user
func (s *NotificationService) MarkAllAsRead(userID uuid.UUID) error {
	now := time.Now()
	result := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND deleted_at IS NULL", userID, false).
		Updates(map[string]interface{}{
			"is_read":  true,
			"read_at":  now,
			"updated_at": now,
		})

	if result.Error != nil {
		return errors.New("failed to mark all as read: " + result.Error.Error())
	}

	return nil
}

// DismissNotification dismisses a notification
func (s *NotificationService) DismissNotification(id string, userID uuid.UUID) error {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid notification ID")
	}

	now := time.Now()
	result := s.db.Model(&models.Notification{}).
		Where("id = ? AND user_id = ? AND deleted_at IS NULL", notificationID, userID).
		Updates(map[string]interface{}{
			"is_dismissed": true,
			"dismissed_at": now,
			"updated_at":   now,
		})

	if result.Error != nil {
		return errors.New("failed to dismiss notification: " + result.Error.Error())
	}

	if result.RowsAffected == 0 {
		return errors.New("notification not found")
	}

	return nil
}

// GetNotificationStats returns notification statistics for a user
func (s *NotificationService) GetNotificationStats(userID uuid.UUID) (*dto.NotificationStatsResponse, error) {
	stats := &dto.NotificationStatsResponse{
		ByPriority: make(map[string]int64),
		ByType:     make(map[string]int64),
	}

	// Get total count
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Count(&stats.Total).Error; err != nil {
		return nil, errors.New("failed to count total notifications: " + err.Error())
	}

	// Get unread count
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND deleted_at IS NULL", userID, false).
		Count(&stats.Unread).Error; err != nil {
		return nil, errors.New("failed to count unread notifications: " + err.Error())
	}

	// Get read count
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_read = ? AND deleted_at IS NULL", userID, true).
		Count(&stats.Read).Error; err != nil {
		return nil, errors.New("failed to count read notifications: " + err.Error())
	}

	// Get dismissed count
	if err := s.db.Model(&models.Notification{}).
		Where("user_id = ? AND is_dismissed = ? AND deleted_at IS NULL", userID, true).
		Count(&stats.Dismissed).Error; err != nil {
		return nil, errors.New("failed to count dismissed notifications: " + err.Error())
	}

	// Get counts by priority
	var priorityResults []struct {
		Priority string
		Count    int64
	}
	if err := s.db.Model(&models.Notification{}).
		Select("priority, count(*) as count").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group("priority").
		Scan(&priorityResults).Error; err != nil {
		return nil, errors.New("failed to get counts by priority: " + err.Error())
	}
	for _, r := range priorityResults {
		stats.ByPriority[r.Priority] = r.Count
	}

	// Get counts by type
	var typeResults []struct {
		Type  string
		Count int64
	}
	if err := s.db.Model(&models.Notification{}).
		Select("type, count(*) as count").
		Where("user_id = ? AND deleted_at IS NULL", userID).
		Group("type").
		Scan(&typeResults).Error; err != nil {
		return nil, errors.New("failed to get counts by type: " + err.Error())
	}
	for _, r := range typeResults {
		stats.ByType[r.Type] = r.Count
	}

	return stats, nil
}

// DeleteNotification soft deletes a notification
func (s *NotificationService) DeleteNotification(id string, userID uuid.UUID) error {
	notificationID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid notification ID")
	}

	var notification models.Notification
	if err := s.db.Where("id = ? AND user_id = ? AND deleted_at IS NULL", notificationID, userID).
		First(&notification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("notification not found")
		}
		return errors.New("failed to fetch notification: " + err.Error())
	}

	if err := s.db.Delete(&notification).Error; err != nil {
		return errors.New("failed to delete notification: " + err.Error())
	}

	return nil
}

// CleanupExpiredNotifications removes expired notifications
func (s *NotificationService) CleanupExpiredNotifications() error {
	now := time.Now()
	result := s.db.Unscoped().
		Where("expires_at IS NOT NULL AND expires_at < ?", now).
		Delete(&models.Notification{})

	if result.Error != nil {
		return errors.New("failed to cleanup expired notifications: " + result.Error.Error())
	}

	if result.RowsAffected > 0 {
		fmt.Printf("🧹 Cleaned up %d expired notifications\n", result.RowsAffected)
	}

	return nil
}

// validateNotificationRequest validates the notification request
func (s *NotificationService) validateNotificationRequest(req *dto.CreateNotificationRequest) error {
	if req.UserID == "" {
		return errors.New("user ID is required")
	}
	if req.Type == "" {
		return errors.New("notification type is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if len(req.Title) < 3 {
		return errors.New("title must be at least 3 characters")
	}
	if req.Message == "" {
		return errors.New("message is required")
	}
	if req.Priority != "" && req.Priority != "low" && req.Priority != "normal" && req.Priority != "high" && req.Priority != "urgent" {
		return errors.New("priority must be 'low', 'normal', 'high', or 'urgent'")
	}
	return nil
}

// toNotificationResponse converts model to response DTO
func (s *NotificationService) toNotificationResponse(notification *models.Notification) *dto.NotificationResponse {
	response := &dto.NotificationResponse{
		ID:          notification.ID.String(),
		UserID:      notification.UserID.String(),
		Type:        notification.Type,
		Title:       notification.Title,
		Message:     notification.Message,
		Icon:        notification.Icon,
		Color:       notification.Color,
		Link:        notification.Link,
		IsRead:      notification.IsRead,
		IsDismissed: notification.IsDismissed,
		Priority:    notification.Priority,
		SentAt:      notification.SentAt,
		ReadAt:      notification.ReadAt,
		DismissedAt: notification.DismissedAt,
		ExpiresAt:   notification.ExpiresAt,
		Metadata:    notification.Metadata,
		CreatedBy:   notification.CreatedBy.String(),
		CreatedAt:   notification.CreatedAt,
		UpdatedAt:   notification.UpdatedAt,
	}

	// Add user details if preloaded
	if notification.User.ID != uuid.Nil {
		response.User = &dto.UserResponse{
			ID:        notification.User.ID.String(),
			FirstName: notification.User.FirstName,
			LastName:  notification.User.LastName,
			Email:     notification.User.Email,
			Phone:     notification.User.Phone,
			Role:      notification.User.Role,
		}
	}

	return response
}
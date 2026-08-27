package dto

import (
	"time"
)

// CreateNotificationRequest represents the request body for creating a notification
type CreateNotificationRequest struct {
	UserID   string `json:"user_id" binding:"required"`
	Type     string `json:"type" binding:"required,oneof=system announcement assignment grade enrollment message alert reminder"`
	Title    string `json:"title" binding:"required,min=3,max=255"`
	Message  string `json:"message" binding:"required"`
	Link     string `json:"link"`
	Priority string `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
	ExpiresAt string `json:"expires_at"`
	Metadata string `json:"metadata"`
}

// BulkCreateNotificationRequest represents the request body for bulk creating notifications
type BulkCreateNotificationRequest struct {
	UserIDs  []string `json:"user_ids" binding:"required,min=1"`
	Type     string   `json:"type" binding:"required,oneof=system announcement assignment grade enrollment message alert reminder"`
	Title    string   `json:"title" binding:"required,min=3,max=255"`
	Message  string   `json:"message" binding:"required"`
	Link     string   `json:"link"`
	Priority string   `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
	ExpiresAt string  `json:"expires_at"`
	Metadata string   `json:"metadata"`
}

// UpdateNotificationRequest represents the request body for updating a notification
type UpdateNotificationRequest struct {
	IsRead      *bool  `json:"is_read"`
	IsDismissed *bool  `json:"is_dismissed"`
	Priority    string `json:"priority" binding:"omitempty,oneof=low normal high urgent"`
}

// NotificationResponse represents the notification response
type NotificationResponse struct {
	ID          string     `json:"id"`
	UserID      string     `json:"user_id"`
	Type        string     `json:"type"`
	Title       string     `json:"title"`
	Message     string     `json:"message"`
	Link        string     `json:"link"`
	IsRead      bool       `json:"is_read"`
	IsDismissed bool       `json:"is_dismissed"`
	Priority    string     `json:"priority"`
	SentAt      time.Time  `json:"sent_at"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	DismissedAt *time.Time `json:"dismissed_at,omitempty"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
	Metadata    string     `json:"metadata"`
	CreatedBy   string     `json:"created_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	User        *UserResponse `json:"user,omitempty"`
}

// NotificationListResponse represents paginated notification list response
type NotificationListResponse struct {
	Notifications []NotificationResponse `json:"notifications"`
	UnreadCount   int64                  `json:"unread_count"`
	Total         int64                  `json:"total"`
	Page          int                    `json:"page"`
	Limit         int                    `json:"limit"`
	TotalPages    int                    `json:"total_pages"`
}

// NotificationQueryParams represents query parameters for filtering notifications
type NotificationQueryParams struct {
	Type       string `form:"type" binding:"omitempty,oneof=system announcement assignment grade enrollment message alert reminder"`
	Priority   string `form:"priority" binding:"omitempty,oneof=low normal high urgent"`
	IsRead     *bool  `form:"is_read"`
	IsDismissed *bool `form:"is_dismissed"`
	Page       int    `form:"page" binding:"min=1"`
	Limit      int    `form:"limit" binding:"min=1,max=100"`
	SortBy     string `form:"sort_by" binding:"omitempty,oneof=priority sent_at created_at"`
	SortOrder  string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// NotificationStatsResponse represents notification statistics
type NotificationStatsResponse struct {
	Total      int64            `json:"total"`
	Unread     int64            `json:"unread"`
	Read       int64            `json:"read"`
	Dismissed  int64            `json:"dismissed"`
	ByPriority map[string]int64 `json:"by_priority"`
	ByType     map[string]int64 `json:"by_type"`
}

// MarkReadRequest represents the request body for marking notifications as read
type MarkReadRequest struct {
	NotificationIDs []string `json:"notification_ids"`
	MarkAll         bool     `json:"mark_all"`
}
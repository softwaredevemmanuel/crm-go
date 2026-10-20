package dto

import (
	"time"

	"github.com/google/uuid"

	"crm-go/models"
)

// ============================================================
// Requests
// ============================================================

type SubscribeRequest struct {
	Email  string `json:"email" binding:"required,email,max=255" example:"jane@example.com"`
	Name   string `json:"name" binding:"omitempty,max=255" example:"Jane Doe"`
	Source string `json:"source" binding:"omitempty,oneof=homepage_footer blog course_page other" example:"homepage_footer"`
	Tags   string `json:"tags" binding:"omitempty,max=500" example:"courses,news"`
}

type ConfirmSubscriptionRequest struct {
	Token string `json:"token" binding:"required" example:"a1b2c3d4e5f6..."`
}

type UnsubscribeRequest struct {
	Token  string `json:"token" binding:"required"`
	Reason string `json:"reason" binding:"omitempty,max=255"`
}

type AdminCreateSubscriberRequest struct {
	Email  string `json:"email" binding:"required,email,max=255"`
	Name   string `json:"name" binding:"omitempty,max=255"`
	Status string `json:"status" binding:"omitempty,oneof=subscribed pending unsubscribed"`
	Tags   string `json:"tags" binding:"omitempty,max=500"`
}

type AdminUpdateSubscriberRequest struct {
	Name   *string `json:"name" binding:"omitempty,max=255"`
	Status *string `json:"status" binding:"omitempty,oneof=pending subscribed unsubscribed bounced complained"`
	Tags   *string `json:"tags" binding:"omitempty,max=500"`
}

// ============================================================
// Responses
// ============================================================

type SubscriberResponse struct {
	ID                uuid.UUID  `json:"id"`
	Email             string     `json:"email"`
	Name              string     `json:"name,omitempty"`
	Status            string     `json:"status"`
	Source            string     `json:"source"`
	Tags              string     `json:"tags,omitempty"`
	ConfirmedAt       *time.Time `json:"confirmed_at,omitempty"`
	UnsubscribedAt    *time.Time `json:"unsubscribed_at,omitempty"`
	UnsubscribeReason string     `json:"unsubscribe_reason,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

func FromSubscriberModel(s *models.Subscriber) SubscriberResponse {
	return SubscriberResponse{
		ID:                s.ID,
		Email:             s.Email,
		Name:              s.Name,
		Status:            string(s.Status),
		Source:            string(s.Source),
		Tags:              s.Tags,
		ConfirmedAt:       s.ConfirmedAt,
		UnsubscribedAt:    s.UnsubscribedAt,
		UnsubscribeReason: s.UnsubscribeReason,
		CreatedAt:         s.CreatedAt,
		UpdatedAt:         s.UpdatedAt,
	}
}

type ListSubscribersQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Search    string `form:"search"`
	Status    string `form:"status"`
	Source    string `form:"source"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}

type SingleSubscriberResponse struct {
	Data SubscriberResponse `json:"data"`
}

type PaginatedSubscribersResponse struct {
	Data       []SubscriberResponse `json:"data"`
	Total      int64                `json:"total"`
	Page       int                  `json:"page"`
	Limit      int                  `json:"limit"`
	TotalPages int                  `json:"total_pages"`
}

type SimpleMessage struct {
	Message string `json:"message" example:"Please check your email to confirm your subscription."`
}

type NewsletterStats struct {
	Total        int64 `json:"total"`
	Subscribed   int64 `json:"subscribed"`
	Pending      int64 `json:"pending"`
	Unsubscribed int64 `json:"unsubscribed"`
	Bounced      int64 `json:"bounced"`
	Last30Days   int64 `json:"last_30_days"`
}

// ErrorResponse is the standard error envelope for Swagger docs and API responses.
type ErrorResponse struct {
	Error   string            `json:"error" example:"validation failed"`
	Message string            `json:"message,omitempty" example:"one or more fields are invalid"`
	Fields  map[string]string `json:"fields,omitempty"`
}
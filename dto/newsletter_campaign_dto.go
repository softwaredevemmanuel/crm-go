package dto

import (
	"time"

	"github.com/google/uuid"

	"crm-go/models"
)

// CreateCampaignRequest — payload for creating a newsletter campaign.
type CreateCampaignRequest struct {
	Subject     string   `json:"subject" binding:"required,min=3,max=255" example:"New Course Launch: Intro to Go"`
	Title       string   `json:"title" binding:"omitempty,max=255" example:"Course Announcement"`
	Content     string   `json:"content" binding:"required" example:"<h1>Hello {{.Name}}</h1><p>...</p>"`
	Preview     string   `json:"preview" binding:"omitempty,max=500" example:"Here's what's new this week..."`
	Audience    string   `json:"audience" binding:"omitempty,oneof=all active tagged manual" example:"active"`
	AudienceTag string   `json:"audience_tag" binding:"omitempty,max=255" example:"courses"`
	ManualIDs   []string `json:"manual_ids" binding:"omitempty,dive,uuid4"`
	ScheduledAt *string  `json:"scheduled_at" binding:"omitempty" example:"2025-02-01T09:00:00Z"`
}

// UpdateCampaignRequest — partial update.
type UpdateCampaignRequest struct {
	Subject     *string   `json:"subject" binding:"omitempty,min=3,max=255"`
	Title       *string   `json:"title" binding:"omitempty,max=255"`
	Content     *string   `json:"content"`
	Preview     *string   `json:"preview" binding:"omitempty,max=500"`
	Audience    *string   `json:"audience" binding:"omitempty,oneof=all active tagged manual"`
	AudienceTag *string   `json:"audience_tag" binding:"omitempty,max=255"`
	ManualIDs   *[]string `json:"manual_ids" binding:"omitempty,dive,uuid4"`
	ScheduledAt *string   `json:"scheduled_at"`
}

// CampaignResponse — the shape returned to the client.
type CampaignResponse struct {
	ID             uuid.UUID  `json:"id"`
	Subject        string     `json:"subject"`
	Title          string     `json:"title"`
	Content        string     `json:"content"`
	Preview        string     `json:"preview"`
	Audience       string     `json:"audience"`
	AudienceTag    string     `json:"audience_tag,omitempty"`
	Status         string     `json:"status"`
	ScheduledAt    *time.Time `json:"scheduled_at,omitempty"`
	SentAt         *time.Time `json:"sent_at,omitempty"`
	FailureReason  string     `json:"failure_reason,omitempty"`
	RecipientCount int        `json:"recipient_count"`
	DeliveredCount int        `json:"delivered_count"`
	FailedCount    int        `json:"failed_count"`
	OpenedCount    int        `json:"opened_count"`
	ClickedCount   int        `json:"clicked_count"`
	CreatedBy      *uuid.UUID `json:"created_by,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func FromCampaignModel(c *models.NewsletterCampaign) CampaignResponse {
	return CampaignResponse{
		ID:             c.ID,
		Subject:        c.Subject,
		Title:          c.Title,
		Content:        c.Content,
		Preview:        c.Preview,
		Audience:       string(c.Audience),
		AudienceTag:    c.AudienceTag,
		Status:         string(c.Status),
		ScheduledAt:    c.ScheduledAt,
		SentAt:         c.SentAt,
		FailureReason:  c.FailureReason,
		RecipientCount: c.RecipientCount,
		DeliveredCount: c.DeliveredCount,
		FailedCount:    c.FailedCount,
		OpenedCount:    c.OpenedCount,
		ClickedCount:   c.ClickedCount,
		CreatedBy:      c.CreatedBy,
		CreatedAt:      c.CreatedAt,
		UpdatedAt:      c.UpdatedAt,
	}
}

// ListCampaignsQuery — filters + pagination.
type ListCampaignsQuery struct {
	Page      int    `form:"page,default=1" binding:"min=1"`
	Limit     int    `form:"limit,default=20" binding:"min=1,max=100"`
	Search    string `form:"search"`
	Status    string `form:"status"`
	SortBy    string `form:"sort_by,default=created_at"`
	SortOrder string `form:"sort_order,default=desc"`
}

// SingleCampaignResponse wraps one campaign.
type SingleCampaignResponse struct {
	Data CampaignResponse `json:"data"`
}

// PaginatedCampaignsResponse wraps a page of campaigns.
type PaginatedCampaignsResponse struct {
	Data       []CampaignResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	Limit      int                `json:"limit"`
	TotalPages int                `json:"total_pages"`
}

// SendCampaignResponse — returned when a campaign is queued/sent.
type SendCampaignResponse struct {
	Message        string `json:"message" example:"Campaign is being sent to 1,240 subscribers."`
	RecipientCount int    `json:"recipient_count" example:"1240"`
}


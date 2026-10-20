package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CampaignStatus string

const (
	CampaignDraft     CampaignStatus = "draft"
	CampaignScheduled CampaignStatus = "scheduled"
	CampaignSending   CampaignStatus = "sending"
	CampaignSent      CampaignStatus = "sent"
	CampaignFailed    CampaignStatus = "failed"
	CampaignCancelled CampaignStatus = "cancelled"
)

type CampaignAudience string

const (
	AudienceAll          CampaignAudience = "all"          // all subscribed
	AudienceActive       CampaignAudience = "active"       // subscribed + confirmed
	AudienceTagged       CampaignAudience = "tagged"       // filter by tags
	AudienceManual       CampaignAudience = "manual"       // specific IDs (stored in JSON)
)

type NewsletterCampaign struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Subject  string    `gorm:"type:varchar(255);not null" json:"subject" example:"New Course Launch: Intro to Go"`
	Title    string    `gorm:"type:varchar(255)" json:"title" example:"Course Announcement"`
	Content  string    `gorm:"type:text;not null" json:"content" example:"<h1>Hello...</h1>"`
	Preview  string    `gorm:"type:varchar(500)" json:"preview" example:"Here's what's new this week..."`

	// Audience targeting
	Audience    CampaignAudience `gorm:"type:varchar(20);not null;default:'active'" json:"audience" example:"active"`
	AudienceTag string           `gorm:"type:varchar(255)" json:"audience_tag,omitempty" example:"courses"`
	ManualIDs   string           `gorm:"type:text" json:"-"` // JSON array of UUIDs

	// Status & scheduling
	Status      CampaignStatus `gorm:"type:varchar(20);not null;default:'draft';index" json:"status" example:"draft"`
	ScheduledAt *time.Time     `json:"scheduled_at,omitempty"`
	SentAt      *time.Time     `json:"sent_at,omitempty"`
	FailureReason string       `gorm:"type:varchar(500)" json:"failure_reason,omitempty"`

	// Delivery metrics
	RecipientCount int `gorm:"default:0" json:"recipient_count" example:"0"`
	DeliveredCount int `gorm:"default:0" json:"delivered_count" example:"0"`
	FailedCount    int `gorm:"default:0" json:"failed_count" example:"0"`
	OpenedCount    int `gorm:"default:0" json:"opened_count" example:"0"`
	ClickedCount   int `gorm:"default:0" json:"clicked_count" example:"0"`

	// Author
	CreatedBy *uuid.UUID `gorm:"type:uuid" json:"created_by,omitempty"`

	CreatedAt time.Time      `json:"created_at" example:"2025-01-15T10:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2025-01-15T10:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (NewsletterCampaign) TableName() string {
	return "newsletter_campaigns"
}
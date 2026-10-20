package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SubscriberStatus string

const (
	SubscriberPending     SubscriberStatus = "pending"     // awaiting email confirmation
	SubscriberSubscribed  SubscriberStatus = "subscribed"  // confirmed
	SubscriberUnsubscribed SubscriberStatus = "unsubscribed"
	SubscriberBounced     SubscriberStatus = "bounced"     // email bounced
	SubscriberComplained  SubscriberStatus = "complained"  // marked as spam
)

type SubscriberSource string

const (
	SourceHomepageFooter SubscriberSource = "homepage_footer"
	SourceBlog           SubscriberSource = "blog"
	SourceCoursePage     SubscriberSource = "course_page"
	SourceAdmin          SubscriberSource = "admin"
	SourceImport         SubscriberSource = "import"
	SourceOther          SubscriberSource = "other"
)

type Subscriber struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Email  string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"email" example:"jane@example.com"`
	Name   string    `gorm:"type:varchar(255)" json:"name,omitempty" example:"Jane Doe"`
	Status SubscriberStatus `gorm:"type:varchar(20);not null;default:'pending';index" json:"status" example:"pending"`
	Source SubscriberSource `gorm:"type:varchar(30);default:'homepage_footer'" json:"source" example:"homepage_footer"`

	// Double opt-in
	ConfirmationToken   string     `gorm:"type:varchar(64);index" json:"-"`
	ConfirmationSentAt  *time.Time `json:"confirmation_sent_at,omitempty"`
	ConfirmedAt         *time.Time `json:"confirmed_at,omitempty"`

	// Unsubscribe
	UnsubscribeToken    string     `gorm:"type:varchar(64);index" json:"-"`
	UnsubscribedAt      *time.Time `json:"unsubscribed_at,omitempty"`
	UnsubscribeReason   string     `gorm:"type:varchar(255)" json:"unsubscribe_reason,omitempty"`

	// Metadata
	IPAddress string `gorm:"type:varchar(45)" json:"-"`
	UserAgent string `gorm:"type:varchar(500)" json:"-"`
	Tags      string `gorm:"type:text" json:"tags,omitempty" example:"courses,news"` // comma-separated

	CreatedAt time.Time      `json:"created_at" example:"2025-01-15T10:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2025-01-15T10:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Subscriber) TableName() string {
	return "subscribers"
}
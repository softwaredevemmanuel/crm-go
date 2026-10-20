package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type CourseLevel string

const (
	LevelBeginner     CourseLevel = "beginner"
	LevelIntermediate CourseLevel = "intermediate"
	LevelAdvanced     CourseLevel = "advanced"
)

type CourseStatus string

const (
	StatusDraft     CourseStatus = "draft"
	StatusPublished CourseStatus = "published"
	StatusArchived  CourseStatus = "archived"
)

type Course struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string    `gorm:"size:255;not null;index" json:"title"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"`
	Slug        string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`

	// Media
	Image         string `gorm:"type:text" json:"image"`          // cover (kept for backwards compat)
	ThumbnailURL  string `gorm:"type:text" json:"thumbnail_url"`  // small grid image
	CoverImageURL string `gorm:"type:text" json:"cover_image_url"`// hero banner
	VideoURL      string `gorm:"type:text" json:"video_url"`      // full lesson/trailer
	IntroVideoURL string `gorm:"type:text" json:"intro_video_url"`// short preview


	// Classification
	Level  CourseLevel  `gorm:"type:varchar(20);not null;default:'beginner';index;check:level IN ('beginner','intermediate','advanced')" json:"level"`
	Status CourseStatus `gorm:"type:varchar(20);not null;default:'draft';index;check:status IN ('draft','published','archived')" json:"status"`

	// Content
	LearningOutcomes datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"learning_outcomes"`
	Requirements     datatypes.JSON `gorm:"type:jsonb;default:'[]'" json:"requirements"`

	// Pricing
	IsFree     bool   `gorm:"default:true" json:"is_free"`
	PriceCents int    `gorm:"default:0" json:"price_cents"`
	Currency   string `gorm:"type:varchar(3);default:'NGN'" json:"currency"`

	// Denormalized metrics
	DurationMinutes int     `gorm:"default:0" json:"duration_minutes"`
	TotalLessons    int     `gorm:"default:0" json:"total_lessons"`
	EnrollmentCount int     `gorm:"default:0" json:"enrollment_count"`
	Rating          float64 `gorm:"type:numeric(3,2);default:0" json:"rating"`
	RatingCount     int     `gorm:"default:0" json:"rating_count"`

	// Access
	RequiresEnrollment bool `gorm:"default:true" json:"requires_enrollment"`
	IsRestricted       bool `gorm:"default:false" json:"is_restricted"`

	// Lifecycle
	PublishedAt *time.Time     `json:"published_at,omitempty"`
	ArchivedAt  *time.Time     `json:"archived_at,omitempty"`
	CreatedBy   *uuid.UUID     `gorm:"type:uuid" json:"created_by,omitempty"`
	UpdatedBy   *uuid.UUID     `gorm:"type:uuid" json:"updated_by,omitempty"`

	// Relationships
	Products     []CourseProductTable  `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Categories   []CourseCategoryTable `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	Materials    []CourseMaterial      `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	// Enrollments  []Enrollment          `gorm:"constraint:OnDelete:CASCADE;" json:"-"`
	// Reviews      []CourseReview        `gorm:"constraint:OnDelete:CASCADE;" json:"-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Course) TableName() string { return "courses" }

type CourseMiniResponse struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Slug         string    `json:"slug"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Level        string    `json:"level"`
}

type CourseListResponse struct {
	ID              uuid.UUID `json:"id"`
	Title           string    `json:"title"`
	Slug            string    `json:"slug"`
	ThumbnailURL    string    `json:"thumbnail_url"`
	Level           string    `json:"level"`
	Status          string    `json:"status"`
	Rating          float64   `json:"rating"`
	EnrollmentCount int       `json:"enrollment_count"`
	DurationMinutes int       `json:"duration_minutes"`
	IsFree          bool      `json:"is_free"`
}

type CourseResponse struct {
	ID               uuid.UUID      `json:"id"`
	Title            string         `json:"title"`
	Code             string         `json:"code"`
	Slug             string         `json:"slug"`
	Description      string         `json:"description"`
	Image            string         `json:"image"`
	ThumbnailURL     string         `json:"thumbnail_url"`
	CoverImageURL    string         `json:"cover_image_url"`
	VideoURL         string         `json:"video_url"`
	IntroVideoURL    string         `json:"intro_video_url"`
	Level            string         `json:"level"`
	Status           string         `json:"status"`
	LearningOutcomes datatypes.JSON `json:"learning_outcomes"`
	Requirements     datatypes.JSON `json:"requirements"`
	IsFree           bool           `json:"is_free"`
	PriceCents       int            `json:"price_cents"`
	Currency         string         `json:"currency"`
	DurationMinutes  int            `json:"duration_minutes"`
	TotalLessons     int            `json:"total_lessons"`
	EnrollmentCount  int            `json:"enrollment_count"`
	Rating           float64        `json:"rating"`
	RatingCount      int            `json:"rating_count"`
	PublishedAt      *time.Time     `json:"published_at,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"
	"os"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
	"crm-go/utils"
)

var (
	ErrCampaignNotFound     = errors.New("campaign not found")
	ErrCampaignAlreadySent  = errors.New("campaign has already been sent")
	ErrCampaignNoRecipients = errors.New("no recipients match this audience")
	ErrCampaignNotDraft     = errors.New("only draft campaigns can be edited")
)

type NewsletterCampaignService struct {
	DB           *gorm.DB
	emailService utils.EmailService // reuse the interface from newsletter_service.go
}

func NewNewsletterCampaignService(db *gorm.DB, emailService utils.EmailService) *NewsletterCampaignService {
	return &NewsletterCampaignService{DB: db, emailService: emailService}
}

// ---------- helpers ----------
func (s *NewsletterCampaignService) buildRecipientQuery(audience models.CampaignAudience, tag string, manualIDs []string) *gorm.DB {
	q := s.DB.Model(&models.Subscriber{})

	switch audience {
	case models.AudienceAll:
		// All non-deleted subscribers
	case models.AudienceActive:
		q = q.Where("status = ?", models.SubscriberSubscribed)
	case models.AudienceTagged:
		q = q.Where("status = ?", models.SubscriberSubscribed).
			Where("tags LIKE ?", "%"+tag+"%")
	case models.AudienceManual:
		q = q.Where("id IN ?", manualIDs)
	default:
		q = q.Where("status = ?", models.SubscriberSubscribed)
	}

	return q
}

// ---------- Create ----------
func (s *NewsletterCampaignService) Create(req dto.CreateCampaignRequest, createdBy *uuid.UUID) (*models.NewsletterCampaign, error) {
	audience := models.CampaignAudience(req.Audience)
	if audience == "" {
		audience = models.AudienceActive
	}

	// Serialize manual IDs
	var manualJSON string
	if len(req.ManualIDs) > 0 {
		b, _ := json.Marshal(req.ManualIDs)
		manualJSON = string(b)
	}

	// Parse scheduled_at if provided
	var scheduledAt *time.Time
	if req.ScheduledAt != nil && *req.ScheduledAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
		if err != nil {
			return nil, fmt.Errorf("invalid scheduled_at (expected RFC3339): %w", err)
		}
		scheduledAt = &t
	}

	status := models.CampaignDraft
	if scheduledAt != nil {
		status = models.CampaignScheduled
	}

	campaign := &models.NewsletterCampaign{
		Subject:     strings.TrimSpace(req.Subject),
		Title:       strings.TrimSpace(req.Title),
		Content:     req.Content,
		Preview:     strings.TrimSpace(req.Preview),
		Audience:    audience,
		AudienceTag: strings.TrimSpace(req.AudienceTag),
		ManualIDs:   manualJSON,
		Status:      status,
		ScheduledAt: scheduledAt,
		CreatedBy:   createdBy,
	}

	if err := s.DB.Create(campaign).Error; err != nil {
		return nil, err
	}
	return campaign, nil
}

// ---------- Get ----------
func (s *NewsletterCampaignService) GetByID(id uuid.UUID) (*models.NewsletterCampaign, error) {
	var c models.NewsletterCampaign
	if err := s.DB.First(&c, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCampaignNotFound
		}
		return nil, err
	}
	return &c, nil
}

// ---------- List ----------
func (s *NewsletterCampaignService) List(q dto.ListCampaignsQuery) ([]models.NewsletterCampaign, int64, error) {
	var (
		campaigns []models.NewsletterCampaign
		total     int64
	)

	query := s.DB.Model(&models.NewsletterCampaign{})

	if q.Search != "" {
		s := "%" + strings.ToLower(q.Search) + "%"
		query = query.Where(
			"LOWER(subject) LIKE ? OR LOWER(title) LIKE ?",
			s, s,
		)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]bool{
		"created_at": true, "updated_at": true, "sent_at": true, "subject": true,
	}
	sortBy := q.SortBy
	if !allowedSort[sortBy] {
		sortBy = "created_at"
	}
	sortOrder := strings.ToLower(q.SortOrder)
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}
	query = query.Order(fmt.Sprintf("%s %s", sortBy, sortOrder))

	offset := (q.Page - 1) * q.Limit
	if err := query.Offset(offset).Limit(q.Limit).Find(&campaigns).Error; err != nil {
		return nil, 0, err
	}

	return campaigns, total, nil
}

// ---------- Update ----------
func (s *NewsletterCampaignService) Update(id uuid.UUID, req dto.UpdateCampaignRequest) (*models.NewsletterCampaign, error) {
	c, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c.Status == models.CampaignSent || c.Status == models.CampaignSending {
		return nil, ErrCampaignAlreadySent
	}

	if req.Subject != nil {
		c.Subject = strings.TrimSpace(*req.Subject)
	}
	if req.Title != nil {
		c.Title = strings.TrimSpace(*req.Title)
	}
	if req.Content != nil {
		c.Content = *req.Content
	}
	if req.Preview != nil {
		c.Preview = strings.TrimSpace(*req.Preview)
	}
	if req.Audience != nil {
		c.Audience = models.CampaignAudience(*req.Audience)
	}
	if req.AudienceTag != nil {
		c.AudienceTag = strings.TrimSpace(*req.AudienceTag)
	}
	if req.ManualIDs != nil {
		b, _ := json.Marshal(*req.ManualIDs)
		c.ManualIDs = string(b)
	}
	if req.ScheduledAt != nil {
		if *req.ScheduledAt == "" {
			c.ScheduledAt = nil
			if c.Status == models.CampaignScheduled {
				c.Status = models.CampaignDraft
			}
		} else {
			t, err := time.Parse(time.RFC3339, *req.ScheduledAt)
			if err != nil {
				return nil, fmt.Errorf("invalid scheduled_at (expected RFC3339): %w", err)
			}
			c.ScheduledAt = &t
			if c.Status == models.CampaignDraft {
				c.Status = models.CampaignScheduled
			}
		}
	}

	if err := s.DB.Save(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

// ---------- Delete ----------
func (s *NewsletterCampaignService) Delete(id uuid.UUID) error {
	c, err := s.GetByID(id)
	if err != nil {
		return err
	}
	if c.Status == models.CampaignSending {
		return errors.New("cannot delete a campaign that is currently sending")
	}
	return s.DB.Delete(c).Error
}

// ---------- Preview recipients ----------
func (s *NewsletterCampaignService) PreviewRecipients(audience, tag string, manualIDs []string) (int64, error) {
	q := s.buildRecipientQuery(models.CampaignAudience(audience), tag, manualIDs)
	var count int64
	if err := q.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ---------- Send ----------
// Send dispatches the campaign to its audience.
// In a real deployment this should push to a queue (Redis / SQS / etc.).
func (s *NewsletterCampaignService) Send(id uuid.UUID) (int, error) {
	c, err := s.GetByID(id)
	if err != nil {
		return 0, err
	}
	if c.Status == models.CampaignSent || c.Status == models.CampaignSending {
		return 0, ErrCampaignAlreadySent
	}

	// Decode manual IDs
	var manualIDs []string
	if c.ManualIDs != "" {
		_ = json.Unmarshal([]byte(c.ManualIDs), &manualIDs)
	}

	// Resolve recipients
	var recipients []models.Subscriber
	q := s.buildRecipientQuery(c.Audience, c.AudienceTag, manualIDs)
	if err := q.Find(&recipients).Error; err != nil {
		return 0, err
	}
	if len(recipients) == 0 {
		return 0, ErrCampaignNoRecipients
	}

	// Mark as sending
	now := time.Now()
	c.Status = models.CampaignSending
	c.SentAt = &now
	c.RecipientCount = len(recipients)
	if err := s.DB.Save(c).Error; err != nil {
		return 0, err
	}

	// Dispatch asynchronously (best-effort)
	go s.dispatch(c, recipients)

	return len(recipients), nil
}

func (s *NewsletterCampaignService) dispatch(c *models.NewsletterCampaign, recipients []models.Subscriber) {
	var delivered, failed int

	for i := range recipients {
		sub := recipients[i]

		// Personalize
		content := strings.ReplaceAll(c.Content, "{{.Name}}", sub.Name)
		content = strings.ReplaceAll(content, "{{.Email}}", sub.Email)

		// Unsubscribe URL
		unsubURL := ""
		if sub.UnsubscribeToken != "" {
			unsubURL = fmt.Sprintf(
				"%s/unsubscribe-newsletter?token=%s",
				frontendBaseURL(),
				sub.UnsubscribeToken,
			)
		}

		err := s.emailService.SendCampaignEmail(
			sub.Email,
			sub.Name,
			c.Subject,
			c.Title,
			content,
			c.Preview,
			unsubURL,
		)
		if err != nil {
			failed++
			continue
		}
		delivered++
	}

	// Update metrics
	c.DeliveredCount = delivered
	c.FailedCount = failed
	if failed == len(recipients) {
		c.Status = models.CampaignFailed
		c.FailureReason = "all deliveries failed"
	} else {
		c.Status = models.CampaignSent
	}
	_ = s.DB.Save(c).Error
}

// ---------- Cancel scheduled ----------
func (s *NewsletterCampaignService) Cancel(id uuid.UUID) (*models.NewsletterCampaign, error) {
	c, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c.Status != models.CampaignScheduled && c.Status != models.CampaignDraft {
		return nil, errors.New("only draft or scheduled campaigns can be cancelled")
	}
	c.Status = models.CampaignCancelled
	c.ScheduledAt = nil
	if err := s.DB.Save(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}


func frontendBaseURL() string {
	u := os.Getenv("FRONTEND_URL")
	if u == "" { u = "http://localhost:3000" }
	return u
}
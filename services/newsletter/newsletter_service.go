package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"log"
	"crm-go/dto"
	"crm-go/utils"
	"crm-go/models"
)

var (
	ErrSubscriberNotFound     = errors.New("subscriber not found")
	ErrSubscriberAlreadyThere = errors.New("this email is already subscribed")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrEmailAlreadyPending    = errors.New("this email has a pending confirmation")
)



type NewsletterService struct {
	DB           *gorm.DB
	emailService utils.EmailService
}

func NewNewsletterService(
	db *gorm.DB,
	emailService utils.EmailService,
) *NewsletterService {
	return &NewsletterService{
		DB:           db,
		emailService: emailService,
	}
}

// ---------- helpers ----------
func randomToken(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// ---------- Subscribe ----------
func (s *NewsletterService) Subscribe(req dto.SubscribeRequest, ip, userAgent string) (*models.Subscriber, error) {
	email := normalizeEmail(req.Email)

	var existing models.Subscriber
	err := s.DB.Where("email = ?", email).First(&existing).Error

	// Already subscribed → idempotent success
	if err == nil {
		if existing.Status == models.SubscriberSubscribed {
			return &existing, nil
		}
		if existing.Status == models.SubscriberPending {
			// Regenerate confirmation token & resend
			existing.ConfirmationToken = randomToken(32)
			now := time.Now()
			existing.ConfirmationSentAt = &now
			existing.Source = models.SubscriberSource(req.Source)
			if req.Name != "" {
				existing.Name = req.Name
			}
			if err := s.DB.Save(&existing).Error; err != nil {
				return nil, err
			}

			s.sendConfirmationAsync(existing.Email, existing.Name, existing.ConfirmationToken)
			return &existing, nil
		}

		// Unsubscribed / bounced / complained → reactivate as pending
		existing.Status = models.SubscriberPending
		existing.ConfirmationToken = randomToken(32)
		now := time.Now()
		existing.ConfirmationSentAt = &now
		existing.UnsubscribedAt = nil
		existing.UnsubscribeReason = ""
		if req.Name != "" {
			existing.Name = req.Name
		}
		if err := s.DB.Save(&existing).Error; err != nil {
			return nil, err
		}

		s.sendConfirmationAsync(existing.Email, existing.Name, existing.ConfirmationToken)
		return &existing, nil
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	source := models.SubscriberSource(req.Source)
	if source == "" {
		source = models.SourceOther
	}

	sub := &models.Subscriber{
		Email:             email,
		Name:              strings.TrimSpace(req.Name),
		Status:            models.SubscriberPending,
		Source:            source,
		Tags:              strings.TrimSpace(req.Tags),
		ConfirmationToken: randomToken(32),
		IPAddress:         ip,
		UserAgent:         userAgent,
	}
	now := time.Now()
	sub.ConfirmationSentAt = &now

	if err := s.DB.Create(sub).Error; err != nil {
		return nil, err
	}

	s.sendConfirmationAsync(sub.Email, sub.Name, sub.ConfirmationToken)

	return sub, nil
}

// sendConfirmationAsync fires the confirmation email in a background goroutine
// so slow SMTP calls never block the HTTP response.
func (s *NewsletterService) sendConfirmationAsync(to, name, token string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("[newsletter] panic while sending confirmation email: %v", r)
			}
		}()

		if err := s.emailService.SendConfirmationEmail(to, name, token); err != nil {
			log.Printf("[newsletter] failed to send confirmation email to %s: %v", to, err)
		}
	}()
}

// ---------- Confirm ----------
func (s *NewsletterService) ConfirmByToken(token string) (*models.Subscriber, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	var sub models.Subscriber
	if err := s.DB.Where("confirmation_token = ?", token).First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if sub.Status == models.SubscriberSubscribed {
		return &sub, nil // idempotent
	}

	now := time.Now()
	sub.Status = models.SubscriberSubscribed
	sub.ConfirmedAt = &now
	sub.ConfirmationToken = "" // burn the token
	sub.UnsubscribeToken = randomToken(32)

	if err := s.DB.Save(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

// ---------- Unsubscribe ----------
func (s *NewsletterService) UnsubscribeByToken(token, reason string) (*models.Subscriber, error) {
	if token == "" {
		return nil, ErrInvalidToken
	}

	var sub models.Subscriber
	if err := s.DB.Where("unsubscribe_token = ?", token).First(&sub).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidToken
		}
		return nil, err
	}

	if sub.Status == models.SubscriberUnsubscribed {
		return &sub, nil // idempotent
	}

	now := time.Now()
	sub.Status = models.SubscriberUnsubscribed
	sub.UnsubscribedAt = &now
	sub.UnsubscribeReason = strings.TrimSpace(reason)

	if err := s.DB.Save(&sub).Error; err != nil {
		return nil, err
	}
	return &sub, nil
}

// ---------- Admin: List ----------
func (s *NewsletterService) List(q dto.ListSubscribersQuery) ([]models.Subscriber, int64, error) {
	var (
		subs  []models.Subscriber
		total int64
	)

	query := s.DB.Model(&models.Subscriber{})

	if q.Search != "" {
		s := "%" + strings.ToLower(q.Search) + "%"
		query = query.Where("LOWER(email) LIKE ? OR LOWER(name) LIKE ?", s, s)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if q.Source != "" {
		query = query.Where("source = ?", q.Source)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	allowedSort := map[string]bool{
		"created_at": true, "updated_at": true, "email": true, "status": true,
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
	if err := query.Offset(offset).Limit(q.Limit).Find(&subs).Error; err != nil {
		return nil, 0, err
	}

	return subs, total, nil
}

// ---------- Admin: Get ----------
func (s *NewsletterService) GetByID(id uuid.UUID) (*models.Subscriber, error) {
	var sub models.Subscriber
	if err := s.DB.First(&sub, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrSubscriberNotFound
		}
		return nil, err
	}
	return &sub, nil
}

// ---------- Admin: Create (bypasses opt-in) ----------
func (s *NewsletterService) AdminCreate(req dto.AdminCreateSubscriberRequest) (*models.Subscriber, error) {
	email := normalizeEmail(req.Email)

	var existing models.Subscriber
	if err := s.DB.Where("email = ?", email).First(&existing).Error; err == nil {
		return nil, ErrSubscriberAlreadyThere
	}

	status := models.SubscriberSubscribed
	if req.Status != "" {
		status = models.SubscriberStatus(req.Status)
	}

	sub := &models.Subscriber{
		Email:  email,
		Name:   strings.TrimSpace(req.Name),
		Status: status,
		Source: models.SourceAdmin,
		Tags:   strings.TrimSpace(req.Tags),
	}
	if status == models.SubscriberSubscribed {
		now := time.Now()
		sub.ConfirmedAt = &now
		sub.UnsubscribeToken = randomToken(32)
	} else if status == models.SubscriberPending {
		sub.ConfirmationToken = randomToken(32)
		now := time.Now()
		sub.ConfirmationSentAt = &now
	}

	if err := s.DB.Create(sub).Error; err != nil {
		return nil, err
	}
	return sub, nil
}

// ---------- Admin: Update ----------
func (s *NewsletterService) AdminUpdate(id uuid.UUID, req dto.AdminUpdateSubscriberRequest) (*models.Subscriber, error) {
	sub, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		sub.Name = strings.TrimSpace(*req.Name)
	}
	if req.Tags != nil {
		sub.Tags = strings.TrimSpace(*req.Tags)
	}
	if req.Status != nil {
		newStatus := models.SubscriberStatus(*req.Status)
		now := time.Now()
		// Handle transition side effects
		switch newStatus {
		case models.SubscriberSubscribed:
			if sub.ConfirmedAt == nil {
				sub.ConfirmedAt = &now
			}
			if sub.UnsubscribeToken == "" {
				sub.UnsubscribeToken = randomToken(32)
			}
		case models.SubscriberUnsubscribed:
			if sub.UnsubscribedAt == nil {
				sub.UnsubscribedAt = &now
			}
		}
		sub.Status = newStatus
	}

	if err := s.DB.Save(sub).Error; err != nil {
		return nil, err
	}
	return sub, nil
}

// ---------- Admin: Delete ----------
func (s *NewsletterService) Delete(id uuid.UUID) error {
	sub, err := s.GetByID(id)
	if err != nil {
		return err
	}
	return s.DB.Delete(sub).Error
}

// ---------- Admin: Stats ----------
func (s *NewsletterService) Stats() (*dto.NewsletterStats, error) {
	stats := &dto.NewsletterStats{}

	if err := s.DB.Model(&models.Subscriber{}).Count(&stats.Total).Error; err != nil {
		return nil, err
	}
	_ = s.DB.Model(&models.Subscriber{}).
		Where("status = ?", models.SubscriberSubscribed).
		Count(&stats.Subscribed).Error
	_ = s.DB.Model(&models.Subscriber{}).
		Where("status = ?", models.SubscriberPending).
		Count(&stats.Pending).Error
	_ = s.DB.Model(&models.Subscriber{}).
		Where("status = ?", models.SubscriberUnsubscribed).
		Count(&stats.Unsubscribed).Error
	_ = s.DB.Model(&models.Subscriber{}).
		Where("status = ?", models.SubscriberBounced).
		Count(&stats.Bounced).Error

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)
	_ = s.DB.Model(&models.Subscriber{}).
		Where("created_at >= ?", thirtyDaysAgo).
		Count(&stats.Last30Days).Error

	return stats, nil
}
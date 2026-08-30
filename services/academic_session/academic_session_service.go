// services/academic_session_service.go
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

type AcademicSessionService struct {
	db *gorm.DB
}

func NewAcademicSessionService(db *gorm.DB) *AcademicSessionService {
	return &AcademicSessionService{db: db}
}

// CreateAcademicSession creates a new academic session
func (s *AcademicSessionService) CreateAcademicSession(req *dto.CreateAcademicSessionRequest, userID uuid.UUID) (*dto.AcademicSessionResponse, error) {
	// Validate input
	if err := s.validateSessionRequest(req); err != nil {
		return nil, err
	}

	// Check if code already exists
	var existing models.AcademicSession
	if err := s.db.Where("code = ? AND deleted_at IS NULL", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("academic session code already exists")
	}

	// Parse dates
	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("invalid start date format. Use YYYY-MM-DD")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("invalid end date format. Use YYYY-MM-DD")
	}

	// Validate date range
	if endDate.Before(startDate) {
		return nil, errors.New("end date must be after start date")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// If IsCurrent is true, set all other sessions to false
	if req.IsCurrent {
		if err := s.db.Model(&models.AcademicSession{}).
			Where("deleted_at IS NULL").
			Update("is_current", false).Error; err != nil {
			return nil, errors.New("failed to update current sessions: " + err.Error())
		}
	}

	// Create academic session
	session := &models.AcademicSession{
		ID:           uuid.New(),
		AcademicYear: req.AcademicYear,
		Code:         req.Code,
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		IsCurrent:    req.IsCurrent,
		Description:  req.Description,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, errors.New("failed to create academic session: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Creator").First(session, session.ID).Error; err != nil {
		return nil, errors.New("failed to load session details: " + err.Error())
	}

	return s.toSessionResponse(session), nil
}

// GetAllSessions retrieves all academic sessions with pagination and filters
func (s *AcademicSessionService) GetAllSessions(params *dto.AcademicSessionQueryParams) (*dto.AcademicSessionListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.AcademicSession{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsCurrent != nil {
		query = query.Where("is_current = ?", *params.IsCurrent)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(academic_year) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			search, search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
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

	// Execute with preloads
	var sessions []models.AcademicSession
	if err := query.Preload("Creator").Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch sessions: %w", err)
	}

	// Convert to response
	responses := make([]dto.AcademicSessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = *s.toSessionResponse(&session)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.AcademicSessionListResponse{
		Sessions:   responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetSessionByID retrieves a single academic session by ID
func (s *AcademicSessionService) GetSessionByID(id string) (*dto.AcademicSessionResponse, error) {
	sessionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid session ID")
	}

	var session models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).
		Preload("Creator").
		Preload("Terms").
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic session not found")
		}
		return nil, errors.New("failed to fetch session: " + err.Error())
	}

	return s.toSessionResponse(&session), nil
}

// GetCurrentSession retrieves the current academic session
func (s *AcademicSessionService) GetCurrentSession() (*dto.AcademicSessionResponse, error) {
	var session models.AcademicSession
	if err := s.db.Where("is_current = ? AND deleted_at IS NULL", true).
		Preload("Creator").
		Preload("Terms").
		First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no current academic session found")
		}
		return nil, errors.New("failed to fetch current session: " + err.Error())
	}

	return s.toSessionResponse(&session), nil
}

// GetActiveSessions retrieves all active academic sessions
func (s *AcademicSessionService) GetActiveSessions() ([]dto.AcademicSessionResponse, error) {
	var sessions []models.AcademicSession
	if err := s.db.Where("status = ? AND deleted_at IS NULL", "active").
		Preload("Creator").
		Order("created_at DESC").
		Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch active sessions: %w", err)
	}

	responses := make([]dto.AcademicSessionResponse, len(sessions))
	for i, session := range sessions {
		responses[i] = *s.toSessionResponse(&session)
	}

	return responses, nil
}

// GetSessionStats retrieves statistics for academic sessions
func (s *AcademicSessionService) GetSessionStats() (*dto.AcademicSessionStats, error) {
	var stats dto.AcademicSessionStats

	// Count total sessions
	if err := s.db.Model(&models.AcademicSession{}).Where("deleted_at IS NULL").Count(&stats.TotalSessions).Error; err != nil {
		return nil, fmt.Errorf("failed to count total sessions: %w", err)
	}

	// Count by status
	if err := s.db.Model(&models.AcademicSession{}).Where("status = ? AND deleted_at IS NULL", "active").Count(&stats.ActiveSessions).Error; err != nil {
		return nil, fmt.Errorf("failed to count active sessions: %w", err)
	}
	if err := s.db.Model(&models.AcademicSession{}).Where("status = ? AND deleted_at IS NULL", "inactive").Count(&stats.InactiveSessions).Error; err != nil {
		return nil, fmt.Errorf("failed to count inactive sessions: %w", err)
	}
	if err := s.db.Model(&models.AcademicSession{}).Where("status = ? AND deleted_at IS NULL", "completed").Count(&stats.CompletedSessions).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed sessions: %w", err)
	}

	// Count current sessions
	if err := s.db.Model(&models.AcademicSession{}).Where("is_current = ? AND deleted_at IS NULL", true).Count(&stats.CurrentSessions).Error; err != nil {
		return nil, fmt.Errorf("failed to count current sessions: %w", err)
	}

	// Count total terms
	if err := s.db.Model(&models.Term{}).Where("deleted_at IS NULL").Count(&stats.TotalTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to count total terms: %w", err)
	}

	return &stats, nil
}

// UpdateAcademicSession updates an existing academic session
func (s *AcademicSessionService) UpdateAcademicSession(id string, req *dto.UpdateAcademicSessionRequest) (*dto.AcademicSessionResponse, error) {
	sessionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid session ID")
	}

	// Find existing session
	var session models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic session not found")
		}
		return nil, errors.New("failed to fetch session: " + err.Error())
	}

	// Update fields
	if req.AcademicYear != "" {
		session.AcademicYear = req.AcademicYear
	}

	if req.Code != "" {
		// Check if code already exists for another session
		var existing models.AcademicSession
		if err := s.db.Where("code = ? AND id != ? AND deleted_at IS NULL", req.Code, sessionID).First(&existing).Error; err == nil {
			return nil, errors.New("academic session code already exists")
		}
		session.Code = req.Code
	}

	if req.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start date format. Use YYYY-MM-DD")
		}
		session.StartDate = startDate
	}

	if req.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format. Use YYYY-MM-DD")
		}
		session.EndDate = endDate
	}

	// Validate date range
	if !session.EndDate.IsZero() && !session.StartDate.IsZero() {
		if session.EndDate.Before(session.StartDate) {
			return nil, errors.New("end date must be after start date")
		}
	}

	if req.Status != "" {
		if req.Status != "active" && req.Status != "inactive" && req.Status != "completed" {
			return nil, errors.New("status must be 'active', 'inactive', or 'completed'")
		}
		session.Status = req.Status
	}

	if req.IsCurrent != nil {
		// If setting to true, set all other sessions to false
		if *req.IsCurrent {
			if err := s.db.Model(&models.AcademicSession{}).
				Where("id != ? AND deleted_at IS NULL", sessionID).
				Update("is_current", false).Error; err != nil {
				return nil, errors.New("failed to update current sessions: " + err.Error())
			}
		}
		session.IsCurrent = *req.IsCurrent
	}

	if req.Description != "" {
		session.Description = req.Description
	}

	session.UpdatedAt = time.Now()

	if err := s.db.Save(&session).Error; err != nil {
		return nil, errors.New("failed to update academic session: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Creator").First(&session, session.ID).Error; err != nil {
		return nil, errors.New("failed to load session details: " + err.Error())
	}

	return s.toSessionResponse(&session), nil
}

// DeleteAcademicSession soft deletes an academic session
func (s *AcademicSessionService) DeleteAcademicSession(id string) error {
	sessionID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid session ID")
	}

	var session models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("academic session not found")
		}
		return errors.New("failed to fetch session: " + err.Error())
	}

	// Check if session has related terms
	var termCount int64
	if err := s.db.Model(&models.Term{}).Where("academic_session_id = ? AND deleted_at IS NULL", sessionID).Count(&termCount).Error; err != nil {
		return errors.New("failed to check related terms: " + err.Error())
	}

	if termCount > 0 {
		return errors.New("cannot delete academic session with existing terms")
	}

	if err := s.db.Delete(&session).Error; err != nil {
		return errors.New("failed to delete academic session: " + err.Error())
	}

	return nil
}

// validateSessionRequest validates the session request
func (s *AcademicSessionService) validateSessionRequest(req *dto.CreateAcademicSessionRequest) error {
	if req.AcademicYear == "" {
		return errors.New("academic year is required")
	}
	if req.Code == "" {
		return errors.New("code is required")
	}
	if req.StartDate == "" {
		return errors.New("start date is required")
	}
	if req.EndDate == "" {
		return errors.New("end date is required")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "completed" {
		return errors.New("status must be 'active', 'inactive', or 'completed'")
	}
	return nil
}

// toSessionResponse converts model to response DTO
func (s *AcademicSessionService) toSessionResponse(session *models.AcademicSession) *dto.AcademicSessionResponse {
	response := &dto.AcademicSessionResponse{
		ID:           session.ID.String(),
		AcademicYear: session.AcademicYear,
		Code:         session.Code,
		StartDate:    session.StartDate,
		EndDate:      session.EndDate,
		Status:       session.Status,
		IsCurrent:    session.IsCurrent,
		Description:  session.Description,
		CreatedBy:    session.CreatedBy.String(),
		CreatedAt:    session.CreatedAt,
		UpdatedAt:    session.UpdatedAt,
	}

	// Add creator details if preloaded
	if session.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        session.Creator.ID.String(),
			FirstName: session.Creator.FirstName,
			LastName:  session.Creator.LastName,
			Email:     session.Creator.Email,
			Phone:     session.Creator.Phone,
			Role:      session.Creator.Role,
		}
	}

	// Add terms if preloaded
	if len(session.Terms) > 0 {
		terms := make([]dto.TermResponse, len(session.Terms))
		for i, term := range session.Terms {
			terms[i] = dto.TermResponse{
				ID:         term.ID.String(),
				Name:       term.Name,
				Code:       term.Code,
				TermNumber: term.TermNumber,
				StartDate:  term.StartDate,
				EndDate:    term.EndDate,
				IsCurrent:  term.IsCurrent,
				Status:     term.Status,
			}
		}
		response.Terms = terms
	}

	return response
}
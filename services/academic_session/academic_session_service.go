// services/academic_session_service.go
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"crm-go/dto"
	"crm-go/models"
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

	// Check if session with same academic year and term already exists
	var existing models.AcademicSession

	if err := s.db.Where(
		"academic_year = ? AND term = ?",
		req.AcademicYear,
		req.Term,
	).First(&existing).Error; err == nil {

		return nil, errors.New(
			"academic session with this academic year and term already exists",
		)
	}

	if err := s.db.Where(
		"is_current = ?",
		req.IsCurrent,
	).First(&existing).Error; err == nil {

		return nil, errors.New(
			"an academic session is currently running",
		)
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}

	// If this session is set as current, unset any existing current sessions
	if req.IsCurrent {
		if err := s.db.Model(&models.AcademicSession{}).
			Where("is_current = ?", true).
			Update("is_current", false).Error; err != nil {
			return nil, errors.New("failed to update current sessions: " + err.Error())
		}
	}

	// Create new academic session
	session := &models.AcademicSession{
		ID:           uuid.New(),
		AcademicYear: strings.TrimSpace(req.AcademicYear),
		Code:         strings.ToUpper(strings.TrimSpace(req.Code)),
		Term:         strings.TrimSpace(req.Term),
		StartDate:    startDate,
		EndDate:      endDate,
		Status:       status,
		IsCurrent:    req.IsCurrent,
		Description:  strings.TrimSpace(req.Description),
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(session).Error; err != nil {
		return nil, errors.New("failed to create academic session: " + err.Error())
	}

	return s.toSessionResponse(session), nil
}

// GetAllAcademicSessions retrieves all academic sessions with pagination and filters
func (s *AcademicSessionService) GetAllAcademicSessions(params *dto.AcademicSessionQueryParams) (*dto.AcademicSessionListResponse, error) {
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
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(year) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsCurrent != nil {
		query = query.Where("is_current = ?", *params.IsCurrent)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count academic sessions: %w", err)
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
	var sessions []models.AcademicSession
	if err := query.Find(&sessions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch academic sessions: %w", err)
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

// GetAcademicSessionByID retrieves a single academic session by ID
func (s *AcademicSessionService) GetAcademicSessionByID(id string) (*dto.AcademicSessionResponse, error) {
	sessionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid academic session ID")
	}

	var session models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic session not found")
		}
		return nil, errors.New("failed to fetch academic session: " + err.Error())
	}

	return s.toSessionResponse(&session), nil
}

// GetCurrentAcademicSession retrieves the current academic session
func (s *AcademicSessionService) GetCurrentAcademicSession() (*dto.AcademicSessionResponse, error) {
	var session models.AcademicSession
	if err := s.db.Where("is_current = ? AND deleted_at IS NULL", true).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no current academic session found")
		}
		return nil, errors.New("failed to fetch current academic session: " + err.Error())
	}

	return s.toSessionResponse(&session), nil
}

// UpdateAcademicSession updates an existing academic session
func (s *AcademicSessionService) UpdateAcademicSession(
	id string,
	req *dto.UpdateAcademicSessionRequest,
	userID uuid.UUID,
) (*dto.AcademicSessionResponse, error) {

	sessionID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid academic session ID")
	}

	// Find existing session
	var session models.AcademicSession

	if err := s.db.
		Where("id = ? AND deleted_at IS NULL", sessionID).
		First(&session).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic session not found")
		}

		return nil, errors.New(
			"failed to fetch academic session: " + err.Error(),
		)
	}

	// ---------------------------------------------------------
	// Determine the final Academic Year and Term
	// ---------------------------------------------------------

	academicYear := session.AcademicYear
	term := session.Term

	if req.AcademicYear != "" {
		academicYear = strings.TrimSpace(req.AcademicYear)
	}

	if req.Term != "" {
		term = strings.TrimSpace(req.Term)
	}

	// ---------------------------------------------------------
	// Check duplicate Academic Year + Term combination
	// ---------------------------------------------------------

	var existing models.AcademicSession

	result := s.db.
		Where(
			"academic_year = ? AND term = ? AND id != ? AND deleted_at IS NULL",
			academicYear,
			term,
			sessionID,
		).
		First(&existing)

	if result.Error == nil {
		return nil, errors.New(
			"academic session with this academic year and term already exists",
		)
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, errors.New(
			"failed to check academic session: " + result.Error.Error(),
		)
	}


	if err := s.db.Where(
		"is_current = ?",
		req.IsCurrent,
	).First(&existing).Error; err == nil {

		return nil, errors.New(
			"an active academic session is currently running",
		)
	}
	// ---------------------------------------------------------
	// Check duplicate Code
	// ---------------------------------------------------------

	if req.Code != "" && strings.TrimSpace(req.Code) != session.Code {

		code := strings.ToUpper(strings.TrimSpace(req.Code))

		var existingCode models.AcademicSession

		result := s.db.
			Where(
				"code = ? AND id != ? AND deleted_at IS NULL",
				code,
				sessionID,
			).
			First(&existingCode)

		if result.Error == nil {
			return nil, errors.New(
				"academic session with this code already exists",
			)
		}

		if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New(
				"failed to check academic session code: " + result.Error.Error(),
			)
		}
	}

	// ---------------------------------------------------------
	// Parse dates if provided
	// ---------------------------------------------------------

	var startDate, endDate time.Time

	if req.StartDate != "" {
		startDate, err = time.Parse("2006-01-02", req.StartDate)

		if err != nil {
			return nil, errors.New(
				"invalid start date format. Use YYYY-MM-DD",
			)
		}
	}

	if req.EndDate != "" {
		endDate, err = time.Parse("2006-01-02", req.EndDate)

		if err != nil {
			return nil, errors.New(
				"invalid end date format. Use YYYY-MM-DD",
			)
		}
	}

	// ---------------------------------------------------------
	// Validate date range
	// ---------------------------------------------------------

	finalStartDate := session.StartDate
	finalEndDate := session.EndDate

	if req.StartDate != "" {
		finalStartDate = startDate
	}

	if req.EndDate != "" {
		finalEndDate = endDate
	}

	if finalEndDate.Before(finalStartDate) {
		return nil, errors.New(
			"end date must be after start date",
		)
	}

	// ---------------------------------------------------------
	// If this session becomes current,
	// unset other current sessions
	// ---------------------------------------------------------

	if req.IsCurrent != nil && *req.IsCurrent && !session.IsCurrent {

		if err := s.db.
			Model(&models.AcademicSession{}).
			Where(
				"is_current = ? AND id != ? AND deleted_at IS NULL",
				true,
				sessionID,
			).
			Update("is_current", false).Error; err != nil {

			return nil, errors.New(
				"failed to update current sessions: " + err.Error(),
			)
		}
	}

	// ---------------------------------------------------------
	// Update fields
	// ---------------------------------------------------------

	if req.AcademicYear != "" {
		session.AcademicYear = academicYear
	}

	if req.Term != "" {
		session.Term = term
	}

	if req.Code != "" {
		session.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}

	if req.StartDate != "" {
		session.StartDate = startDate
	}

	if req.EndDate != "" {
		session.EndDate = endDate
	}

	if req.Status != "" {
		session.Status = req.Status
	}

	if req.IsCurrent != nil {
		session.IsCurrent = *req.IsCurrent
	}

	if req.Description != "" {
		session.Description = strings.TrimSpace(req.Description)
	}

	session.UpdatedAt = time.Now()

	// ---------------------------------------------------------
	// Save
	// ---------------------------------------------------------

	if err := s.db.Save(&session).Error; err != nil {
		return nil, errors.New(
			"failed to update academic session: " + err.Error(),
		)
	}

	return s.toSessionResponse(&session), nil
}

// DeleteAcademicSession soft deletes an academic session
func (s *AcademicSessionService) DeleteAcademicSession(id string) error {
	sessionID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid academic session ID")
	}

	var session models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).First(&session).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("academic session not found")
		}
		return errors.New("failed to fetch academic session: " + err.Error())
	}

	if session.IsCurrent {
		return errors.New("cannot delete the current academic session")
	}

	if err := s.db.Delete(&session).Error; err != nil {
		return errors.New("failed to delete academic session: " + err.Error())
	}

	return nil
}

// validateSessionRequest validates the academic session request
func (s *AcademicSessionService) validateSessionRequest(req *dto.CreateAcademicSessionRequest) error {
	if req.AcademicYear == "" {
		return errors.New("session year is required")
	}
	if !strings.Contains(req.AcademicYear, "/") {
		return errors.New("academic year must be in format YYYY/YYYY (e.g., 2024/2025)")
	}
	if req.Code == "" {
		return errors.New("session code is required")
	}
	if len(req.Code) < 2 {
		return errors.New("session code must be at least 2 characters")
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
	now := time.Now()
	daysRemaining := 0
	isActive := false

	if session.Status == "active" {
		if now.After(session.StartDate) && now.Before(session.EndDate) {
			isActive = true
			daysRemaining = int(session.EndDate.Sub(now).Hours() / 24)
		}
	}

	return &dto.AcademicSessionResponse{
		ID:            session.ID.String(),
		AcademicYear:  session.AcademicYear,
		Code:          session.Code,
		Term:          session.Term,
		StartDate:     session.StartDate,
		EndDate:       session.EndDate,
		Status:        session.Status,
		IsCurrent:     session.IsCurrent,
		Description:   session.Description,
		CreatedBy:     session.CreatedBy.String(),
		CreatedAt:     session.CreatedAt,
		UpdatedAt:     session.UpdatedAt,
		DaysRemaining: daysRemaining,
		IsActive:      isActive,
	}
}

// services/term_service.go
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

type TermService struct {
	db *gorm.DB
}

func NewTermService(db *gorm.DB) *TermService {
	return &TermService{db: db}
}

// CreateTerm creates a new term
func (s *TermService) CreateTerm(req *dto.CreateTermRequest, userID uuid.UUID) (*dto.TermResponse, error) {
	// Validate input
	if err := s.validateTermRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	academicSessionID, err := uuid.Parse(req.AcademicSessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID format")
	}

	// Check if academic session exists
	var academicSession models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", academicSessionID).First(&academicSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("academic session not found")
		}
		return nil, errors.New("failed to verify academic session: " + err.Error())
	}

	// Check if term number already exists for this academic session
	var existing models.Term
	if err := s.db.Where("academic_session_id = ? AND term_number = ? AND deleted_at IS NULL",
		academicSessionID, req.TermNumber).First(&existing).Error; err == nil {
		return nil, errors.New("term number already exists for this academic session")
	}

	// Parse dates
	var startDate, endDate *time.Time
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start date format. Use YYYY-MM-DD")
		}
		startDate = &start
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format. Use YYYY-MM-DD")
		}
		endDate = &end
	}

	// If IsCurrent is true, set all other terms to false for this academic session
	if req.IsCurrent {
		if err := s.db.Model(&models.Term{}).
			Where("academic_session_id = ? AND deleted_at IS NULL", academicSessionID).
			Update("is_current", false).Error; err != nil {
			return nil, errors.New("failed to update current terms: " + err.Error())
		}
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create term
	term := &models.Term{
		ID:                uuid.New(),
		AcademicSessionID: academicSessionID,
		Name:              req.Name,
		Code:              req.Code,
		TermNumber:        req.TermNumber,
		StartDate:         startDate,
		EndDate:           endDate,
		IsCurrent:         req.IsCurrent,
		Status:            status,
		Description:       req.Description,
		CreatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(term).Error; err != nil {
		return nil, errors.New("failed to create term: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("AcademicSession").Preload("Creator").First(term, term.ID).Error; err != nil {
		return nil, errors.New("failed to load term details: " + err.Error())
	}

	return s.toTermResponse(term), nil
}



// GetAllTerms retrieves all terms with pagination and filters
func (s *TermService) GetAllTerms(params *dto.TermQueryParams) (*dto.TermListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "term_number"
	}
	if params.SortOrder == "" {
		params.SortOrder = "asc"
	}

	// Build query
	query := s.db.Model(&models.Term{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.AcademicSessionID != "" {
		sessionID, err := uuid.Parse(params.AcademicSessionID)
		if err == nil {
			query = query.Where("academic_session_id = ?", sessionID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsCurrent != nil {
		query = query.Where("is_current = ?", *params.IsCurrent)
	}

	if params.TermNumber > 0 {
		query = query.Where("term_number = ?", params.TermNumber)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(code) LIKE ? OR LOWER(description) LIKE ?", 
			search, search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count terms: %w", err)
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
	var terms []models.Term
	if err := query.Preload("AcademicSession").Preload("Creator").Find(&terms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch terms: %w", err)
	}

	// Convert to response
	responses := make([]dto.TermResponse, len(terms))
	for i, term := range terms {
		responses[i] = *s.toTermResponse(&term)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.TermListResponse{
		Terms:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetTermByID retrieves a single term by ID
func (s *TermService) GetTermByID(id string) (*dto.TermResponse, error) {
	termID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid term ID")
	}

	var term models.Term
	if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).
		Preload("AcademicSession").
		Preload("Creator").
		First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("term not found")
		}
		return nil, errors.New("failed to fetch term: " + err.Error())
	}

	return s.toTermResponse(&term), nil
}

// GetTermsByAcademicSession retrieves all terms for a specific academic session
func (s *TermService) GetTermsByAcademicSession(sessionID string) ([]dto.TermResponse, error) {
	sID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID")
	}

	var terms []models.Term
	if err := s.db.Where("academic_session_id = ? AND deleted_at IS NULL", sID).
		Preload("AcademicSession").
		Preload("Creator").
		Order("term_number ASC").
		Find(&terms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch terms: %w", err)
	}

	responses := make([]dto.TermResponse, len(terms))
	for i, term := range terms {
		responses[i] = *s.toTermResponse(&term)
	}

	return responses, nil
}

// GetCurrentTerm retrieves the current term for an academic session
func (s *TermService) GetCurrentTerm(sessionID string) (*dto.TermResponse, error) {
	sID, err := uuid.Parse(sessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID")
	}

	var term models.Term
	if err := s.db.Where("academic_session_id = ? AND is_current = ? AND deleted_at IS NULL", sID, true).
		Preload("AcademicSession").
		Preload("Creator").
		First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no current term found for this academic session")
		}
		return nil, errors.New("failed to fetch current term: " + err.Error())
	}

	return s.toTermResponse(&term), nil
}

// GetTermStats retrieves statistics for terms
func (s *TermService) GetTermStats(filter map[string]interface{}) (*dto.TermStats, error) {
	query := s.db.Model(&models.Term{}).Where("deleted_at IS NULL")

	// Apply filters
	if sessionID, ok := filter["academic_session_id"].(string); ok && sessionID != "" {
		if id, err := uuid.Parse(sessionID); err == nil {
			query = query.Where("academic_session_id = ?", id)
		}
	}

	var stats dto.TermStats

	// Count total
	if err := query.Count(&stats.TotalTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to count total terms: %w", err)
	}

	// Count by status
	if err := query.Where("status = ?", "active").Count(&stats.ActiveTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to count active terms: %w", err)
	}
	if err := query.Where("status = ?", "inactive").Count(&stats.InactiveTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to count inactive terms: %w", err)
	}
	if err := query.Where("status = ?", "completed").Count(&stats.CompletedTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to count completed terms: %w", err)
	}
	if err := query.Where("is_current = ?", true).Count(&stats.CurrentTerms).Error; err != nil {
		return nil, fmt.Errorf("failed to count current terms: %w", err)
	}

	// Get average terms per session
	var sessionsCount int64
	if err := s.db.Model(&models.AcademicSession{}).Where("deleted_at IS NULL").Count(&sessionsCount).Error; err != nil {
		return nil, fmt.Errorf("failed to count sessions: %w", err)
	}
	if sessionsCount > 0 {
		stats.TermsPerSession = stats.TotalTerms / sessionsCount
	}

	return &stats, nil
}

// UpdateTerm updates an existing term
func (s *TermService) UpdateTerm(id string, req *dto.UpdateTermRequest) (*dto.TermResponse, error) {
	termID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid term ID")
	}

	// Find existing term
	var term models.Term
	if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("term not found")
		}
		return nil, errors.New("failed to fetch term: " + err.Error())
	}

	// Update fields
	if req.Name != "" {
		term.Name = req.Name
	}

	if req.Code != "" {
		term.Code = req.Code
	}

	if req.TermNumber > 0 {
		// Check if another term with this number exists for the same academic session
		var existing models.Term
		if err := s.db.Where("academic_session_id = ? AND term_number = ? AND id != ? AND deleted_at IS NULL",
			term.AcademicSessionID, req.TermNumber, termID).First(&existing).Error; err == nil {
			return nil, errors.New("term number already exists for this academic session")
		}
		term.TermNumber = req.TermNumber
	}

	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return nil, errors.New("invalid start date format. Use YYYY-MM-DD")
		}
		term.StartDate = &start
	}

	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return nil, errors.New("invalid end date format. Use YYYY-MM-DD")
		}
		term.EndDate = &end
	}

	if req.IsCurrent != nil {
		// If setting to true, set all other terms in the same academic session to false
		if *req.IsCurrent {
			if err := s.db.Model(&models.Term{}).
				Where("academic_session_id = ? AND id != ? AND deleted_at IS NULL", 
					term.AcademicSessionID, termID).
				Update("is_current", false).Error; err != nil {
				return nil, errors.New("failed to update current terms: " + err.Error())
			}
		}
		term.IsCurrent = *req.IsCurrent
	}

	if req.Status != "" {
		if req.Status != "active" && req.Status != "inactive" && req.Status != "completed" {
			return nil, errors.New("status must be 'active', 'inactive', or 'completed'")
		}
		term.Status = req.Status
	}

	if req.Description != "" {
		term.Description = req.Description
	}

	term.UpdatedAt = time.Now()

	if err := s.db.Save(&term).Error; err != nil {
		return nil, errors.New("failed to update term: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("AcademicSession").Preload("Creator").First(&term, term.ID).Error; err != nil {
		return nil, errors.New("failed to load term details: " + err.Error())
	}

	return s.toTermResponse(&term), nil
}

// DeleteTerm soft deletes a term
func (s *TermService) DeleteTerm(id string) error {
	termID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid term ID")
	}

	var term models.Term
	if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("term not found")
		}
		return errors.New("failed to fetch term: " + err.Error())
	}

	if err := s.db.Delete(&term).Error; err != nil {
		return errors.New("failed to delete term: " + err.Error())
	}

	return nil
}

// validateTermRequest validates the term request
func (s *TermService) validateTermRequest(req *dto.CreateTermRequest) error {
	if req.AcademicSessionID == "" {
		return errors.New("academic session ID is required")
	}
	if req.Name == "" {
		return errors.New("name is required")
	}
	if req.TermNumber < 1 || req.TermNumber > 3 {
		return errors.New("term number must be between 1 and 3")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "completed" {
		return errors.New("status must be 'active', 'inactive', or 'completed'")
	}
	return nil
}

// toTermResponse converts model to response DTO
func (s *TermService) toTermResponse(term *models.Term) *dto.TermResponse {
	response := &dto.TermResponse{
		ID:                term.ID.String(),
		AcademicSessionID: term.AcademicSessionID.String(),
		Name:              term.Name,
		Code:              term.Code,
		TermNumber:        term.TermNumber,
		StartDate:         term.StartDate,
		EndDate:           term.EndDate,
		IsCurrent:         term.IsCurrent,
		Status:            term.Status,
		Description:       term.Description,
		CreatedBy:         term.CreatedBy.String(),
		CreatedAt:         term.CreatedAt,
		UpdatedAt:         term.UpdatedAt,
	}

	// Add academic session details if preloaded
	if term.AcademicSession.ID != uuid.Nil {
		response.AcademicSession = &dto.AcademicSessionResponse{
			ID:           term.AcademicSession.ID.String(),
			AcademicYear: term.AcademicSession.AcademicYear,
			Code:         term.AcademicSession.Code,
			StartDate:    term.AcademicSession.StartDate,
			EndDate:      term.AcademicSession.EndDate,
			Status:       term.AcademicSession.Status,
			IsCurrent:    term.AcademicSession.IsCurrent,
			Description:  term.AcademicSession.Description,
		}
	}

	// Add creator details if preloaded
	if term.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        term.Creator.ID.String(),
			FirstName: term.Creator.FirstName,
			LastName:  term.Creator.LastName,
			Email:     term.Creator.Email,
			Phone:     term.Creator.Phone,
			Role:      term.Creator.Role,
		}
	}

	return response
}
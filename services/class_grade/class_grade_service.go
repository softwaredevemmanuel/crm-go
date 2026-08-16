// services/class_grade_service.go
package services

import (
	"errors"
	"strings"
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"crm-go/models"
	"crm-go/dto"
)

type ClassGradeService struct {
	db *gorm.DB
}

func NewClassGradeService(db *gorm.DB) *ClassGradeService {
	return &ClassGradeService{db: db}
}

// CreateClassGrade creates a new class grade
func (s *ClassGradeService) CreateClassGrade(req *dto.CreateClassGradeRequest, userID uuid.UUID) (*dto.ClassGradeResponse, error) {
	// Validate input
	if err := s.validateClassGradeRequest(req); err != nil {
		return nil, err
	}

	// Parse AcademicSession ID
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

	// Check if class grade with same name already exists
	var existing models.ClassGrade
	if err := s.db.Where("name = ?", req.Name).First(&existing).Error; err == nil {
		return nil, errors.New("class grade with this name already exists")
	}

	// Check if class grade with same code already exists
	if err := s.db.Where("code = ?", req.Code).First(&existing).Error; err == nil {
		return nil, errors.New("class grade with this code already exists")
	}

	// Check if class grade with same level and academic session already exists
	if err := s.db.Where("level = ? AND academic_session_id = ?", req.Level, academicSessionID).First(&existing).Error; err == nil {
		return nil, errors.New("class grade with this level and academic session already exists")
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create new class grade
	classGrade := &models.ClassGrade{
		ID:                uuid.New(),
		Name:              strings.TrimSpace(req.Name),
		Code:              strings.ToUpper(strings.TrimSpace(req.Code)),
		Level:             req.Level,
		Description:       strings.TrimSpace(req.Description),
		AcademicSessionID: academicSessionID,
		Capacity:          req.Capacity,
		Status:            status,
		CreatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	// Set default capacity if not provided
	if classGrade.Capacity == 0 {
		classGrade.Capacity = 30
	}

	// Save to database
	if err := s.db.Create(classGrade).Error; err != nil {
		return nil, errors.New("failed to create class grade: " + err.Error())
	}

	// Preload AcademicSession for response
	if err := s.db.Preload("AcademicSession").First(classGrade, classGrade.ID).Error; err != nil {
		return nil, errors.New("failed to load class grade details: " + err.Error())
	}

	// Convert to response DTO
	return s.toClassGradeResponse(classGrade), nil
}

// validateClassGradeRequest validates the class grade request
func (s *ClassGradeService) validateClassGradeRequest(req *dto.CreateClassGradeRequest) error {
	if req.Name == "" {
		return errors.New("class grade name is required")
	}
	if len(req.Name) < 2 {
		return errors.New("class grade name must be at least 2 characters")
	}
	if req.Code == "" {
		return errors.New("class grade code is required")
	}
	if len(req.Code) < 2 {
		return errors.New("class grade code must be at least 2 characters")
	}
	if req.Level < 1 || req.Level > 12 {
		return errors.New("level must be between 1 and 12")
	}
	if req.AcademicSessionID == "" {
		return errors.New("academic session is required")
	}
	if req.Capacity < 1 {
		return errors.New("capacity must be at least 1")
	}
	if req.Capacity > 100 {
		return errors.New("capacity cannot exceed 100")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "archived" {
		return errors.New("status must be 'active', 'inactive', or 'archived'")
	}
	return nil
}

// toClassGradeResponse converts model to response DTO
func (s *ClassGradeService) toClassGradeResponse(classGrade *models.ClassGrade) *dto.ClassGradeResponse {
	response := &dto.ClassGradeResponse{
		ID:          classGrade.ID.String(),
		Name:        classGrade.Name,
		Code:        classGrade.Code,
		Level:       classGrade.Level,
		Description: classGrade.Description,
		AcademicSessionID: classGrade.AcademicSessionID.String(),
		Capacity:    classGrade.Capacity,
		Status:      classGrade.Status,
		CreatedBy:   classGrade.CreatedBy.String(),
		CreatedAt:   classGrade.CreatedAt,
		UpdatedAt:   classGrade.UpdatedAt,
	}

	// Add academic session details if preloaded
	if classGrade.AcademicSession.ID != uuid.Nil {
		response.AcademicSession = &dto.AcademicSessionResponse{
			ID:          classGrade.AcademicSession.ID.String(),
			AcademicYear:        classGrade.AcademicSession.AcademicYear,
			Code:        classGrade.AcademicSession.Code,
			StartDate:   classGrade.AcademicSession.StartDate,
			EndDate:     classGrade.AcademicSession.EndDate,
			Status:      classGrade.AcademicSession.Status,
			IsCurrent:   classGrade.AcademicSession.IsCurrent,
			Description: classGrade.AcademicSession.Description,
		}
	}

	return response
}

// GetAllClassGrades retrieves all class grades with pagination and filters
func (s *ClassGradeService) GetAllClassGrades(params *dto.ClassGradeQueryParams) (*dto.ClassGradeListResponse, error) {
	// Set default values
	if params.Page == 0 {
		params.Page = 1
	}
	if params.Limit == 0 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "created_at"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.ClassGrade{}).Where("class_grades.deleted_at IS NULL")

	// Apply filters
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where(
			"LOWER(class_grades.name) LIKE ? OR LOWER(class_grades.code) LIKE ? OR LOWER(class_grades.description) LIKE ?",
			searchTerm, searchTerm, searchTerm,
		)
	}

	if params.Level > 0 {
		query = query.Where("class_grades.level = ?", params.Level)
	}

	if params.AcademicSessionID != "" {
		query = query.Where("class_grades.academic_session_id = ?", params.AcademicSessionID)
	}

	if params.Status != "" {
		query = query.Where("class_grades.status = ?", params.Status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count class grades: " + err.Error())
	}

	// Apply sorting
	sortColumn := params.SortBy
	if sortColumn == "" {
		sortColumn = "created_at"
	}
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order("class_grades." + sortColumn + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query with preload
	var classGrades []models.ClassGrade
	if err := query.Preload("AcademicSession").Find(&classGrades).Error; err != nil {
		return nil, errors.New("failed to fetch class grades: " + err.Error())
	}

	// Convert to response DTO
	var classGradeResponses []dto.ClassGradeResponse
	for _, classGrade := range classGrades {
		classGradeResponses = append(classGradeResponses, *s.toClassGradeResponse(&classGrade))
	}

	// Calculate total pages
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	return &dto.ClassGradeListResponse{
		ClassGrades: classGradeResponses,
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
	}, nil
}

// GetClassGradeByID retrieves a single class grade by ID
func (s *ClassGradeService) GetClassGradeByID(id string) (*dto.ClassGradeResponse, error) {
	classGradeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid class grade ID")
	}

	var classGrade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classGradeID).
		Preload("AcademicSession").
		First(&classGrade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class grade not found")
		}
		return nil, errors.New("failed to fetch class grade: " + err.Error())
	}

	return s.toClassGradeResponse(&classGrade), nil
}

// GetClassGradesByAcademicSession retrieves all class grades for a specific academic session
func (s *ClassGradeService) GetClassGradesByAcademicSession(academicSessionID string) ([]dto.ClassGradeResponse, error) {
	sessionID, err := uuid.Parse(academicSessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID")
	}

	var classGrades []models.ClassGrade
	if err := s.db.Where("academic_session_id = ? AND deleted_at IS NULL", sessionID).
		Preload("AcademicSession").
		Order("level ASC").
		Find(&classGrades).Error; err != nil {
		return nil, errors.New("failed to fetch class grades: " + err.Error())
	}

	responses := make([]dto.ClassGradeResponse, len(classGrades))
	for i, classGrade := range classGrades {
		responses[i] = *s.toClassGradeResponse(&classGrade)
	}

	return responses, nil
}

// GetLevels retrieves all unique levels
func (s *ClassGradeService) GetLevels() ([]int, error) {
	var levels []int
	if err := s.db.Model(&models.ClassGrade{}).
		Where("deleted_at IS NULL").
		Distinct("level").
		Order("level ASC").
		Pluck("level", &levels).Error; err != nil {
		return nil, errors.New("failed to fetch levels: " + err.Error())
	}
	return levels, nil
}

// UpdateClassGrade updates an existing class grade
func (s *ClassGradeService) UpdateClassGrade(id string, req *dto.UpdateClassGradeRequest, userID uuid.UUID) (*dto.ClassGradeResponse, error) {
	// Parse class grade ID
	classGradeID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid class grade ID")
	}

	// Find existing class grade
	var classGrade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classGradeID).First(&classGrade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class grade not found")
		}
		return nil, errors.New("failed to fetch class grade: " + err.Error())
	}

	// Check for duplicate code if being updated
	if req.Code != "" && req.Code != classGrade.Code {
		var existing models.ClassGrade
		if err := s.db.Where("code = ? AND id != ? AND deleted_at IS NULL", req.Code, classGradeID).First(&existing).Error; err == nil {
			return nil, errors.New("class grade with this code already exists")
		}
	}

	// Check for duplicate name if being updated
	if req.Name != "" && req.Name != classGrade.Name {
		var existing models.ClassGrade
		if err := s.db.Where("name = ? AND id != ? AND deleted_at IS NULL", req.Name, classGradeID).First(&existing).Error; err == nil {
			return nil, errors.New("class grade with this name already exists")
		}
	}

	// Check for duplicate level + academic session if being updated
	if req.Level > 0 && req.AcademicSessionID != "" {
		sessionID, err := uuid.Parse(req.AcademicSessionID)
		if err != nil {
			return nil, errors.New("invalid academic session ID format")
		}
		var existing models.ClassGrade
		if err := s.db.Where("level = ? AND academic_session_id = ? AND id != ? AND deleted_at IS NULL",
			req.Level, sessionID, classGradeID).First(&existing).Error; err == nil {
			return nil, errors.New("class grade with this level and academic session already exists")
		}
	}

	// Update fields
	if req.Name != "" {
		classGrade.Name = strings.TrimSpace(req.Name)
	}
	if req.Code != "" {
		classGrade.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	}
	if req.Level > 0 {
		classGrade.Level = req.Level
	}
	if req.Description != "" {
		classGrade.Description = strings.TrimSpace(req.Description)
	}
	if req.AcademicSessionID != "" {
		sessionID, err := uuid.Parse(req.AcademicSessionID)
		if err != nil {
			return nil, errors.New("invalid academic session ID format")
		}
		// Verify academic session exists
		var academicSession models.AcademicSession
		if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).First(&academicSession).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("academic session not found")
			}
			return nil, errors.New("failed to verify academic session: " + err.Error())
		}
		classGrade.AcademicSessionID = sessionID
	}
	if req.Capacity > 0 {
		classGrade.Capacity = req.Capacity
	}
	if req.Status != "" {
		classGrade.Status = req.Status
	}

	// Update timestamp
	classGrade.UpdatedAt = time.Now()

	// Save to database
	if err := s.db.Save(&classGrade).Error; err != nil {
		return nil, errors.New("failed to update class grade: " + err.Error())
	}

	// Preload AcademicSession for response
	if err := s.db.Preload("AcademicSession").First(&classGrade, classGrade.ID).Error; err != nil {
		return nil, errors.New("failed to load class grade details: " + err.Error())
	}

	// Convert to response DTO
	return s.toClassGradeResponse(&classGrade), nil
}

// DeleteClassGrade soft deletes a class grade
func (s *ClassGradeService) DeleteClassGrade(id string, userID uuid.UUID) error {
	// Parse class grade ID
	classGradeID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid class grade ID")
	}

	// Find existing class grade
	var classGrade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classGradeID).First(&classGrade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("class grade not found")
		}
		return errors.New("failed to fetch class grade: " + err.Error())
	}

	// Check if class grade has students enrolled (if you have a relationship)
	// Uncomment and implement if you have a Student model
	// var studentCount int64
	// if err := s.db.Model(&models.Student{}).Where("class_grade_id = ?", classGradeID).Count(&studentCount).Error; err != nil {
	// 	return errors.New("failed to check class grade usage: " + err.Error())
	// }
	// if studentCount > 0 {
	// 	return errors.New("cannot delete class grade: it has students enrolled")
	// }

	// Perform soft delete
	if err := s.db.Delete(&classGrade).Error; err != nil {
		return errors.New("failed to delete class grade: " + err.Error())
	}

	return nil
}
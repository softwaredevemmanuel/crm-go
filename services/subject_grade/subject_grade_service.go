// services/subject_grade_service.go
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

type SubjectGradeService struct {
	db *gorm.DB
}

func NewSubjectGradeService(db *gorm.DB) *SubjectGradeService {
	return &SubjectGradeService{db: db}
}

// CreateSubjectGrade creates a new subject-grade relationship
func (s *SubjectGradeService) CreateSubjectGrade(req *dto.CreateSubjectGradeRequest, userID uuid.UUID) (*dto.SubjectGradeResponse, error) {
	// Validate input
	if err := s.validateSubjectGradeRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// Check if relationship already exists
	var existing models.SubjectGrade
	if err := s.db.Where("subject_id = ? AND grade_id = ? AND academic_year = ?", 
		subjectID, gradeID, req.AcademicYear).First(&existing).Error; err == nil {
		return nil, errors.New("this subject-grade relationship already exists for the academic year")
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create relationship
	subjectGrade := &models.SubjectGrade{
		ID:           uuid.New(),
		SubjectID:    subjectID,
		GradeID:      gradeID,
		AcademicYear: strings.TrimSpace(req.AcademicYear),
		Status:       status,
		IsRequired:   req.IsRequired,
		Credits:      req.Credits,
		Description:  strings.TrimSpace(req.Description),
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Save to database
	if err := s.db.Create(subjectGrade).Error; err != nil {
		return nil, errors.New("failed to create subject-grade relationship: " + err.Error())
	}

	// Convert to response DTO
	return s.toSubjectGradeResponse(subjectGrade), nil
}

// BulkCreateSubjectGrades creates multiple subject-grade relationships
func (s *SubjectGradeService) BulkCreateSubjectGrades(req *dto.BulkCreateSubjectGradeRequest, userID uuid.UUID) (*dto.BulkCreateSubjectGradeResponse, error) {
	// Parse grade ID
	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("class grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	response := &dto.BulkCreateSubjectGradeResponse{
		Created: []dto.SubjectGradeResponse{},
		Errors:  []dto.BulkCreateError{},
	}

	for _, subjectIDStr := range req.SubjectIDs {
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "invalid subject ID",
			})
			continue
		}

		// Check if subject exists
		var subject models.Subject
		if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "subject not found",
			})
			continue
		}

		// Check if relationship already exists
		var existing models.SubjectGrade
		if err := s.db.Where("subject_id = ? AND grade_id = ? AND academic_year = ?", 
			subjectID, gradeID, req.AcademicYear).First(&existing).Error; err == nil {
			response.FailedCount++
			response.Errors = append(response.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "relationship already exists for this academic year",
			})
			continue
		}

		// Set default status
		status := req.Status
		if status == "" {
			status = "active"
		}

		// Create relationship
		subjectGrade := &models.SubjectGrade{
			ID:           uuid.New(),
			SubjectID:    subjectID,
			GradeID:      gradeID,
			AcademicYear: strings.TrimSpace(req.AcademicYear),
			Status:       status,
			IsRequired:   req.IsRequired,
			Credits:      req.Credits,
			Description:  strings.TrimSpace(req.Description),
			CreatedBy:    userID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.db.Create(subjectGrade).Error; err != nil {
			response.FailedCount++
			response.Errors = append(response.Errors, dto.BulkCreateError{
				SubjectID: subjectIDStr,
				Error:     "failed to create relationship: " + err.Error(),
			})
			continue
		}

		response.SuccessCount++
		response.Created = append(response.Created, *s.toSubjectGradeResponse(subjectGrade))
	}

	return response, nil
}

// GetAllSubjectGrades retrieves all subject-grade relationships with pagination and filters
func (s *SubjectGradeService) GetAllSubjectGrades(params *dto.SubjectGradeQueryParams) (*dto.SubjectGradeListResponse, error) {
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
	query := s.db.Model(&models.SubjectGrade{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("subject_id = ?", subjectID)
		}
	}

	if params.GradeID != "" {
		gradeID, err := uuid.Parse(params.GradeID)
		if err == nil {
			query = query.Where("grade_id = ?", gradeID)
		}
	}

	if params.AcademicYear != "" {
		query = query.Where("academic_year = ?", params.AcademicYear)
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.IsRequired != nil {
		query = query.Where("is_required = ?", *params.IsRequired)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, errors.New("failed to count subject-grade relationships: " + err.Error())
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
	query = query.Order(sortColumn + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute query with preload
	var subjectGrades []models.SubjectGrade
	if err := query.Preload("Subject").Preload("Grade").Find(&subjectGrades).Error; err != nil {
		return nil, errors.New("failed to fetch subject-grade relationships: " + err.Error())
	}

	// Convert to response DTO
	var subjectGradeResponses []dto.SubjectGradeResponse
	for _, sg := range subjectGrades {
		subjectGradeResponses = append(subjectGradeResponses, *s.toSubjectGradeResponse(&sg))
	}

	// Calculate total pages
	totalPages := int(total) / params.Limit
	if int(total)%params.Limit != 0 {
		totalPages++
	}

	return &dto.SubjectGradeListResponse{
		SubjectGrades: subjectGradeResponses,
		Total:         total,
		Page:          params.Page,
		Limit:         params.Limit,
		TotalPages:    totalPages,
	}, nil
}

// GetSubjectGradeByID retrieves a single subject-grade relationship by ID
func (s *SubjectGradeService) GetSubjectGradeByID(id string) (*dto.SubjectGradeResponse, error) {
	sgID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid subject-grade relationship ID")
	}

	var subjectGrade models.SubjectGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sgID).Preload("Subject").Preload("Grade").First(&subjectGrade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject-grade relationship not found")
		}
		return nil, errors.New("failed to fetch subject-grade relationship: " + err.Error())
	}

	return s.toSubjectGradeResponse(&subjectGrade), nil
}

// GetSubjectsByGrade retrieves all subjects for a specific grade
func (s *SubjectGradeService) GetSubjectsByGrade(gradeID string, academicYear string) ([]dto.SubjectGradeResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	query := s.db.Model(&models.SubjectGrade{}).
		Where("grade_id = ? AND deleted_at IS NULL", gID).
		Preload("Subject").
		Preload("Grade")

	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}

	var subjectGrades []models.SubjectGrade
	if err := query.Find(&subjectGrades).Error; err != nil {
		return nil, errors.New("failed to fetch subjects for grade: " + err.Error())
	}

	var responses []dto.SubjectGradeResponse
	for _, sg := range subjectGrades {
		responses = append(responses, *s.toSubjectGradeResponse(&sg))
	}

	return responses, nil
}

// GetGradesBySubject retrieves all grades for a specific subject
func (s *SubjectGradeService) GetGradesBySubject(subjectID string, academicYear string) ([]dto.SubjectGradeResponse, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	query := s.db.Model(&models.SubjectGrade{}).
		Where("subject_id = ? AND deleted_at IS NULL", sID).
		Preload("Subject").
		Preload("Grade")

	if academicYear != "" {
		query = query.Where("academic_year = ?", academicYear)
	}

	var subjectGrades []models.SubjectGrade
	if err := query.Find(&subjectGrades).Error; err != nil {
		return nil, errors.New("failed to fetch grades for subject: " + err.Error())
	}

	var responses []dto.SubjectGradeResponse
	for _, sg := range subjectGrades {
		responses = append(responses, *s.toSubjectGradeResponse(&sg))
	}

	return responses, nil
}

// DeleteSubjectGrade soft deletes a subject-grade relationship
func (s *SubjectGradeService) DeleteSubjectGrade(id string) error {
	sgID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid subject-grade relationship ID")
	}

	var subjectGrade models.SubjectGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", sgID).First(&subjectGrade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject-grade relationship not found")
		}
		return errors.New("failed to fetch subject-grade relationship: " + err.Error())
	}

	if err := s.db.Delete(&subjectGrade).Error; err != nil {
		return errors.New("failed to delete subject-grade relationship: " + err.Error())
	}

	return nil
}

// validateSubjectGradeRequest validates the subject-grade request
func (s *SubjectGradeService) validateSubjectGradeRequest(req *dto.CreateSubjectGradeRequest) error {
	if req.SubjectID == "" {
		return errors.New("subject ID is required")
	}
	if req.GradeID == "" {
		return errors.New("grade ID is required")
	}
	if req.AcademicYear == "" {
		return errors.New("academic year is required")
	}
	if !strings.Contains(req.AcademicYear, "/") {
		return errors.New("academic year must be in format YYYY/YYYY")
	}
	if req.Credits < 0 || req.Credits > 10 {
		return errors.New("credits must be between 0 and 10")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && req.Status != "archived" {
		return errors.New("status must be 'active', 'inactive', or 'archived'")
	}
	return nil
}

// toSubjectGradeResponse converts model to response DTO
func (s *SubjectGradeService) toSubjectGradeResponse(sg *models.SubjectGrade) *dto.SubjectGradeResponse {
	response := &dto.SubjectGradeResponse{
		ID:           sg.ID.String(),
		SubjectID:    sg.SubjectID.String(),
		GradeID:      sg.GradeID.String(),
		AcademicYear: sg.AcademicYear,
		Status:       sg.Status,
		IsRequired:   sg.IsRequired,
		Credits:      sg.Credits,
		Description:  sg.Description,
		CreatedBy:    sg.CreatedBy.String(),
		CreatedAt:    sg.CreatedAt,
		UpdatedAt:    sg.UpdatedAt,
	}

	// Add subject details if preloaded
	if sg.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:          sg.Subject.ID.String(),
			Name:        sg.Subject.Name,
			Code:        sg.Subject.Code,
			Description: sg.Subject.Description,
			Credits:     sg.Subject.Credits,
			Status:      sg.Subject.Status,
			CreatedAt:   sg.Subject.CreatedAt,
			UpdatedAt:   sg.Subject.UpdatedAt,
		}
	}

	// Add grade details if preloaded
	if sg.Grade.ID != uuid.Nil {
		response.Grade = &dto.ClassGradeResponse{
			ID:           sg.Grade.ID.String(),
			Name:         sg.Grade.Name,
			Code:         sg.Grade.Code,
			Level:        sg.Grade.Level,
			Description:  sg.Grade.Description,
			AcademicSessionID: sg.Grade.AcademicSessionID.String(),
			Capacity:     sg.Grade.Capacity,
			Status:       sg.Grade.Status,
			CreatedAt:    sg.Grade.CreatedAt,
			UpdatedAt:    sg.Grade.UpdatedAt,
		}
	}

	return response
}
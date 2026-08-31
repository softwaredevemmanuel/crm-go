// services/exam_service.go
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

type ExamService struct {
	db *gorm.DB
}

func NewExamService(db *gorm.DB) *ExamService {
	return &ExamService{db: db}
}

// CreateExam creates a new exam
func (s *ExamService) CreateExam(req *dto.CreateExamRequest, userID uuid.UUID) (*dto.ExamResponse, error) {
	// Validate input
	if err := s.validateExamRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	academicSessionID, err := uuid.Parse(req.AcademicSessionID)
	if err != nil {
		return nil, errors.New("invalid academic session ID format")
	}

	termID, err := uuid.Parse(req.TermID)
	if err != nil {
		return nil, errors.New("invalid term ID format")
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	classID, err := uuid.Parse(req.ClassID)
	if err != nil {
		return nil, errors.New("invalid class ID format")
	}



	// Verify all entities exist
	if err := s.verifyEntities(academicSessionID, termID, subjectID, classID); err != nil {
		return nil, err
	}

	// Parse exam date if provided
	var examDate *time.Time
	if req.ExamDate != "" {
		date, err := time.Parse("2006-01-02", req.ExamDate)
		if err != nil {
			return nil, errors.New("invalid exam date format. Use YYYY-MM-DD")
		}
		examDate = &date
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Create exam
	exam := &models.Exam{
		ID:                uuid.New(),
		AcademicSessionID: academicSessionID,
		TermID:            termID,
		SubjectID:         subjectID,
		ClassID:           classID,
		Title:             req.Title,
		ExamType:          req.ExamType,
		ExamDate:          examDate,
		Duration:          req.Duration,
		TotalMarks:        req.TotalMarks,
		Status:            status,
		CreatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(exam).Error; err != nil {
		return nil, errors.New("failed to create exam: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(exam, exam.ID).Error; err != nil {
		return nil, errors.New("failed to load exam details: " + err.Error())
	}

	return s.toExamResponse(exam), nil
}

// BulkCreateExams creates multiple exams
func (s *ExamService) BulkCreateExams(req *dto.BulkCreateExamsRequest, userID uuid.UUID) (*dto.BulkExamResult, error) {
	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	result := &dto.BulkExamResult{
		Created: []dto.ExamResponse{},
		Errors:  []dto.BulkExamError{},
	}

	for _, examReq := range req.Exams {
		// Validate
		if examReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: "title is required",
			})
			continue
		}

		// Parse UUIDs
		academicSessionID, err := uuid.Parse(examReq.AcademicSessionID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: "invalid academic session ID format",
			})
			continue
		}

		termID, err := uuid.Parse(examReq.TermID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: "invalid term ID format",
			})
			continue
		}

		subjectID, err := uuid.Parse(examReq.SubjectID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: "invalid subject ID format",
			})
			continue
		}

		classID, err := uuid.Parse(examReq.ClassID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: "invalid class ID format",
			})
			continue
		}


		// Verify entities
		if err := s.verifyEntities(academicSessionID, termID, subjectID, classID); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: err.Error(),
			})
			continue
		}

		// Parse exam date if provided
		var examDate *time.Time
		if examReq.ExamDate != "" {
			date, err := time.Parse("2006-01-02", examReq.ExamDate)
			if err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors, dto.BulkExamError{
					Title: examReq.Title,
					Error: "invalid exam date format. Use YYYY-MM-DD",
				})
				continue
			}
			examDate = &date
		}

		// Create exam
		exam := &models.Exam{
			ID:                uuid.New(),
			AcademicSessionID: academicSessionID,
			TermID:            termID,
			SubjectID:         subjectID,
			ClassID:           classID,
			Title:             examReq.Title,
			ExamType:          examReq.ExamType,
			ExamDate:          examDate,
			Duration:          examReq.Duration,
			TotalMarks:        examReq.TotalMarks,
			Status:            status,
			CreatedBy:         userID,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := s.db.Create(exam).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: "failed to create exam: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(exam, exam.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkExamError{
				Title: examReq.Title,
				Error: "failed to load exam details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toExamResponse(exam))
	}

	return result, nil
}

// GetAllExams retrieves all exams with pagination and filters
func (s *ExamService) GetAllExams(params *dto.ExamQueryParams) (*dto.ExamListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "exam_date"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.Exam{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.AcademicSessionID != "" {
		sessionID, err := uuid.Parse(params.AcademicSessionID)
		if err == nil {
			query = query.Where("academic_session_id = ?", sessionID)
		}
	}

	if params.TermID != "" {
		termID, err := uuid.Parse(params.TermID)
		if err == nil {
			query = query.Where("term_id = ?", termID)
		}
	}

	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("subject_id = ?", subjectID)
		}
	}

	if params.ClassID != "" {
		classID, err := uuid.Parse(params.ClassID)
		if err == nil {
			query = query.Where("class_id = ?", classID)
		}
	}


	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	if params.Search != "" {
		search := "%" + strings.ToLower(params.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(exam_type) LIKE ?", search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count exams: %w", err)
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
	var exams []models.Exam
	if err := query.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").Find(&exams).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exams: %w", err)
	}

	// Convert to response
	responses := make([]dto.ExamResponse, len(exams))
	for i, exam := range exams {
		responses[i] = *s.toExamResponse(&exam)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.ExamListResponse{
		Exams:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetExamByID retrieves a single exam by ID
func (s *ExamService) GetExamByID(id string) (*dto.ExamResponse, error) {
	examID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid exam ID")
	}

	var exam models.Exam
	if err := s.db.Where("id = ? AND deleted_at IS NULL", examID).
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
		Preload("Creator").
		First(&exam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("exam not found")
		}
		return nil, errors.New("failed to fetch exam: " + err.Error())
	}

	return s.toExamResponse(&exam), nil
}

// GetExamsBySubject retrieves all exams for a specific subject
func (s *ExamService) GetExamsBySubject(subjectID string) ([]dto.ExamResponse, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	var exams []models.Exam
	if err := s.db.Where("subject_id = ? AND deleted_at IS NULL", sID).
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
		Preload("Creator").
		Order("exam_date DESC").
		Find(&exams).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exams: %w", err)
	}

	responses := make([]dto.ExamResponse, len(exams))
	for i, exam := range exams {
		responses[i] = *s.toExamResponse(&exam)
	}

	return responses, nil
}

// GetExamsByClass retrieves all exams for a specific class
func (s *ExamService) GetExamsByClass(classID string) ([]dto.ExamResponse, error) {
	cID, err := uuid.Parse(classID)
	if err != nil {
		return nil, errors.New("invalid class ID")
	}

	var exams []models.Exam
	if err := s.db.Where("class_id = ? AND deleted_at IS NULL", cID).
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
		Preload("Creator").
		Order("exam_date DESC").
		Find(&exams).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch exams: %w", err)
	}

	responses := make([]dto.ExamResponse, len(exams))
	for i, exam := range exams {
		responses[i] = *s.toExamResponse(&exam)
	}

	return responses, nil
}

// UpdateExam updates an existing exam
func (s *ExamService) UpdateExam(id string, req *dto.UpdateExamRequest) (*dto.ExamResponse, error) {
	examID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid exam ID")
	}

	// Find existing exam
	var exam models.Exam
	if err := s.db.Where("id = ? AND deleted_at IS NULL", examID).First(&exam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("exam not found")
		}
		return nil, errors.New("failed to fetch exam: " + err.Error())
	}

	// Update fields
	if req.AcademicSessionID != "" {
		sessionID, err := uuid.Parse(req.AcademicSessionID)
		if err != nil {
			return nil, errors.New("invalid academic session ID format")
		}
		var session models.AcademicSession
		if err := s.db.Where("id = ? AND deleted_at IS NULL", sessionID).First(&session).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("academic session not found")
			}
			return nil, errors.New("failed to verify academic session: " + err.Error())
		}
		exam.AcademicSessionID = sessionID
	}

	if req.TermID != "" {
		termID, err := uuid.Parse(req.TermID)
		if err != nil {
			return nil, errors.New("invalid term ID format")
		}
		var term models.Term
		if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).First(&term).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("term not found")
			}
			return nil, errors.New("failed to verify term: " + err.Error())
		}
		exam.TermID = termID
	}

	if req.SubjectID != "" {
		subjectID, err := uuid.Parse(req.SubjectID)
		if err != nil {
			return nil, errors.New("invalid subject ID format")
		}
		var subject models.Subject
		if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("subject not found")
			}
			return nil, errors.New("failed to verify subject: " + err.Error())
		}
		exam.SubjectID = subjectID
	}

	if req.ClassID != "" {
		classID, err := uuid.Parse(req.ClassID)
		if err != nil {
			return nil, errors.New("invalid class ID format")
		}
		var class models.ClassGrade
		if err := s.db.Where("id = ? AND deleted_at IS NULL", classID).First(&class).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("class not found")
			}
			return nil, errors.New("failed to verify class: " + err.Error())
		}
		exam.ClassID = classID
	}



	if req.Title != "" {
		exam.Title = req.Title
	}

	if req.ExamType != "" {
		exam.ExamType = req.ExamType
	}

	if req.ExamDate != "" {
		date, err := time.Parse("2006-01-02", req.ExamDate)
		if err != nil {
			return nil, errors.New("invalid exam date format. Use YYYY-MM-DD")
		}
		exam.ExamDate = &date
	}

	if req.Duration > 0 {
		exam.Duration = req.Duration
	}

	if req.TotalMarks > 0 {
		exam.TotalMarks = req.TotalMarks
	}

	if req.Status != "" {
		if req.Status != "draft" && req.Status != "published" && req.Status != "completed" {
			return nil, errors.New("status must be 'draft', 'published', or 'completed'")
		}
		exam.Status = req.Status
	}

	exam.UpdatedAt = time.Now()

	if err := s.db.Save(&exam).Error; err != nil {
		return nil, errors.New("failed to update exam: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(&exam, exam.ID).Error; err != nil {
		return nil, errors.New("failed to load exam details: " + err.Error())
	}

	return s.toExamResponse(&exam), nil
}

// DeleteExam soft deletes an exam
func (s *ExamService) DeleteExam(id string) error {
	examID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid exam ID")
	}

	var exam models.Exam
	if err := s.db.Where("id = ? AND deleted_at IS NULL", examID).First(&exam).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("exam not found")
		}
		return errors.New("failed to fetch exam: " + err.Error())
	}

	if err := s.db.Delete(&exam).Error; err != nil {
		return errors.New("failed to delete exam: " + err.Error())
	}

	return nil
}

// verifyEntities verifies that all referenced entities exist
func (s *ExamService) verifyEntities(academicSessionID, termID, subjectID, classID uuid.UUID) error {
	// Check academic session
	var academicSession models.AcademicSession
	if err := s.db.Where("id = ? AND deleted_at IS NULL", academicSessionID).First(&academicSession).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("academic session not found")
		}
		return errors.New("failed to verify academic session: " + err.Error())
	}

	// Check term
	var term models.Term
	if err := s.db.Where("id = ? AND deleted_at IS NULL", termID).First(&term).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("term not found")
		}
		return errors.New("failed to verify term: " + err.Error())
	}

	// Check subject
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("subject not found")
		}
		return errors.New("failed to verify subject: " + err.Error())
	}

	// Check class
	var class models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", classID).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("class not found")
		}
		return errors.New("failed to verify class: " + err.Error())
	}

	return nil
}

// validateExamRequest validates the exam request
func (s *ExamService) validateExamRequest(req *dto.CreateExamRequest) error {
	if req.AcademicSessionID == "" {
		return errors.New("academic session ID is required")
	}
	if req.TermID == "" {
		return errors.New("term ID is required")
	}
	if req.SubjectID == "" {
		return errors.New("subject ID is required")
	}
	if req.ClassID == "" {
		return errors.New("class ID is required")
	}
	if req.Title == "" {
		return errors.New("title is required")
	}
	if req.Status != "" && req.Status != "draft" && req.Status != "published" && req.Status != "completed" {
		return errors.New("status must be 'draft', 'published', or 'completed'")
	}
	if req.TotalMarks < 0 {
		return errors.New("total marks cannot be negative")
	}
	return nil
}

// toExamResponse converts model to response DTO
func (s *ExamService) toExamResponse(exam *models.Exam) *dto.ExamResponse {
	response := &dto.ExamResponse{
		ID:                exam.ID.String(),
		AcademicSessionID: exam.AcademicSessionID.String(),
		TermID:            exam.TermID.String(),
		SubjectID:         exam.SubjectID.String(),
		ClassID:           exam.ClassID.String(),
		Title:             exam.Title,
		ExamType:          exam.ExamType,
		ExamDate:          exam.ExamDate,
		Duration:          exam.Duration,
		TotalMarks:        exam.TotalMarks,
		Status:            exam.Status,
		CreatedBy:         exam.CreatedBy.String(),
		CreatedAt:         exam.CreatedAt,
		UpdatedAt:         exam.UpdatedAt,
	}

	// Add academic session details if preloaded
	if exam.AcademicSession.ID != uuid.Nil {
		response.AcademicSession = &dto.AcademicSessionResponse{
			ID:           exam.AcademicSession.ID.String(),
			AcademicYear: exam.AcademicSession.AcademicYear,
			Code:         exam.AcademicSession.Code,
			StartDate:    exam.AcademicSession.StartDate,
			EndDate:      exam.AcademicSession.EndDate,
			Status:       exam.AcademicSession.Status,
			IsCurrent:    exam.AcademicSession.IsCurrent,
		}
	}

	// Add term details if preloaded
	if exam.Term.ID != uuid.Nil {
		response.Term = &dto.TermResponse{
			ID:         exam.Term.ID.String(),
			Name:       exam.Term.Name,
			Code:       exam.Term.Code,
			TermNumber: exam.Term.TermNumber,
			StartDate:  exam.Term.StartDate,
			EndDate:    exam.Term.EndDate,
			IsCurrent:  exam.Term.IsCurrent,
			Status:     exam.Term.Status,
		}
	}

	// Add subject details if preloaded
	if exam.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:           exam.Subject.ID.String(),
			Name:         exam.Subject.Name,
			Code:         exam.Subject.Code,
			Description:  exam.Subject.Description,
			DepartmentID: exam.Subject.DepartmentID.String(),
		}
	}

	// Add class details if preloaded
	if exam.Class.ID != uuid.Nil {
		response.Class = &dto.ClassGradeResponse{
			ID:          exam.Class.ID.String(),
			Name:        exam.Class.Name,
			Code:        exam.Class.Code,
			Level:       exam.Class.Level,
			Description: exam.Class.Description,
			Status:      exam.Class.Status,
		}
	}



	// Add creator details if preloaded
	if exam.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        exam.Creator.ID.String(),
			FirstName: exam.Creator.FirstName,
			LastName:  exam.Creator.LastName,
			Email:     exam.Creator.Email,
			Phone:     exam.Creator.Phone,
			Role:      exam.Creator.Role,
		}
	}

	return response
}
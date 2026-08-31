// services/test_service.go
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

type TestService struct {
	db *gorm.DB
}

func NewTestService(db *gorm.DB) *TestService {
	return &TestService{db: db}
}

// CreateTest creates a new test
func (s *TestService) CreateTest(req *dto.CreateTestRequest, userID uuid.UUID) (*dto.TestResponse, error) {
	// Validate input
	if err := s.validateTestRequest(req); err != nil {
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

	// Parse test date if provided
	var testDate *time.Time
	if req.TestDate != "" {
		date, err := time.Parse("2006-01-02", req.TestDate)
		if err != nil {
			return nil, errors.New("invalid test date format. Use YYYY-MM-DD")
		}
		testDate = &date
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// Create test
	test := &models.Test{
		ID:                uuid.New(),
		AcademicSessionID: academicSessionID,
		TermID:            termID,
		SubjectID:         subjectID,
		ClassID:           classID,
		Title:             req.Title,
		TestType:          req.TestType,
		TestDate:          testDate,
		Duration:          req.Duration,
		TotalMarks:        req.TotalMarks,
		Status:            status,
		CreatedBy:         userID,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.db.Create(test).Error; err != nil {
		return nil, errors.New("failed to create test: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(test, test.ID).Error; err != nil {
		return nil, errors.New("failed to load test details: " + err.Error())
	}

	return s.toTestResponse(test), nil
}

// BulkCreateTests creates multiple tests
func (s *TestService) BulkCreateTests(req *dto.BulkCreateTestsRequest, userID uuid.UUID) (*dto.BulkTestResult, error) {
	// Set default status
	status := req.Status
	if status == "" {
		status = "draft"
	}

	result := &dto.BulkTestResult{
		Created: []dto.TestResponse{},
		Errors:  []dto.BulkTestError{},
	}

	for _, testReq := range req.Tests {
		// Validate
		if testReq.Title == "" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: "title is required",
			})
			continue
		}

		// Parse UUIDs
		academicSessionID, err := uuid.Parse(testReq.AcademicSessionID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: "invalid academic session ID format",
			})
			continue
		}

		termID, err := uuid.Parse(testReq.TermID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: "invalid term ID format",
			})
			continue
		}

		subjectID, err := uuid.Parse(testReq.SubjectID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: "invalid subject ID format",
			})
			continue
		}

		classID, err := uuid.Parse(testReq.ClassID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: "invalid class ID format",
			})
			continue
		}

		// Verify entities
		if err := s.verifyEntities(academicSessionID, termID, subjectID, classID); err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: err.Error(),
			})
			continue
		}

		// Parse test date if provided
		var testDate *time.Time
		if testReq.TestDate != "" {
			date, err := time.Parse("2006-01-02", testReq.TestDate)
			if err != nil {
				result.FailedCount++
				result.Errors = append(result.Errors, dto.BulkTestError{
					Title: testReq.Title,
					Error: "invalid test date format. Use YYYY-MM-DD",
				})
				continue
			}
			testDate = &date
		}

		// Create test
		test := &models.Test{
			ID:                uuid.New(),
			AcademicSessionID: academicSessionID,
			TermID:            termID,
			SubjectID:         subjectID,
			ClassID:           classID,
			Title:             testReq.Title,
			TestType:          testReq.TestType,
			TestDate:          testDate,
			Duration:          testReq.Duration,
			TotalMarks:        testReq.TotalMarks,
			Status:            status,
			CreatedBy:         userID,
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
		}

		if err := s.db.Create(test).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: "failed to create test: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(test, test.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkTestError{
				Title: testReq.Title,
				Error: "failed to load test details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toTestResponse(test))
	}

	return result, nil
}

// GetAllTests retrieves all tests with pagination and filters
func (s *TestService) GetAllTests(params *dto.TestQueryParams) (*dto.TestListResponse, error) {
	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 || params.Limit > 100 {
		params.Limit = 20
	}
	if params.SortBy == "" {
		params.SortBy = "test_date"
	}
	if params.SortOrder == "" {
		params.SortOrder = "desc"
	}

	// Build query
	query := s.db.Model(&models.Test{}).Where("deleted_at IS NULL")

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
		query = query.Where("LOWER(title) LIKE ? OR LOWER(test_type) LIKE ?", search, search)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count tests: %w", err)
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
	var tests []models.Test
	if err := query.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").Find(&tests).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tests: %w", err)
	}

	// Convert to response
	responses := make([]dto.TestResponse, len(tests))
	for i, test := range tests {
		responses[i] = *s.toTestResponse(&test)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.TestListResponse{
		Tests:      responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetTestByID retrieves a single test by ID
func (s *TestService) GetTestByID(id string) (*dto.TestResponse, error) {
	testID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid test ID")
	}

	var test models.Test
	if err := s.db.Where("id = ? AND deleted_at IS NULL", testID).
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
		Preload("Creator").
		First(&test).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test not found")
		}
		return nil, errors.New("failed to fetch test: " + err.Error())
	}

	return s.toTestResponse(&test), nil
}

// GetTestsBySubject retrieves all tests for a specific subject
func (s *TestService) GetTestsBySubject(subjectID string) ([]dto.TestResponse, error) {
	sID, err := uuid.Parse(subjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID")
	}

	var tests []models.Test
	if err := s.db.Where("subject_id = ? AND deleted_at IS NULL", sID).
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
		Preload("Creator").
		Order("test_date DESC").
		Find(&tests).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tests: %w", err)
	}

	responses := make([]dto.TestResponse, len(tests))
	for i, test := range tests {
		responses[i] = *s.toTestResponse(&test)
	}

	return responses, nil
}

// GetTestsByClass retrieves all tests for a specific class
func (s *TestService) GetTestsByClass(classID string) ([]dto.TestResponse, error) {
	cID, err := uuid.Parse(classID)
	if err != nil {
		return nil, errors.New("invalid class ID")
	}

	var tests []models.Test
	if err := s.db.Where("class_id = ? AND deleted_at IS NULL", cID).
		Preload("AcademicSession").
		Preload("Term").
		Preload("Subject").
		Preload("Class").
		Preload("Creator").
		Order("test_date DESC").
		Find(&tests).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch tests: %w", err)
	}

	responses := make([]dto.TestResponse, len(tests))
	for i, test := range tests {
		responses[i] = *s.toTestResponse(&test)
	}

	return responses, nil
}

// UpdateTest updates an existing test
func (s *TestService) UpdateTest(id string, req *dto.UpdateTestRequest) (*dto.TestResponse, error) {
	testID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid test ID")
	}

	// Find existing test
	var test models.Test
	if err := s.db.Where("id = ? AND deleted_at IS NULL", testID).First(&test).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("test not found")
		}
		return nil, errors.New("failed to fetch test: " + err.Error())
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
		test.AcademicSessionID = sessionID
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
		test.TermID = termID
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
		test.SubjectID = subjectID
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
		test.ClassID = classID
	}


	if req.Title != "" {
		test.Title = req.Title
	}

	if req.TestType != "" {
		test.TestType = req.TestType
	}

	if req.TestDate != "" {
		date, err := time.Parse("2006-01-02", req.TestDate)
		if err != nil {
			return nil, errors.New("invalid test date format. Use YYYY-MM-DD")
		}
		test.TestDate = &date
	}

	if req.Duration > 0 {
		test.Duration = req.Duration
	}

	if req.TotalMarks > 0 {
		test.TotalMarks = req.TotalMarks
	}

	if req.Status != "" {
		if req.Status != "draft" && req.Status != "published" && req.Status != "completed" {
			return nil, errors.New("status must be 'draft', 'published', or 'completed'")
		}
		test.Status = req.Status
	}

	test.UpdatedAt = time.Now()

	if err := s.db.Save(&test).Error; err != nil {
		return nil, errors.New("failed to update test: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("AcademicSession").Preload("Term").Preload("Subject").Preload("Class").Preload("Creator").First(&test, test.ID).Error; err != nil {
		return nil, errors.New("failed to load test details: " + err.Error())
	}

	return s.toTestResponse(&test), nil
}

// DeleteTest soft deletes a test
func (s *TestService) DeleteTest(id string) error {
	testID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid test ID")
	}

	var test models.Test
	if err := s.db.Where("id = ? AND deleted_at IS NULL", testID).First(&test).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("test not found")
		}
		return errors.New("failed to fetch test: " + err.Error())
	}

	if err := s.db.Delete(&test).Error; err != nil {
		return errors.New("failed to delete test: " + err.Error())
	}

	return nil
}

// verifyEntities verifies that all referenced entities exist
func (s *TestService) verifyEntities(academicSessionID, termID, subjectID, classID uuid.UUID) error {
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


// validateTestRequest validates the test request
func (s *TestService) validateTestRequest(req *dto.CreateTestRequest) error {
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

// toTestResponse converts model to response DTO
func (s *TestService) toTestResponse(test *models.Test) *dto.TestResponse {
	response := &dto.TestResponse{
		ID:                test.ID.String(),
		AcademicSessionID: test.AcademicSessionID.String(),
		TermID:            test.TermID.String(),
		SubjectID:         test.SubjectID.String(),
		ClassID:           test.ClassID.String(),
		Title:             test.Title,
		TestType:          test.TestType,
		TestDate:          test.TestDate,
		Duration:          test.Duration,
		TotalMarks:        test.TotalMarks,
		Status:            test.Status,
		CreatedBy:         test.CreatedBy.String(),
		CreatedAt:         test.CreatedAt,
		UpdatedAt:         test.UpdatedAt,
	}

	// Add academic session details if preloaded
	if test.AcademicSession.ID != uuid.Nil {
		response.AcademicSession = &dto.AcademicSessionResponse{
			ID:           test.AcademicSession.ID.String(),
			AcademicYear: test.AcademicSession.AcademicYear,
			Code:         test.AcademicSession.Code,
			StartDate:    test.AcademicSession.StartDate,
			EndDate:      test.AcademicSession.EndDate,
			Status:       test.AcademicSession.Status,
			IsCurrent:    test.AcademicSession.IsCurrent,
		}
	}

	// Add term details if preloaded
	if test.Term.ID != uuid.Nil {
		response.Term = &dto.TermResponse{
			ID:         test.Term.ID.String(),
			Name:       test.Term.Name,
			Code:       test.Term.Code,
			TermNumber: test.Term.TermNumber,
			StartDate:  test.Term.StartDate,
			EndDate:    test.Term.EndDate,
			IsCurrent:  test.Term.IsCurrent,
			Status:     test.Term.Status,
		}
	}

	// Add subject details if preloaded
	if test.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:           test.Subject.ID.String(),
			Name:         test.Subject.Name,
			Code:         test.Subject.Code,
			Description:  test.Subject.Description,
			DepartmentID: test.Subject.DepartmentID.String(),
		}
	}

	// Add class details if preloaded
	if test.Class.ID != uuid.Nil {
		response.Class = &dto.ClassGradeResponse{
			ID:          test.Class.ID.String(),
			Name:        test.Class.Name,
			Code:        test.Class.Code,
			Level:       test.Class.Level,
			Description: test.Class.Description,
			Status:      test.Class.Status,
		}
	}


	// Add creator details if preloaded
	if test.Creator.ID != uuid.Nil {
		response.Creator = &dto.UserResponse{
			ID:        test.Creator.ID.String(),
			FirstName: test.Creator.FirstName,
			LastName:  test.Creator.LastName,
			Email:     test.Creator.Email,
			Phone:     test.Creator.Phone,
			Role:      test.Creator.Role,
		}
	}

	return response
}
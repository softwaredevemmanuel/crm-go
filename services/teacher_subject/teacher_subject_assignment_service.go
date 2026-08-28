// services/teacher_subject_assignment_service.go
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

type TeacherSubjectAssignmentService struct {
	db *gorm.DB
}

func NewTeacherSubjectAssignmentService(db *gorm.DB) *TeacherSubjectAssignmentService {
	return &TeacherSubjectAssignmentService{db: db}
}

// CreateAssignment creates a new subject assignment for a teacher
func (s *TeacherSubjectAssignmentService) CreateAssignment(req *dto.CreateTeacherSubjectAssignmentRequest) (*dto.TeacherSubjectAssignmentResponse, error) {
	// Validate input
	if err := s.validateAssignmentRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	subjectID, err := uuid.Parse(req.SubjectID)
	if err != nil {
		return nil, errors.New("invalid subject ID format")
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		return nil, errors.New("invalid teacher ID format")
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// Check if subject exists
	var subject models.Subject
	if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("subject not found")
		}
		return nil, errors.New("failed to verify subject: " + err.Error())
	}

	// Check if teacher exists and is a teacher
	var teacher models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", teacherID).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("failed to verify teacher: " + err.Error())
	}

	if teacher.Position != "teacher" && teacher.Position != "admin" {
		return nil, errors.New("user is not a teacher or admin")
	}

	// Check if assignment already exists
	var existing models.TeacherSubjectAssignment
	if err := s.db.Where("grade_id = ? AND subject_id = ? AND teacher_id = ? AND deleted_at IS NULL",
		gradeID, subjectID, teacherID).First(&existing).Error; err == nil {
		return nil, errors.New("this subject is already assigned to this teacher for this grade")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create assignment
	assignment := &models.TeacherSubjectAssignment{
		ID:        uuid.New(),
		GradeID:   gradeID,
		SubjectID: subjectID,
		TeacherID: teacherID,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(assignment).Error; err != nil {
		return nil, errors.New("failed to create assignment: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Grade").Preload("Subject").Preload("Teacher").First(assignment, assignment.ID).Error; err != nil {
		return nil, errors.New("failed to load assignment details: " + err.Error())
	}

	return s.toAssignmentResponse(assignment), nil
}

// BulkAssignSubjects assigns multiple subjects to a teacher for a grade
func (s *TeacherSubjectAssignmentService) BulkAssignSubjects(req *dto.BulkAssignSubjectsRequest) (*dto.BulkAssignmentResult, error) {
	// Parse UUIDs
	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		return nil, errors.New("invalid teacher ID format")
	}

	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	// Check if teacher exists
	var teacher models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", teacherID).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("failed to verify teacher: " + err.Error())
	}

	if teacher.Position != "teacher" && teacher.Position != "admin" {
		return nil, errors.New("user is not a teacher or admin")
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	result := &dto.BulkAssignmentResult{
		Created: []dto.TeacherSubjectAssignmentResponse{},
		Errors:  []dto.BulkAssignmentError{},
	}

	for _, subjectIDStr := range req.SubjectIDs {
		subjectID, err := uuid.Parse(subjectIDStr)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignmentError{
				SubjectID: subjectIDStr,
				Error:     "invalid subject ID format",
			})
			continue
		}

		// Check if subject exists
		var subject models.Subject
		if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignmentError{
				SubjectID: subjectIDStr,
				Error:     "subject not found",
			})
			continue
		}

		// Check if assignment already exists
		var existing models.TeacherSubjectAssignment
		if err := s.db.Where("grade_id = ? AND subject_id = ? AND teacher_id = ? AND deleted_at IS NULL",
			gradeID, subjectID, teacherID).First(&existing).Error; err == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignmentError{
				SubjectID: subjectIDStr,
				Error:     "subject already assigned to this teacher for this grade",
			})
			continue
		}

		// Create assignment
		assignment := &models.TeacherSubjectAssignment{
			ID:        uuid.New(),
			GradeID:   gradeID,
			SubjectID: subjectID,
			TeacherID: teacherID,
			Status:    status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.db.Create(assignment).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignmentError{
				SubjectID: subjectIDStr,
				Error:     "failed to create assignment: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Grade").Preload("Subject").Preload("Teacher").First(assignment, assignment.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignmentError{
				SubjectID: subjectIDStr,
				Error:     "failed to load assignment details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toAssignmentResponse(assignment))
	}

	return result, nil
}

// GetAllAssignments retrieves all assignments with pagination and filters
func (s *TeacherSubjectAssignmentService) GetAllAssignments(params *dto.TeacherSubjectAssignmentQueryParams) (*dto.TeacherSubjectAssignmentListResponse, error) {
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
	query := s.db.Model(&models.TeacherSubjectAssignment{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.GradeID != "" {
		gradeID, err := uuid.Parse(params.GradeID)
		if err == nil {
			query = query.Where("grade_id = ?", gradeID)
		}
	}

	if params.SubjectID != "" {
		subjectID, err := uuid.Parse(params.SubjectID)
		if err == nil {
			query = query.Where("subject_id = ?", subjectID)
		}
	}

	if params.TeacherID != "" {
		teacherID, err := uuid.Parse(params.TeacherID)
		if err == nil {
			query = query.Where("teacher_id = ?", teacherID)
		}
	}

	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count assignments: %w", err)
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
	var assignments []models.TeacherSubjectAssignment
	if err := query.Preload("Grade").Preload("Subject").Preload("Teacher").Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch assignments: %w", err)
	}

	// Convert to response
	responses := make([]dto.TeacherSubjectAssignmentResponse, len(assignments))
	for i, assignment := range assignments {
		responses[i] = *s.toAssignmentResponse(&assignment)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.TeacherSubjectAssignmentListResponse{
		Assignments: responses,
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
	}, nil
}

// GetAssignmentByID retrieves a single assignment by ID
func (s *TeacherSubjectAssignmentService) GetAssignmentByID(id string) (*dto.TeacherSubjectAssignmentResponse, error) {
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid assignment ID")
	}

	var assignment models.TeacherSubjectAssignment
	if err := s.db.Where("id = ? AND deleted_at IS NULL", assignmentID).
		Preload("Grade").
		Preload("Subject").
		Preload("Teacher").
		First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("assignment not found")
		}
		return nil, errors.New("failed to fetch assignment: " + err.Error())
	}

	return s.toAssignmentResponse(&assignment), nil
}

// GetAssignmentsByTeacher retrieves all assignments for a specific teacher
func (s *TeacherSubjectAssignmentService) GetAssignmentsByTeacher(teacherID string) ([]dto.TeacherSubjectAssignmentResponse, error) {
	tID, err := uuid.Parse(teacherID)
	if err != nil {
		return nil, errors.New("invalid teacher ID")
	}

	var assignments []models.TeacherSubjectAssignment
	if err := s.db.Where("teacher_id = ? AND deleted_at IS NULL AND status = ?", tID, "active").
		Preload("Grade").
		Preload("Subject").
		Preload("Teacher").
		Order("created_at DESC").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch assignments for teacher: %w", err)
	}

	responses := make([]dto.TeacherSubjectAssignmentResponse, len(assignments))
	for i, assignment := range assignments {
		responses[i] = *s.toAssignmentResponse(&assignment)
	}

	return responses, nil
}

// GetAssignmentsByGrade retrieves all assignments for a specific grade
func (s *TeacherSubjectAssignmentService) GetAssignmentsByGrade(gradeID string) ([]dto.TeacherSubjectAssignmentResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	var assignments []models.TeacherSubjectAssignment
	if err := s.db.Where("grade_id = ? AND deleted_at IS NULL", gID).
		Preload("Grade").
		Preload("Subject").
		Preload("Teacher").
		Order("created_at DESC").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch assignments for grade: %w", err)
	}

	responses := make([]dto.TeacherSubjectAssignmentResponse, len(assignments))
	for i, assignment := range assignments {
		responses[i] = *s.toAssignmentResponse(&assignment)
	}

	return responses, nil
}

// UpdateAssignment updates an existing assignment
func (s *TeacherSubjectAssignmentService) UpdateAssignment(id string, req *dto.UpdateTeacherSubjectAssignmentRequest) (*dto.TeacherSubjectAssignmentResponse, error) {
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid assignment ID")
	}

	// Find existing assignment
	var assignment models.TeacherSubjectAssignment
	if err := s.db.Where("id = ? AND deleted_at IS NULL", assignmentID).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("assignment not found")
		}
		return nil, errors.New("failed to fetch assignment: " + err.Error())
	}

	// Update fields
	if req.GradeID != "" {
		gradeID, err := uuid.Parse(req.GradeID)
		if err != nil {
			return nil, errors.New("invalid grade ID format")
		}
		// Verify grade exists
		var grade models.ClassGrade
		if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("grade not found")
			}
			return nil, errors.New("failed to verify grade: " + err.Error())
		}
		assignment.GradeID = gradeID
	}

	if req.SubjectID != "" {
		subjectID, err := uuid.Parse(req.SubjectID)
		if err != nil {
			return nil, errors.New("invalid subject ID format")
		}
		// Verify subject exists
		var subject models.Subject
		if err := s.db.Where("id = ? AND deleted_at IS NULL", subjectID).First(&subject).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("subject not found")
			}
			return nil, errors.New("failed to verify subject: " + err.Error())
		}
		assignment.SubjectID = subjectID
	}

	if req.TeacherID != "" {
		teacherID, err := uuid.Parse(req.TeacherID)
		if err != nil {
			return nil, errors.New("invalid teacher ID format")
		}
		// Verify teacher exists
		var teacher models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", teacherID).First(&teacher).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("teacher not found")
			}
			return nil, errors.New("failed to verify teacher: " + err.Error())
		}
		if teacher.Position != "teacher" && teacher.Position != "admin" {
			return nil, errors.New("user is not a teacher or admin")
		}
		assignment.TeacherID = teacherID
	}

	if req.Status != "" {
		if req.Status != "active" && req.Status != "inactive" {
			return nil, errors.New("status must be 'active' or 'inactive'")
		}
		assignment.Status = req.Status
	}

	assignment.UpdatedAt = time.Now()

	if err := s.db.Save(&assignment).Error; err != nil {
		return nil, errors.New("failed to update assignment: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Grade").Preload("Subject").Preload("Teacher").First(&assignment, assignment.ID).Error; err != nil {
		return nil, errors.New("failed to load assignment details: " + err.Error())
	}

	return s.toAssignmentResponse(&assignment), nil
}

// DeleteAssignment soft deletes an assignment
func (s *TeacherSubjectAssignmentService) DeleteAssignment(id string) error {
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid assignment ID")
	}

	var assignment models.TeacherSubjectAssignment
	if err := s.db.Where("id = ? AND deleted_at IS NULL", assignmentID).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("assignment not found")
		}
		return errors.New("failed to fetch assignment: " + err.Error())
	}

	if err := s.db.Delete(&assignment).Error; err != nil {
		return errors.New("failed to delete assignment: " + err.Error())
	}

	return nil
}

// validateAssignmentRequest validates the assignment request
func (s *TeacherSubjectAssignmentService) validateAssignmentRequest(req *dto.CreateTeacherSubjectAssignmentRequest) error {
	if req.GradeID == "" {
		return errors.New("grade ID is required")
	}
	if req.SubjectID == "" {
		return errors.New("subject ID is required")
	}
	if req.TeacherID == "" {
		return errors.New("teacher ID is required")
	}
	if req.Status != "" && req.Status != "active" && req.Status != "inactive" {
		return errors.New("status must be 'active' or 'inactive'")
	}
	return nil
}

// toAssignmentResponse converts model to response DTO
func (s *TeacherSubjectAssignmentService) toAssignmentResponse(assignment *models.TeacherSubjectAssignment) *dto.TeacherSubjectAssignmentResponse {
	response := &dto.TeacherSubjectAssignmentResponse{
		ID:        assignment.ID.String(),
		GradeID:   assignment.GradeID.String(),
		SubjectID: assignment.SubjectID.String(),
		TeacherID: assignment.TeacherID.String(),
		Status:    assignment.Status,
		CreatedAt: assignment.CreatedAt,
		UpdatedAt: assignment.UpdatedAt,
	}

	// Add grade details if preloaded
	if assignment.Grade.ID != uuid.Nil {
		response.Grade = &dto.ClassGradeResponse{
			ID:          assignment.Grade.ID.String(),
			Name:        assignment.Grade.Name,
			Code:        assignment.Grade.Code,
			Level:       assignment.Grade.Level,
			Description: assignment.Grade.Description,
			Capacity:    assignment.Grade.Capacity,
			Status:      assignment.Grade.Status,
			CreatedAt:   assignment.Grade.CreatedAt,
			UpdatedAt:   assignment.Grade.UpdatedAt,
		}
	}

	// Add subject details if preloaded
	if assignment.Subject.ID != uuid.Nil {
		response.Subject = &dto.SubjectResponse{
			ID:          assignment.Subject.ID.String(),
			Name:        assignment.Subject.Name,
			Code:        assignment.Subject.Code,
			Description: assignment.Subject.Description,
			CreatedAt:   assignment.Subject.CreatedAt,
			UpdatedAt:   assignment.Subject.UpdatedAt,
		}
	}

	// Add teacher details if preloaded
	if assignment.Teacher.ID != uuid.Nil {
		response.Teacher = &dto.UserResponse{
			ID:        assignment.Teacher.ID.String(),
			FirstName: assignment.Teacher.FirstName,
			LastName:  assignment.Teacher.LastName,
			Email:     assignment.Teacher.Email,
			Phone:     assignment.Teacher.Phone,
			Role:      assignment.Teacher.Role,
			Position:  assignment.Teacher.Position,
		}
	}

	return response
}
// services/student_enrollment_service.go
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
	"log"
)

type StudentEnrollmentService struct {
	db *gorm.DB
}

func NewStudentEnrollmentService(db *gorm.DB) *StudentEnrollmentService {
	return &StudentEnrollmentService{db: db}
}

// CreateStudentEnrollment creates a new student enrollment
func (s *StudentEnrollmentService) CreateStudentEnrollment(req *dto.CreateStudentEnrollmentRequest, userID uuid.UUID) (*dto.StudentEnrollmentResponse, error) {
	// Validate input
	if err := s.validateEnrollmentRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, errors.New("invalid student ID format")
	}
		log.Println("✅ Database migrated successfully")

	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	// Check if student exists
	var student models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", studentID).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("failed to verify student: " + err.Error())
	}

		// ✅ VALIDATION: Check if student role is 'student'
	if student.Role != "student" {
		return nil, errors.New("user is not a student")
	}

	// ✅ VALIDATION: Check if student is verified
	if !student.IsVerified {
		return nil, errors.New("student account is not verified. Please verify the student's email before enrollment")
	}

	// ✅ VALIDATION: Check if student is active
	if !student.IsActive {
		return nil, errors.New("student account is inactive. Please activate the student's account before enrollment")
	}



	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// ✅ VALIDATION: Check if grade is active
	if grade.Status != "active" {
		return nil, errors.New("grade is not active. Cannot enroll student in an inactive grade")
	}

	// Check if student is already enrolled in this grade
	var existing models.StudentEnrollment
	if err := s.db.Where("student_id = ? AND grade_id = ? AND deleted_at IS NULL", 
		studentID, gradeID).First(&existing).Error; err == nil {
		return nil, errors.New("student is already enrolled in this grade")
	}



	// Parse graduation date if provided
	var graduationDate *time.Time
	if req.GraduationDate != "" {
		gradDate, err := time.Parse("2006-01-02", req.GraduationDate)
		if err != nil {
			return nil, errors.New("invalid graduation date format. Use YYYY-MM-DD")
		}
		graduationDate = &gradDate
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Create new enrollment
	enrollment := &models.StudentEnrollment{
		ID:             uuid.New(),
		StudentID:      studentID,
		GradeID:        gradeID,
		Status:         status,
		GraduationDate: graduationDate,
		Notes:          strings.TrimSpace(req.Notes),
		IsVerified:  	bool(req.IsVerified),
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.Create(enrollment).Error; err != nil {
		return nil, errors.New("failed to create enrollment: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Student").Preload("Grade").First(enrollment, enrollment.ID).Error; err != nil {
		return nil, errors.New("failed to load enrollment details: " + err.Error())
	}

	return s.toEnrollmentResponse(enrollment), nil
}


// BulkCreateStudentEnrollments creates multiple student enrollments
func (s *StudentEnrollmentService) BulkCreateStudentEnrollments(req *dto.BulkCreateStudentEnrollmentRequest, userID uuid.UUID) (*dto.BulkEnrollmentResult, error) {
	// Parse Grade ID
	gradeID, err := uuid.Parse(req.GradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID format")
	}

	// Check if grade exists
	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", gradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	// ✅ VALIDATION: Check if grade is active
	if grade.Status != "active" {
		return nil, errors.New("grade is not active. Cannot enroll students in an inactive grade")
	}


	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	result := &dto.BulkEnrollmentResult{
		Created: []dto.StudentEnrollmentResponse{},
		Errors:  []dto.BulkEnrollmentError{},
	}

	for _, studentIDStr := range req.StudentIDs {
		studentID, err := uuid.Parse(studentIDStr)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "invalid student ID format",
			})
			continue
		}

		// Check if student exists
		var student models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", studentID).First(&student).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "student not found",
			})
			continue
		}

		// ✅ VALIDATION: Check if student is verified
		if !student.IsVerified {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "student account is not verified",
			})
			continue
		}

		// ✅ VALIDATION: Check if student is active
		if !student.IsActive {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "student account is inactive",
			})
			continue
		}

		// ✅ VALIDATION: Check if student role is 'student'
		if student.Role != "student" {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "user is not a student",
			})
			continue
		}

		// Check if student is already enrolled
		var existing models.StudentEnrollment
		if err := s.db.Where("student_id = ? AND grade_id = ? AND deleted_at IS NULL", 
			studentID, gradeID).First(&existing).Error; err == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "student already enrolled in this grade",
			})
			continue
		}

		// Create enrollment
		enrollment := &models.StudentEnrollment{
			ID:             uuid.New(),
			StudentID:      studentID,
			GradeID:        gradeID,
			Status:         status,
			Notes:          strings.TrimSpace(req.Notes),
			IsVerified:  	bool(req.IsVerified),
			CreatedBy:      userID,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		if err := s.db.Create(enrollment).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "failed to create enrollment: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Student").Preload("Grade").First(enrollment, enrollment.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "failed to load enrollment details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toEnrollmentResponse(enrollment))
	}

	return result, nil
}


// GetAllStudentEnrollments retrieves all student enrollments with pagination and filters
func (s *StudentEnrollmentService) GetAllStudentEnrollments(params *dto.StudentEnrollmentQueryParams) (*dto.StudentEnrollmentListResponse, error) {
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
	query := s.db.Model(&models.StudentEnrollment{}).Where("student_enrollments.deleted_at IS NULL")

	// Apply filters
	if params.StudentID != "" {
		studentID, err := uuid.Parse(params.StudentID)
		if err == nil {
			query = query.Where("student_enrollments.student_id = ?", studentID)
		}
	}

	if params.GradeID != "" {
		gradeID, err := uuid.Parse(params.GradeID)
		if err == nil {
			query = query.Where("student_enrollments.grade_id = ?", gradeID)
		}
	}

	if params.Status != "" {
		query = query.Where("student_enrollments.status = ?", params.Status)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count enrollments: %w", err)
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order("student_enrollments." + params.SortBy + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var enrollments []models.StudentEnrollment
	if err := query.Preload("Student").Preload("Grade").Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enrollments: %w", err)
	}

	// Convert to response
	responses := make([]dto.StudentEnrollmentResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = *s.toEnrollmentResponse(&enrollment)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.StudentEnrollmentListResponse{
		Enrollments: responses,
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
	}, nil
}

// GetStudentEnrollmentByID retrieves a single student enrollment by ID
func (s *StudentEnrollmentService) GetStudentEnrollmentByID(id string) (*dto.StudentEnrollmentResponse, error) {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid enrollment ID")
	}

	var enrollment models.StudentEnrollment
	if err := s.db.Where("id = ? AND deleted_at IS NULL", enrollmentID).
		Preload("Student").
		Preload("Grade").
		First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, errors.New("failed to fetch enrollment: " + err.Error())
	}

	return s.toEnrollmentResponse(&enrollment), nil
}

// GetEnrollmentsByStudent retrieves all verified enrollments for a specific student
func (s *StudentEnrollmentService) GetEnrollmentsByStudent(studentID string) ([]dto.StudentEnrollmentResponse, error) {
	sID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, errors.New("invalid student ID")
	}

	var enrollments []models.StudentEnrollment

	if err := s.db.
		Where(
			"student_id = ? AND deleted_at IS NULL AND is_verified = ?",
			sID,
			true,
		).
		Preload("Student").
		Preload("Grade").
		Order("created_at DESC").
		Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enrollments for student: %w", err)
	}

	responses := make([]dto.StudentEnrollmentResponse, len(enrollments))

	for i, enrollment := range enrollments {
		responses[i] = *s.toEnrollmentResponse(&enrollment)
	}

	return responses, nil
}

// GetEnrollmentsByGrade retrieves all enrollments for a specific grade
func (s *StudentEnrollmentService) GetEnrollmentsByGrade(gradeID string) ([]dto.StudentEnrollmentResponse, error) {
	gID, err := uuid.Parse(gradeID)
	if err != nil {
		return nil, errors.New("invalid grade ID")
	}

	var enrollments []models.StudentEnrollment
	if err := s.db.Where("grade_id = ? AND deleted_at IS NULL", gID).
		Preload("Student").
		Preload("Grade").
		Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enrollments for grade: %w", err)
	}

	responses := make([]dto.StudentEnrollmentResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = *s.toEnrollmentResponse(&enrollment)
	}

	return responses, nil
}

// GetCurrentEnrollmentByStudent retrieves the current active enrollment for a student
func (s *StudentEnrollmentService) GetCurrentEnrollmentByStudent(studentID string) (*dto.StudentEnrollmentResponse, error) {
	sID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, errors.New("invalid student ID")
	}

	var enrollment models.StudentEnrollment
	if err := s.db.Where("student_id = ? AND status = ? AND deleted_at IS NULL", sID, "active").
		Preload("Student").
		Preload("Grade").
		Order("created_at DESC").
		First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active enrollment found for this student")
		}
		return nil, errors.New("failed to fetch current enrollment: " + err.Error())
	}

	return s.toEnrollmentResponse(&enrollment), nil
}

// UpdateStudentEnrollment updates an existing student enrollment
func (s *StudentEnrollmentService) UpdateStudentEnrollment(id string, req *dto.UpdateStudentEnrollmentRequest, userID uuid.UUID) (*dto.StudentEnrollmentResponse, error) {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid enrollment ID")
	}

	// Find existing enrollment
	var enrollment models.StudentEnrollment
	if err := s.db.Where("id = ? AND deleted_at IS NULL", enrollmentID).First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, errors.New("failed to fetch enrollment: " + err.Error())
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
		enrollment.GradeID = gradeID
	}

	if req.Status != "" {
		enrollment.Status = req.Status
	}

	if req.GraduationDate != "" {
		gradDate, err := time.Parse("2006-01-02", req.GraduationDate)
		if err != nil {
			return nil, errors.New("invalid graduation date format. Use YYYY-MM-DD")
		}
		enrollment.GraduationDate = &gradDate
	}

	if req.Notes != "" {
		enrollment.Notes = strings.TrimSpace(req.Notes)
	}

	enrollment.UpdatedAt = time.Now()

	if err := s.db.Save(&enrollment).Error; err != nil {
		return nil, errors.New("failed to update enrollment: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Student").Preload("Grade").First(&enrollment, enrollment.ID).Error; err != nil {
		return nil, errors.New("failed to load enrollment details: " + err.Error())
	}

	return s.toEnrollmentResponse(&enrollment), nil
}

// DeleteStudentEnrollment soft deletes a student enrollment
func (s *StudentEnrollmentService) DeleteStudentEnrollment(id string) error {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid enrollment ID")
	}

	var enrollment models.StudentEnrollment
	if err := s.db.Where("id = ? AND deleted_at IS NULL", enrollmentID).First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("enrollment not found")
		}
		return errors.New("failed to fetch enrollment: " + err.Error())
	}

	if err := s.db.Delete(&enrollment).Error; err != nil {
		return errors.New("failed to delete enrollment: " + err.Error())
	}

	return nil
}

// validateEnrollmentRequest validates the enrollment request
func (s *StudentEnrollmentService) validateEnrollmentRequest(req *dto.CreateStudentEnrollmentRequest) error {
	if req.StudentID == "" {
		return errors.New("student ID is required")
	}
	if req.GradeID == "" {
		return errors.New("grade ID is required")
	}

	if req.Status != "" && req.Status != "active" && req.Status != "inactive" && 
		req.Status != "transferred" && req.Status != "graduated" && req.Status != "withdrawn" {
		return errors.New("status must be 'active', 'inactive', 'transferred', 'graduated', or 'withdrawn'")
	}
	return nil
}

// toEnrollmentResponse converts model to response DTO
func (s *StudentEnrollmentService) toEnrollmentResponse(enrollment *models.StudentEnrollment) *dto.StudentEnrollmentResponse {
	response := &dto.StudentEnrollmentResponse{
		ID:             enrollment.ID.String(),
		StudentID:      enrollment.StudentID.String(),
		GradeID:        enrollment.GradeID.String(),
		Status:         enrollment.Status,
		GraduationDate: enrollment.GraduationDate,
		Notes:          enrollment.Notes,
		IsVerified: 	enrollment.IsVerified,
		CreatedBy:      enrollment.CreatedBy.String(),
		CreatedAt:      enrollment.CreatedAt,
		UpdatedAt:      enrollment.UpdatedAt,
	}

	// Add student details if preloaded
	if enrollment.Student.ID != uuid.Nil {
		response.Student = &dto.UserResponse{
			ID:        enrollment.Student.ID.String(),
			FirstName: enrollment.Student.FirstName,
			LastName:  enrollment.Student.LastName,
			Email:     enrollment.Student.Email,
			Phone:     enrollment.Student.Phone,
			Role:      enrollment.Student.Role,
			IsVerified:    enrollment.Student.IsVerified,
		}
	}

	// Add grade details if preloaded
	if enrollment.Grade.ID != uuid.Nil {
		response.Grade = &dto.ClassGradeResponse{
			ID:          enrollment.Grade.ID.String(),
			Name:        enrollment.Grade.Name,
			Code:        enrollment.Grade.Code,
			Level:       enrollment.Grade.Level,
			Description: enrollment.Grade.Description,
			AcademicSessionID: enrollment.Grade.AcademicSessionID.String(),
			Capacity:    enrollment.Grade.Capacity,
			Status:      enrollment.Grade.Status,
			CreatedAt:   enrollment.Grade.CreatedAt,
			UpdatedAt:   enrollment.Grade.UpdatedAt,
		}
	}

	return response
}


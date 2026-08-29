// services/student_enrollment_service.go
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

	armID, err := uuid.Parse(req.ArmID)
	if err != nil {
		return nil, errors.New("invalid arm ID format")
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

	// Check if arm exists
	var arm models.Arm
	if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("arm not found")
		}
		return nil, errors.New("failed to verify arm: " + err.Error())
	}

	// ✅ VALIDATION: Check if arm is active
	if arm.Status != "active" {
		return nil, errors.New("arm is not active. Cannot enroll student in an inactive arm")
	}

	// ✅ VALIDATION: Check if grade exists and is active (through arm)
	if arm.GradeID == uuid.Nil {
		return nil, errors.New("arm has no grade assigned")
	}

	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", arm.GradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found for this arm")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	if grade.Status != "active" {
		return nil, errors.New("grade is not active. Cannot enroll student in an inactive grade")
	}

	// Check if student is already enrolled in this arm
	var existing models.StudentEnrollment
	if err := s.db.Where("student_id = ? AND arm_id = ? AND deleted_at IS NULL",
		studentID, armID).First(&existing).Error; err == nil {
		return nil, errors.New("student is already enrolled in this arm")
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
		ArmID:          armID,
		Status:         status,
		GraduationDate: graduationDate,
		Notes:          strings.TrimSpace(req.Notes),
		IsVerified:     req.IsVerified,
		CreatedBy:      userID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.db.Create(enrollment).Error; err != nil {
		return nil, errors.New("failed to create enrollment: " + err.Error())
	}

	// ✅ Preload relationships for response - with proper preloads
	if err := s.db.
		Preload("Student").
		Preload("Arm").
		Preload("Arm.Grade").
		Preload("Arm.ClassTeacher").
		Preload("Arm.ClassTeacher.Teacher").
		First(enrollment, enrollment.ID).Error; err != nil {
		return nil, errors.New("failed to load enrollment details: " + err.Error())
	}

	return s.toEnrollmentResponseWithGrade(enrollment), nil
}


// BulkCreateStudentEnrollments creates multiple student enrollments
func (s *StudentEnrollmentService) BulkCreateStudentEnrollments(req *dto.BulkCreateStudentEnrollmentsRequest, userID uuid.UUID) (*dto.BulkEnrollmentResult, error) {
	// Parse Arm ID
	armID, err := uuid.Parse(req.ArmID)
	if err != nil {
		return nil, errors.New("invalid arm ID format")
	}

	// Check if arm exists
	var arm models.Arm
	if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("arm not found")
		}
		return nil, errors.New("failed to verify arm: " + err.Error())
	}

	// ✅ VALIDATION: Check if arm is active
	if arm.Status != "active" {
		return nil, errors.New("arm is not active. Cannot enroll students in an inactive arm")
	}

	// ✅ VALIDATION: Check if grade exists and is active (through arm)
	if arm.GradeID == uuid.Nil {
		return nil, errors.New("arm has no grade assigned")
	}

	var grade models.ClassGrade
	if err := s.db.Where("id = ? AND deleted_at IS NULL", arm.GradeID).First(&grade).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("grade not found for this arm")
		}
		return nil, errors.New("failed to verify grade: " + err.Error())
	}

	if grade.Status != "active" {
		return nil, errors.New("grade is not active. Cannot enroll students in an inactive grade")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
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
		if err := s.db.Where("student_id = ? AND arm_id = ? AND deleted_at IS NULL",
			studentID, armID).First(&existing).Error; err == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "student already enrolled in this arm",
			})
			continue
		}

		// Create enrollment
		enrollment := &models.StudentEnrollment{
			ID:             uuid.New(),
			StudentID:      studentID,
			ArmID:          armID,
			Status:         status,
			GraduationDate: graduationDate,
			Notes:          strings.TrimSpace(req.Notes),
			IsVerified:     req.IsVerified,
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
		if err := s.db.
			Preload("Student").
			Preload("Arm").
			Preload("Arm.Grade").
			Preload("Arm.ClassTeacher").
			Preload("Arm.ClassTeacher.Teacher").
			First(enrollment, enrollment.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkEnrollmentError{
				StudentID: studentIDStr,
				Error:     "failed to load enrollment details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toEnrollmentResponseWithGrade(enrollment))
	}

	return result, nil
}

// validateEnrollmentRequest validates the enrollment request
func (s *StudentEnrollmentService) validateEnrollmentRequest(req *dto.CreateStudentEnrollmentRequest) error {
	if req.StudentID == "" {
		return errors.New("student ID is required")
	}
	if req.ArmID == "" {
		return errors.New("arm ID is required")
	}

	if req.Status != "" && req.Status != "active" && req.Status != "inactive" &&
		req.Status != "transferred" && req.Status != "graduated" && req.Status != "withdrawn" {
		return errors.New("status must be 'active', 'inactive', 'transferred', 'graduated', or 'withdrawn'")
	}
	return nil
}

// services/student_enrollment/student_enrollment_service.go

// toEnrollmentResponseWithGrade converts model to response DTO with grade info
func (s *StudentEnrollmentService) toEnrollmentResponseWithGrade(enrollment *models.StudentEnrollment) *dto.StudentEnrollmentResponse {
	response := &dto.StudentEnrollmentResponse{
		ID:             enrollment.ID.String(),
		StudentID:      enrollment.StudentID.String(),
		ArmID:          enrollment.ArmID.String(),
		Status:         enrollment.Status,
		GraduationDate: enrollment.GraduationDate,
		Notes:          enrollment.Notes,
		IsVerified:     enrollment.IsVerified,
		CreatedBy:      enrollment.CreatedBy.String(),
		CreatedAt:      enrollment.CreatedAt,
		UpdatedAt:      enrollment.UpdatedAt,
	}

	// Add student details if preloaded
	if enrollment.Student.ID != uuid.Nil {
		response.Student = &dto.UserResponse{
			ID:          enrollment.Student.ID.String(),
			FirstName:   enrollment.Student.FirstName,
			LastName:    enrollment.Student.LastName,
			MiddleName:  enrollment.Student.MiddleName,
			Email:       enrollment.Student.Email,
			Phone:       enrollment.Student.Phone,
			Role:        enrollment.Student.Role,
			Position:    enrollment.Student.Position,
			Picture:     enrollment.Student.Picture,
			IsVerified:  enrollment.Student.IsVerified,
			IsActive:    enrollment.Student.IsActive,
			Location:    enrollment.Student.Location,
			LastLoginAt: enrollment.Student.LastLoginAt,
			CreatedAt:   enrollment.Student.CreatedAt,
			UpdatedAt:   enrollment.Student.UpdatedAt,
		}
	}

	// Add arm details with grade if preloaded
	if enrollment.Arm.ID != uuid.Nil {
		armResponse := &dto.ArmResponse{
			ID:          enrollment.Arm.ID.String(),
			Name:        enrollment.Arm.Name,
			Code:        enrollment.Arm.Code,
			Description: enrollment.Arm.Description,
			GradeID:     enrollment.Arm.GradeID.String(),
			Capacity:    enrollment.Arm.Capacity,
			Status:      enrollment.Arm.Status,
			CreatedBy:   enrollment.Arm.CreatedBy.String(),
			CreatedAt:   enrollment.Arm.CreatedAt,
			UpdatedAt:   enrollment.Arm.UpdatedAt,
		}

		// ✅ Check if Grade exists before accessing
		if enrollment.Arm.Grade.ID != uuid.Nil {
			armResponse.Grade = &dto.ClassGradeResponse{
				ID:          enrollment.Arm.Grade.ID.String(),
				Name:        enrollment.Arm.Grade.Name,
				Code:        enrollment.Arm.Grade.Code,
				Level:       enrollment.Arm.Grade.Level,
				Description: enrollment.Arm.Grade.Description,
				Status:      enrollment.Arm.Grade.Status,
				CreatedAt:   enrollment.Arm.Grade.CreatedAt,
				UpdatedAt:   enrollment.Arm.Grade.UpdatedAt,
			}
		}

		// ✅ Check if ClassTeacher exists before accessing
		if enrollment.Arm.ClassTeacher.ID != uuid.Nil {
			if enrollment.Arm.ClassTeacher.Teacher.ID != uuid.Nil {
				armResponse.ClassTeacher = &dto.UserResponse{
					ID:          enrollment.Arm.ClassTeacher.Teacher.ID.String(),
					FirstName:   enrollment.Arm.ClassTeacher.Teacher.FirstName,
					LastName:    enrollment.Arm.ClassTeacher.Teacher.LastName,
					MiddleName:  enrollment.Arm.ClassTeacher.Teacher.MiddleName,
					Email:       enrollment.Arm.ClassTeacher.Teacher.Email,
					Phone:       enrollment.Arm.ClassTeacher.Teacher.Phone,
					Role:        enrollment.Arm.ClassTeacher.Teacher.Role,
					Position:    enrollment.Arm.ClassTeacher.Teacher.Position,
					Picture:     enrollment.Arm.ClassTeacher.Teacher.Picture,
					IsVerified:  enrollment.Arm.ClassTeacher.Teacher.IsVerified,
					IsActive:    enrollment.Arm.ClassTeacher.Teacher.IsActive,
					Location:    enrollment.Arm.ClassTeacher.Teacher.Location,
					LastLoginAt: enrollment.Arm.ClassTeacher.Teacher.LastLoginAt,
					CreatedAt:   enrollment.Arm.ClassTeacher.Teacher.CreatedAt,
					UpdatedAt:   enrollment.Arm.ClassTeacher.Teacher.UpdatedAt,
				}
			} else {
				// ClassTeacher exists but Teacher relationship is not loaded
				armResponse.ClassTeacher = &dto.UserResponse{
					ID: enrollment.Arm.ClassTeacher.TeacherID.String(),
				}
			}
		}

		response.Arm = armResponse
	}

	return response
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

	if params.ArmID != "" {
		armID, err := uuid.Parse(params.ArmID)
		if err == nil {
			query = query.Where("student_enrollments.arm_id = ?", armID)
		}
	}

	if params.Status != "" {
		query = query.Where("student_enrollments.status = ?", params.Status)
	}

	if params.IsVerified != nil {
		query = query.Where("student_enrollments.is_verified = ?", *params.IsVerified)
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

	// Execute with preloads - Preload Arm and its Grade
	var enrollments []models.StudentEnrollment
	if err := query.
		Preload("Student").
		Preload("Arm").
		Preload("Arm.Grade").
		Preload("Arm.ClassTeacher").
		Preload("Arm.ClassTeacher.Teacher").
		Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enrollments: %w", err)
	}

	// Convert to response
	responses := make([]dto.StudentEnrollmentResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = *s.toEnrollmentResponseWithGrade(&enrollment)
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

// GetEnrollmentsByStudent retrieves all verified enrollments for a specific student with grade info
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
		Preload("Arm").
		Preload("Arm.Grade").
		Preload("Arm.ClassTeacher").
		Preload("Arm.ClassTeacher.Teacher").
		Order("created_at DESC").
		Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enrollments for student: %w", err)
	}

	responses := make([]dto.StudentEnrollmentResponse, len(enrollments))

	for i, enrollment := range enrollments {
		responses[i] = *s.toEnrollmentResponseWithGrade(&enrollment)
	}

	return responses, nil
}

// GetEnrollmentsByArm retrieves all enrollments for a specific arm with grade info
func (s *StudentEnrollmentService) GetEnrollmentsByArm(armID string) ([]dto.StudentEnrollmentResponse, error) {
	aID, err := uuid.Parse(armID)
	if err != nil {
		return nil, errors.New("invalid arm ID")
	}

	var enrollments []models.StudentEnrollment
	if err := s.db.Where("arm_id = ? AND deleted_at IS NULL", aID).
		Preload("Student").
		Preload("Arm").
		Preload("Arm.Grade").
		Preload("Arm.ClassTeacher").
		Preload("Arm.ClassTeacher.Teacher").
		Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enrollments for arm: %w", err)
	}

	responses := make([]dto.StudentEnrollmentResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = *s.toEnrollmentResponseWithGrade(&enrollment)
	}

	return responses, nil
}

// GetStudentEnrollmentByID retrieves a single student enrollment by ID with grade info
func (s *StudentEnrollmentService) GetStudentEnrollmentByID(id string) (*dto.StudentEnrollmentResponse, error) {
	enrollmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid enrollment ID")
	}

	var enrollment models.StudentEnrollment
	if err := s.db.Where("id = ? AND deleted_at IS NULL", enrollmentID).
		Preload("Student").
		Preload("Arm").
		Preload("Arm.Grade").
		Preload("Arm.ClassTeacher").
		Preload("Arm.ClassTeacher.Teacher").
		First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("enrollment not found")
		}
		return nil, errors.New("failed to fetch enrollment: " + err.Error())
	}

	return s.toEnrollmentResponseWithGrade(&enrollment), nil
}

// GetCurrentEnrollmentByStudent retrieves the current active enrollment for a student with grade info
func (s *StudentEnrollmentService) GetCurrentEnrollmentByStudent(studentID string) (*dto.StudentEnrollmentResponse, error) {
	sID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, errors.New("invalid student ID")
	}

	var enrollment models.StudentEnrollment
	if err := s.db.Where("student_id = ? AND status = ? AND deleted_at IS NULL", sID, "active").
		Preload("Student").
		Preload("Arm").
		Preload("Arm.Grade").
		Preload("Arm.ClassTeacher").
		Preload("Arm.ClassTeacher.Teacher").
		Order("created_at DESC").
		First(&enrollment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no active enrollment found for this student")
		}
		return nil, errors.New("failed to fetch current enrollment: " + err.Error())
	}

	return s.toEnrollmentResponseWithGrade(&enrollment), nil
}

// toEnrollmentResponse converts model to response DTO
func (s *StudentEnrollmentService) toEnrollmentResponse(enrollment *models.StudentEnrollment) *dto.StudentEnrollmentResponse {
	response := &dto.StudentEnrollmentResponse{
		ID:             enrollment.ID.String(),
		StudentID:      enrollment.StudentID.String(),
		ArmID:          enrollment.ArmID.String(),
		Status:         enrollment.Status,
		GraduationDate: enrollment.GraduationDate,
		Notes:          enrollment.Notes,
		IsVerified:     enrollment.IsVerified,
		CreatedBy:      enrollment.CreatedBy.String(),
		CreatedAt:      enrollment.CreatedAt,
		UpdatedAt:      enrollment.UpdatedAt,
	}

	// Add student details if preloaded
	if enrollment.Student.ID != uuid.Nil {
		response.Student = &dto.UserResponse{
			ID:          enrollment.Student.ID.String(),
			FirstName:   enrollment.Student.FirstName,
			LastName:    enrollment.Student.LastName,
			Email:       enrollment.Student.Email,
			Phone:       enrollment.Student.Phone,
			Role:        enrollment.Student.Role,
			IsVerified:  enrollment.Student.IsVerified,
			IsActive:    enrollment.Student.IsActive,
			Picture:     enrollment.Student.Picture,
			Location:    enrollment.Student.Location,
			MiddleName:  enrollment.Student.MiddleName,
			LastLoginAt: enrollment.Student.LastLoginAt,
			CreatedAt:   enrollment.Student.CreatedAt,
			UpdatedAt:   enrollment.Student.UpdatedAt,
		}
	}

	// Add arm details with grade if preloaded
	if enrollment.Arm.ID != uuid.Nil {
		armResponse := &dto.ArmResponse{
			ID:          enrollment.Arm.ID.String(),
			Name:        enrollment.Arm.Name,
			Code:        enrollment.Arm.Code,
			Description: enrollment.Arm.Description,
			GradeID:     enrollment.Arm.GradeID.String(),
			Capacity:    enrollment.Arm.Capacity,
			Status:      enrollment.Arm.Status,
			CreatedBy:   enrollment.Arm.CreatedBy.String(),
			CreatedAt:   enrollment.Arm.CreatedAt,
			UpdatedAt:   enrollment.Arm.UpdatedAt,
		}

		// Add grade details if preloaded
		if enrollment.Arm.Grade.ID != uuid.Nil {
			armResponse.Grade = &dto.ClassGradeResponse{
				ID:          enrollment.Arm.Grade.ID.String(),
				Name:        enrollment.Arm.Grade.Name,
				Code:        enrollment.Arm.Grade.Code,
				Level:       enrollment.Arm.Grade.Level,
				Description: enrollment.Arm.Grade.Description,
				Status:      enrollment.Arm.Grade.Status,
				CreatedAt:   enrollment.Arm.Grade.CreatedAt,
				UpdatedAt:   enrollment.Arm.Grade.UpdatedAt,
			}
		}

		// Add class teacher details if preloaded
		if enrollment.Arm.ClassTeacher.ID != uuid.Nil && enrollment.Arm.ClassTeacher.Teacher.ID != uuid.Nil {
			armResponse.ClassTeacher = &dto.UserResponse{
				ID:         enrollment.Arm.ClassTeacher.Teacher.ID.String(),
				FirstName:  enrollment.Arm.ClassTeacher.Teacher.FirstName,
				LastName:   enrollment.Arm.ClassTeacher.Teacher.LastName,
				Email:      enrollment.Arm.ClassTeacher.Teacher.Email,
				Phone:      enrollment.Arm.ClassTeacher.Teacher.Phone,
				Role:       enrollment.Arm.ClassTeacher.Teacher.Role,
				IsVerified: enrollment.Arm.ClassTeacher.Teacher.IsVerified,
				IsActive:   enrollment.Arm.ClassTeacher.Teacher.IsActive,
				Picture:    enrollment.Arm.ClassTeacher.Teacher.Picture,
			}
		}

		response.Arm = armResponse
	}

	return response
}

// GetEnrollmentsByArmAndStatus retrieves enrollments for an arm filtered by status
func (s *StudentEnrollmentService) GetEnrollmentsByArmAndStatus(armID, status string) ([]dto.StudentEnrollmentResponse, error) {
	aID, err := uuid.Parse(armID)
	if err != nil {
		return nil, errors.New("invalid arm ID")
	}

	var enrollments []models.StudentEnrollment
	if err := s.db.Where("arm_id = ? AND status = ? AND deleted_at IS NULL", aID, status).
		Preload("Student").
		Preload("Arm").
		Find(&enrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enrollments for arm: %w", err)
	}

	responses := make([]dto.StudentEnrollmentResponse, len(enrollments))
	for i, enrollment := range enrollments {
		responses[i] = *s.toEnrollmentResponse(&enrollment)
	}

	return responses, nil
}

// GetEnrollmentStats retrieves statistics for student enrollments
func (s *StudentEnrollmentService) GetEnrollmentStats(filter map[string]interface{}) (*dto.StudentEnrollmentStats, error) {
	query := s.db.Model(&models.StudentEnrollment{}).Where("deleted_at IS NULL")

	// Apply filters
	if armID, ok := filter["arm_id"].(string); ok && armID != "" {
		if id, err := uuid.Parse(armID); err == nil {
			query = query.Where("arm_id = ?", id)
		}
	}
	if studentID, ok := filter["student_id"].(string); ok && studentID != "" {
		if id, err := uuid.Parse(studentID); err == nil {
			query = query.Where("student_id = ?", id)
		}
	}

	var stats dto.StudentEnrollmentStats

	// Count total
	if err := query.Count(&stats.TotalEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count total enrollments: %w", err)
	}

	// Count by status
	if err := query.Where("status = ?", "active").Count(&stats.ActiveEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count active enrollments: %w", err)
	}
	if err := query.Where("status = ?", "inactive").Count(&stats.InactiveEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count inactive enrollments: %w", err)
	}
	if err := query.Where("status = ?", "graduated").Count(&stats.GraduatedEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count graduated enrollments: %w", err)
	}
	if err := query.Where("status = ?", "transferred").Count(&stats.TransferredEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count transferred enrollments: %w", err)
	}
	if err := query.Where("status = ?", "withdrawn").Count(&stats.WithdrawnEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count withdrawn enrollments: %w", err)
	}

	// Count verified/unverified
	if err := query.Where("is_verified = ?", true).Count(&stats.VerifiedEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count verified enrollments: %w", err)
	}
	if err := query.Where("is_verified = ?", false).Count(&stats.UnverifiedEnrollments).Error; err != nil {
		return nil, fmt.Errorf("failed to count unverified enrollments: %w", err)
	}

	return &stats, nil
}

// services/student_enrollment/student_enrollment_service.go

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
	if req.ArmID != "" {
		armID, err := uuid.Parse(req.ArmID)
		if err != nil {
			return nil, errors.New("invalid arm ID format")
		}
		// Verify arm exists
		var arm models.Arm
		if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("arm not found")
			}
			return nil, errors.New("failed to verify arm: " + err.Error())
		}
		enrollment.ArmID = armID
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

	if req.IsVerified != nil {
		enrollment.IsVerified = *req.IsVerified
	}

	enrollment.UpdatedAt = time.Now()

	if err := s.db.Save(&enrollment).Error; err != nil {
		return nil, errors.New("failed to update enrollment: " + err.Error())
	}

	// ✅ FIX: Preload all relationships including Arm.Grade and ClassTeacher
	if err := s.db.
		Preload("Student").
		Preload("Arm").
		Preload("Arm.Grade").
		Preload("Arm.ClassTeacher").
		Preload("Arm.ClassTeacher.Teacher").
		First(&enrollment, enrollment.ID).Error; err != nil {
		return nil, errors.New("failed to load enrollment details: " + err.Error())
	}

	// ✅ Use the safe function that handles nil Grade and ClassTeacher
	return s.toEnrollmentResponseWithGrade(&enrollment), nil
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

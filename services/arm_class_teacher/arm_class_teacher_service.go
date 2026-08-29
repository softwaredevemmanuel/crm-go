// services/arm_class_teacher_service.go
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

type ArmClassTeacherService struct {
	db *gorm.DB
}

func NewArmClassTeacherService(db *gorm.DB) *ArmClassTeacherService {
	return &ArmClassTeacherService{db: db}
}

// CreateAssignment creates a new arm class teacher assignment
func (s *ArmClassTeacherService) CreateAssignment(req *dto.CreateArmClassTeacherRequest, userID uuid.UUID) (*dto.ArmClassTeacherResponse, error) {
	// Validate input
	if err := s.validateAssignmentRequest(req); err != nil {
		return nil, err
	}

	// Parse UUIDs
	armID, err := uuid.Parse(req.ArmID)
	if err != nil {
		return nil, errors.New("invalid arm ID format")
	}

	teacherID, err := uuid.Parse(req.TeacherID)
	if err != nil {
		return nil, errors.New("invalid teacher ID format")
	}


	// Check if arm exists
	var arm models.Arm
	if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("arm not found")
		}
		return nil, errors.New("failed to verify arm: " + err.Error())
	}

	// Check if teacher exists and is a teacher
	var teacher models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", teacherID).First(&teacher).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("teacher not found")
		}
		return nil, errors.New("failed to verify teacher: " + err.Error())
	}

	if teacher.Role != "teacher" && teacher.Role != "admin" && teacher.Role != "staff" {
		return nil, errors.New("user is not a teacher or admin")
	}


	// Check if assignment already exists
	var existing models.ArmClassTeacher
	if err := s.db.Where("arm_id = ? AND teacher_id = ? AND deleted_at IS NULL",
		armID, teacherID).First(&existing).Error; err == nil {
		return nil, errors.New("this teacher is already assigned to this arm")
	}

	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	// Check if arm already has a class teacher
	var existingArmTeacher models.ArmClassTeacher
	if err := s.db.Where("arm_id = ? AND status = ? AND deleted_at IS NULL",
		armID, "active").First(&existingArmTeacher).Error; err == nil {
		return nil, errors.New("this arm already has an active class teacher")
	}

	// Create assignment
	assignment := &models.ArmClassTeacher{
		ID:        uuid.New(),
		ArmID:     armID,
		TeacherID: teacherID,
		Status:    status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.db.Create(assignment).Error; err != nil {
		return nil, errors.New("failed to create assignment: " + err.Error())
	}

	// Preload relationships for response
	if err := s.db.Preload("Arm").Preload("Teacher").First(assignment, assignment.ID).Error; err != nil {
		return nil, errors.New("failed to load assignment details: " + err.Error())
	}

	return s.toAssignmentResponse(assignment), nil
}

// BulkAssignClassTeachers assigns multiple class teachers to arms
func (s *ArmClassTeacherService) BulkAssignClassTeachers(req *dto.BulkAssignClassTeachersRequest, userID uuid.UUID) (*dto.BulkAssignResult, error) {
	// Set default status
	status := req.Status
	if status == "" {
		status = "active"
	}

	result := &dto.BulkAssignResult{
		Created: []dto.ArmClassTeacherResponse{},
		Errors:  []dto.BulkAssignError{},
	}

	for _, assignReq := range req.Assignments {
		// Parse UUIDs
		armID, err := uuid.Parse(assignReq.ArmID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignError{
				ArmID: assignReq.ArmID,
				Error: "invalid arm ID format",
			})
			continue
		}

		teacherID, err := uuid.Parse(assignReq.TeacherID)
		if err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignError{
				TeacherID: assignReq.TeacherID,
				Error:     "invalid teacher ID format",
			})
			continue
		}



		// Check if arm exists
		var arm models.Arm
		if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignError{
				ArmID: assignReq.ArmID,
				Error: "arm not found",
			})
			continue
		}

		// Check if teacher exists
		var teacher models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", teacherID).First(&teacher).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignError{
				TeacherID: assignReq.TeacherID,
				Error:     "teacher not found",
			})
			continue
		}



		// Check if assignment already exists
		var existing models.ArmClassTeacher
		if err := s.db.Where("arm_id = ? AND teacher_id = ? AND deleted_at IS NULL",
			armID, teacherID).First(&existing).Error; err == nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignError{
				ArmID:   assignReq.ArmID,
				TeacherID: assignReq.TeacherID,
				Error:   "assignment already exists",
			})
			continue
		}

		// Create assignment
		assignment := &models.ArmClassTeacher{
			ID:        uuid.New(),
			ArmID:     armID,
			TeacherID: teacherID,
			Status:    status,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.db.Create(assignment).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignError{
				ArmID:   assignReq.ArmID,
				TeacherID: assignReq.TeacherID,
				Error:   "failed to create assignment: " + err.Error(),
			})
			continue
		}

		// Preload relationships
		if err := s.db.Preload("Arm").Preload("Teacher").First(assignment, assignment.ID).Error; err != nil {
			result.FailedCount++
			result.Errors = append(result.Errors, dto.BulkAssignError{
				ArmID:   assignReq.ArmID,
				TeacherID: assignReq.TeacherID,
				Error:   "failed to load assignment details",
			})
			continue
		}

		result.SuccessCount++
		result.Created = append(result.Created, *s.toAssignmentResponse(assignment))
	}

	return result, nil
}

// GetAllAssignments retrieves all assignments with pagination and filters
func (s *ArmClassTeacherService) GetAllAssignments(params *dto.ArmClassTeacherQueryParams) (*dto.ArmClassTeacherListResponse, error) {
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
	query := s.db.Model(&models.ArmClassTeacher{}).Where("deleted_at IS NULL")

	// Apply filters
	if params.ArmID != "" {
		armID, err := uuid.Parse(params.ArmID)
		if err == nil {
			query = query.Where("arm_id = ?", armID)
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
	var assignments []models.ArmClassTeacher
	if err := query.Preload("Arm").Preload("Teacher").Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch assignments: %w", err)
	}

	// Convert to response
	responses := make([]dto.ArmClassTeacherResponse, len(assignments))
	for i, assignment := range assignments {
		responses[i] = *s.toAssignmentResponse(&assignment)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.ArmClassTeacherListResponse{
		Assignments: responses,
		Total:       total,
		Page:        params.Page,
		Limit:       params.Limit,
		TotalPages:  totalPages,
	}, nil
}

// GetAssignmentByID retrieves a single assignment by ID
func (s *ArmClassTeacherService) GetAssignmentByID(id string) (*dto.ArmClassTeacherResponse, error) {
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid assignment ID")
	}

	var assignment models.ArmClassTeacher
	if err := s.db.Where("id = ? AND deleted_at IS NULL", assignmentID).
		Preload("Arm").
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
func (s *ArmClassTeacherService) GetAssignmentsByTeacher(teacherID string) ([]dto.ArmClassTeacherResponse, error) {
	tID, err := uuid.Parse(teacherID)
	if err != nil {
		return nil, errors.New("invalid teacher ID")
	}

	var assignments []models.ArmClassTeacher
	if err := s.db.Where("teacher_id = ? AND deleted_at IS NULL", tID).
		Preload("Arm").
		Preload("Teacher").
		Order("created_at DESC").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch assignments for teacher: %w", err)
	}

	responses := make([]dto.ArmClassTeacherResponse, len(assignments))
	for i, assignment := range assignments {
		responses[i] = *s.toAssignmentResponse(&assignment)
	}

	return responses, nil
}

// GetAssignmentsByArm retrieves all assignments for a specific arm
func (s *ArmClassTeacherService) GetAssignmentsByArm(armID string) ([]dto.ArmClassTeacherResponse, error) {
	aID, err := uuid.Parse(armID)
	if err != nil {
		return nil, errors.New("invalid arm ID")
	}

	var assignments []models.ArmClassTeacher
	if err := s.db.Where("arm_id = ? AND deleted_at IS NULL", aID).
		Preload("Arm").
		Preload("Teacher").
		Order("created_at DESC").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch assignments for arm: %w", err)
	}

	responses := make([]dto.ArmClassTeacherResponse, len(assignments))
	for i, assignment := range assignments {
		responses[i] = *s.toAssignmentResponse(&assignment)
	}

	return responses, nil
}


// GetArmsWithClassTeachers retrieves all arms with their class teachers
func (s *ArmClassTeacherService) GetArmsWithClassTeachers() ([]dto.ArmWithClassTeacher, error) {
	var arms []models.Arm
	if err := s.db.Where("deleted_at IS NULL").Find(&arms).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch arms: %w", err)
	}

	var assignments []models.ArmClassTeacher
	if err := s.db.Where("status = ? AND deleted_at IS NULL", "active").
		Preload("Arm").
		Preload("Teacher").
		Find(&assignments).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch assignments: %w", err)
	}

	// Create a map of arm_id -> assignment
	assignmentMap := make(map[string]models.ArmClassTeacher)
	for _, assignment := range assignments {
		assignmentMap[assignment.ArmID.String()] = assignment
	}

	result := make([]dto.ArmWithClassTeacher, len(arms))
	for i, arm := range arms {
		result[i] = dto.ArmWithClassTeacher{
			ArmID:   arm.ID.String(),
			ArmName: arm.Name,
			ArmCode: arm.Code,
		}

		if assignment, exists := assignmentMap[arm.ID.String()]; exists {
			result[i].TeacherID = assignment.TeacherID.String()
			result[i].TeacherName = assignment.Teacher.FirstName + " " + assignment.Teacher.LastName
			result[i].TeacherEmail = assignment.Teacher.Email
			result[i].Status = "assigned"
		} else {
			result[i].Status = "unassigned"
		}
	}

	return result, nil
}

// UpdateAssignment updates an existing assignment
func (s *ArmClassTeacherService) UpdateAssignment(id string, req *dto.UpdateArmClassTeacherRequest) (*dto.ArmClassTeacherResponse, error) {
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid assignment ID")
	}

	// Find existing assignment
	var assignment models.ArmClassTeacher
	if err := s.db.Where("id = ? AND deleted_at IS NULL", assignmentID).First(&assignment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("assignment not found")
		}
		return nil, errors.New("failed to fetch assignment: " + err.Error())
	}

	// Update fields
	if req.ArmID != "" {
		armID, err := uuid.Parse(req.ArmID)
		if err != nil {
			return nil, errors.New("invalid arm ID format")
		}
		var arm models.Arm
		if err := s.db.Where("id = ? AND deleted_at IS NULL", armID).First(&arm).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("arm not found")
			}
			return nil, errors.New("failed to verify arm: " + err.Error())
		}
		assignment.ArmID = armID
	}

	if req.TeacherID != "" {
		teacherID, err := uuid.Parse(req.TeacherID)
		if err != nil {
			return nil, errors.New("invalid teacher ID format")
		}
		var teacher models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", teacherID).First(&teacher).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("teacher not found")
			}
			return nil, errors.New("failed to verify teacher: " + err.Error())
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
	if err := s.db.Preload("Arm").Preload("Teacher").First(&assignment, assignment.ID).Error; err != nil {
		return nil, errors.New("failed to load assignment details: " + err.Error())
	}

	return s.toAssignmentResponse(&assignment), nil
}

// DeleteAssignment soft deletes an assignment
func (s *ArmClassTeacherService) DeleteAssignment(id string) error {
	assignmentID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid assignment ID")
	}

	var assignment models.ArmClassTeacher
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
func (s *ArmClassTeacherService) validateAssignmentRequest(req *dto.CreateArmClassTeacherRequest) error {
	if req.ArmID == "" {
		return errors.New("arm ID is required")
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
func (s *ArmClassTeacherService) toAssignmentResponse(assignment *models.ArmClassTeacher) *dto.ArmClassTeacherResponse {
	response := &dto.ArmClassTeacherResponse{
		ID:        assignment.ID.String(),
		ArmID:     assignment.ArmID.String(),
		TeacherID: assignment.TeacherID.String(),
		Status:    assignment.Status,
		CreatedAt: assignment.CreatedAt,
		UpdatedAt: assignment.UpdatedAt,
	}

	// Add arm details if preloaded
	if assignment.Arm.ID != uuid.Nil {
		response.Arm = &dto.ArmResponse{
			ID:      assignment.Arm.ID.String(),
			Name:    assignment.Arm.Name,
			Code:    assignment.Arm.Code,
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
		}
	}





	return response
}
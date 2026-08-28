// services/guardian_service.go
package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"crm-go/models"
	"crm-go/dto"
)

type GuardianService struct {
	db *gorm.DB
}

func NewGuardianService(db *gorm.DB) *GuardianService {
	return &GuardianService{db: db}
}

// CreateGuardian creates a new guardian and a user account for them
func (s *GuardianService) CreateGuardian(req *dto.CreateGuardianRequest, userID uuid.UUID) (*dto.GuardianResponse, error) {
	// Validate input
	if err := s.validateGuardianRequest(req); err != nil {
		return nil, err
	}

	// Parse Student ID
	studentID, err := uuid.Parse(req.StudentID)
	if err != nil {
		return nil, errors.New("invalid student ID")
	}

	// Check if student exists
	var student models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", studentID).First(&student).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, errors.New("failed to verify student: " + err.Error())
	}

	// Check if user with same email already exists
	var existingUser models.User
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("a user with this email already exists")
	}

	// Set default status if not provided
	status := req.Status
	if status == "" {
		status = "active"
	}

	// ========================================================
	// 1. CREATE USER ACCOUNT FOR GUARDIAN
	// ========================================================
	
	defaultPassword := "12345"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(defaultPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password: " + err.Error())
	}

// Generate login_id (can be email or custom)
	emailTrim := strings.ToLower(strings.TrimSpace(req.Email))

	parts := strings.SplitN(emailTrim, "@", 2)

	loginID := parts[0]

	guardianUser := &models.User{
		ID:         uuid.New(),
		FirstName:  strings.TrimSpace(req.FirstName),
		LastName:   strings.TrimSpace(req.LastName),
		MiddleName: strings.TrimSpace(req.MiddleName),
		Email:      strings.ToLower(strings.TrimSpace(req.Email)),
		Password:   string(hashedPassword),
		LoginID:    loginID,
		Role:       "parent",
		Position:   "parent",
		Phone:      strings.TrimSpace(req.Phone),
		Picture:    "",
		Provider:   "local",
		IsVerified: true,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.db.Create(guardianUser).Error; err != nil {
		return nil, errors.New("failed to create user account for guardian: " + err.Error())
	}

	// ========================================================
	// 2. CREATE GUARDIAN RECORD
	// ========================================================

	guardian := &models.Guardian{
		ID:           uuid.New(),
		Occupation:   strings.TrimSpace(req.Occupation),
		Relationship: req.Relationship,
		Address:      strings.TrimSpace(req.Address),
		StudentID:    studentID,
		UserID:       guardianUser.ID,
		Status:       status,
		IsPrimary:    req.IsPrimary,
		IsEmergency:  req.IsEmergency,
		CreatedBy:    userID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.db.Create(guardian).Error; err != nil {
		return nil, errors.New("failed to create guardian: " + err.Error())
	}

	// Preload relationships
	if err := s.db.Preload("Student").Preload("User").First(guardian, guardian.ID).Error; err != nil {
		return nil, errors.New("failed to load guardian details: " + err.Error())
	}

	return s.toGuardianResponse(guardian), nil
}

// GetAllGuardians retrieves all guardians with pagination and filters
func (s *GuardianService) GetAllGuardians(params *dto.GuardianQueryParams) (*dto.GuardianListResponse, error) {
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
	query := s.db.Model(&models.Guardian{}).Where("guardians.deleted_at IS NULL")

	// Apply filters - search by user fields
	if params.Search != "" {
		searchTerm := "%" + strings.ToLower(params.Search) + "%"
		query = query.Joins("JOIN users ON users.id = guardians.user_id").
			Where("LOWER(users.first_name) LIKE ? OR LOWER(users.last_name) LIKE ? OR LOWER(users.email) LIKE ? OR LOWER(users.phone) LIKE ?",
				searchTerm, searchTerm, searchTerm, searchTerm)
	}

	if params.StudentID != "" {
		studentID, err := uuid.Parse(params.StudentID)
		if err == nil {
			query = query.Where("guardians.student_id = ?", studentID)
		}
	}

	if params.Relationship != "" {
		query = query.Where("guardians.relationship = ?", params.Relationship)
	}

	if params.Status != "" {
		query = query.Where("guardians.status = ?", params.Status)
	}

	if params.IsPrimary != nil {
		query = query.Where("guardians.is_primary = ?", *params.IsPrimary)
	}

	if params.IsEmergency != nil {
		query = query.Where("guardians.is_emergency = ?", *params.IsEmergency)
	}

	// Get total count
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("failed to count guardians: %w", err)
	}

	// Apply sorting
	sortDirection := "DESC"
	if strings.ToLower(params.SortOrder) == "asc" {
		sortDirection = "ASC"
	}
	query = query.Order("guardians." + params.SortBy + " " + sortDirection)

	// Apply pagination
	offset := (params.Page - 1) * params.Limit
	query = query.Offset(offset).Limit(params.Limit)

	// Execute with preloads
	var guardians []models.Guardian
	if err := query.Preload("Student").Preload("User").Find(&guardians).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch guardians: %w", err)
	}

	// Convert to response
	responses := make([]dto.GuardianResponse, len(guardians))
	for i, guardian := range guardians {
		responses[i] = *s.toGuardianResponse(&guardian)
	}

	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return &dto.GuardianListResponse{
		Guardians:  responses,
		Total:      total,
		Page:       params.Page,
		Limit:      params.Limit,
		TotalPages: totalPages,
	}, nil
}

// GetGuardianByID retrieves a single guardian by ID
func (s *GuardianService) GetGuardianByID(id string) (*dto.GuardianResponse, error) {
	guardianID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid guardian ID")
	}

	var guardian models.Guardian
	if err := s.db.Where("id = ? AND deleted_at IS NULL", guardianID).
		Preload("Student").
		Preload("User").
		First(&guardian).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("guardian not found")
		}
		return nil, errors.New("failed to fetch guardian: " + err.Error())
	}

	return s.toGuardianResponse(&guardian), nil
}

// GetGuardiansByStudent retrieves all guardians for a specific student
func (s *GuardianService) GetGuardiansByStudent(studentID string) ([]dto.GuardianResponse, error) {
	sID, err := uuid.Parse(studentID)
	if err != nil {
		return nil, errors.New("invalid student ID")
	}

	var guardians []models.Guardian
	if err := s.db.Where("student_id = ? AND deleted_at IS NULL", sID).
		Preload("Student").
		Preload("User").
		Order("is_primary DESC, created_at ASC").
		Find(&guardians).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch guardians for student: %w", err)
	}

	responses := make([]dto.GuardianResponse, len(guardians))
	for i, guardian := range guardians {
		responses[i] = *s.toGuardianResponse(&guardian)
	}

	return responses, nil
}


// UpdateGuardian updates an existing guardian and the associated user account
func (s *GuardianService) UpdateGuardian(id string, req *dto.UpdateGuardianRequest, userID uuid.UUID) (*dto.GuardianResponse, error) {
	guardianID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid guardian ID")
	}

	// Find existing guardian
	var guardian models.Guardian
	if err := s.db.Where("id = ? AND deleted_at IS NULL", guardianID).First(&guardian).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("guardian not found")
		}
		return nil, errors.New("failed to fetch guardian: " + err.Error())
	}

	// Find the associated user
	var user models.User
	if err := s.db.Where("id = ? AND deleted_at IS NULL", guardian.UserID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("associated user account not found")
		}
		return nil, errors.New("failed to fetch user account: " + err.Error())
	}

	// ========================================================
	// UPDATE GUARDIAN FIELDS (only guardian-specific fields)
	// ========================================================
	
	if req.Occupation != "" {
		guardian.Occupation = strings.TrimSpace(req.Occupation)
	}
	if req.Relationship != "" {
		guardian.Relationship = req.Relationship
	}
	if req.Address != "" {
		guardian.Address = strings.TrimSpace(req.Address)
	}
	if req.StudentID != "" {
		sID, err := uuid.Parse(req.StudentID)
		if err != nil {
			return nil, errors.New("invalid student ID")
		}
		var student models.User
		if err := s.db.Where("id = ? AND deleted_at IS NULL", sID).First(&student).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("student not found")
			}
			return nil, errors.New("failed to verify student: " + err.Error())
		}
		guardian.StudentID = sID
	}
	if req.IsPrimary != nil {
		guardian.IsPrimary = *req.IsPrimary
	}
	if req.IsEmergency != nil {
		guardian.IsEmergency = *req.IsEmergency
	}
	if req.Status != "" {
		guardian.Status = req.Status
	}

	// Update guardian timestamp
	guardian.UpdatedAt = time.Now()

	

	// ========================================================
	// SAVE BOTH RECORDS
	// ========================================================

	if err := s.db.Save(&guardian).Error; err != nil {
		return nil, errors.New("failed to update guardian: " + err.Error())
	}


	// Preload relationships
	if err := s.db.Preload("Student").Preload("User").First(&guardian, guardian.ID).Error; err != nil {
		return nil, errors.New("failed to load guardian details: " + err.Error())
	}

	return s.toGuardianResponse(&guardian), nil
}

// DeleteGuardian soft deletes a guardian
func (s *GuardianService) DeleteGuardian(id string) error {
	guardianID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("invalid guardian ID")
	}

	var guardian models.Guardian
	if err := s.db.Where("id = ? AND deleted_at IS NULL", guardianID).First(&guardian).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("guardian not found")
		}
		return errors.New("failed to fetch guardian: " + err.Error())
	}

	if err := s.db.Delete(&guardian).Error; err != nil {
		return errors.New("failed to delete guardian: " + err.Error())
	}

	return nil
}

// validateGuardianRequest validates the guardian request
func (s *GuardianService) validateGuardianRequest(req *dto.CreateGuardianRequest) error {
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	if len(req.FirstName) < 2 {
		return errors.New("first name must be at least 2 characters")
	}
	if req.LastName == "" {
		return errors.New("last name is required")
	}
	if len(req.LastName) < 2 {
		return errors.New("last name must be at least 2 characters")
	}
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Phone == "" {
		return errors.New("phone number is required")
	}
	if req.Relationship == "" {
		return errors.New("relationship is required")
	}
	if req.StudentID == "" {
		return errors.New("student ID is required")
	}
	return nil
}

// toGuardianResponse converts model to response DTO
func (s *GuardianService) toGuardianResponse(guardian *models.Guardian) *dto.GuardianResponse {
	response := &dto.GuardianResponse{
		ID:           guardian.ID.String(),
		Occupation:   guardian.Occupation,
		Relationship: guardian.Relationship,
		Address:      guardian.Address,
		StudentID:    guardian.StudentID.String(),
		UserID:       guardian.UserID.String(),
		IsPrimary:    guardian.IsPrimary,
		IsEmergency:  guardian.IsEmergency,
		Status:       guardian.Status,
		CreatedBy:    guardian.CreatedBy.String(),
		CreatedAt:    guardian.CreatedAt,
		UpdatedAt:    guardian.UpdatedAt,
	}

	// Add student details if preloaded
	if guardian.Student.ID != uuid.Nil {
		response.Student = s.toUserResponse(&guardian.Student)
	}

	// Add user details if preloaded (the guardian's own user account)
	if guardian.User.ID != uuid.Nil {
		response.User = s.toUserResponse(&guardian.User)
	}

	return response
}

// toUserResponse converts User model to UserResponse DTO
func (s *GuardianService) toUserResponse(user *models.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:          user.ID.String(),
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		MiddleName:  user.MiddleName,
		FullName:    strings.TrimSpace(user.FirstName + " " + user.MiddleName + " " + user.LastName),
		Email:       user.Email,
		Phone:       user.Phone,
		Role:        user.Role,
		Position:    user.Position,
		Picture:     user.Picture,
		IsVerified:  user.IsVerified,
		Location:    user.Location,
		LastLoginAt: user.LastLoginAt,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}
}
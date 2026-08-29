// dto/arm_class_teacher_dto.go
package dto

import (
	"time"
)

// CreateArmClassTeacherRequest represents the request to assign a class teacher to an arm
type CreateArmClassTeacherRequest struct {
	ArmID     string `json:"arm_id" binding:"required"`
	TeacherID string `json:"teacher_id" binding:"required"`
	Status    string `json:"status"` // active, inactive (default: active)
}

// UpdateArmClassTeacherRequest represents the request to update an arm class teacher assignment
type UpdateArmClassTeacherRequest struct {
	ArmID     string `json:"arm_id"`
	TeacherID string `json:"teacher_id"`
	Status    string `json:"status"`
}

// ArmClassTeacherResponse represents the response for an arm class teacher assignment
type ArmClassTeacherResponse struct {
	ID        string    `json:"id"`
	ArmID     string    `json:"arm_id"`
	TeacherID string    `json:"teacher_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Nested relationships
	Arm     *ArmResponse       `json:"arm,omitempty"`
	Teacher *UserResponse      `json:"teacher,omitempty"`
	Creator *UserResponse      `json:"creator,omitempty"`
}

// ArmClassTeacherListResponse represents a paginated list of arm class teacher assignments
type ArmClassTeacherListResponse struct {
	Assignments []ArmClassTeacherResponse `json:"assignments"`
	Total       int64                     `json:"total"`
	Page        int                       `json:"page"`
	Limit       int                       `json:"limit"`
	TotalPages  int                       `json:"total_pages"`
}

// ArmClassTeacherQueryParams represents query parameters for filtering assignments
type ArmClassTeacherQueryParams struct {
	ArmID     string `form:"arm_id"`
	TeacherID string `form:"teacher_id"`
	Status    string `form:"status"`
	Page      int    `form:"page" default:"1"`
	Limit     int    `form:"limit" default:"20"`
	SortBy    string `form:"sort_by" default:"created_at"`
	SortOrder string `form:"sort_order" default:"desc"`
}

// ArmWithClassTeacher represents an arm with its class teacher information
type ArmWithClassTeacher struct {
	ArmID        string `json:"arm_id"`
	ArmName      string `json:"arm_name"`
	ArmCode      string `json:"arm_code"`
	TeacherID    string `json:"teacher_id,omitempty"`
	TeacherName  string `json:"teacher_name,omitempty"`
	TeacherEmail string `json:"teacher_email,omitempty"`
	Status       string `json:"status"`
}

// BulkAssignClassTeachersRequest represents the request to bulk assign class teachers
type BulkAssignClassTeachersRequest struct {
	Assignments []struct {
		ArmID     string `json:"arm_id" binding:"required"`
		TeacherID string `json:"teacher_id" binding:"required"`
	} `json:"assignments" binding:"required,min=1"`
	Status string `json:"status"`
}

// BulkAssignResult represents the result of a bulk assignment operation
type BulkAssignResult struct {
	SuccessCount int                           `json:"success_count"`
	FailedCount  int                           `json:"failed_count"`
	Created      []ArmClassTeacherResponse     `json:"created"`
	Errors       []BulkAssignError             `json:"errors"`
}

// BulkAssignError represents an error in bulk assignment
type BulkAssignError struct {
	ArmID   string `json:"arm_id"`
	TeacherID string `json:"teacher_id"`
	Error   string `json:"error"`
}
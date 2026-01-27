package models

import (
    "time"
    
    "github.com/google/uuid"
	"gorm.io/datatypes"
	)

type ActivityLog struct {
    ID         uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Who performed the action
    UserID     uuid.UUID      `gorm:"type:uuid;not null;index"`
    
    // What action was performed
    Action     string         `gorm:"type:varchar(100);not null;index"`
    
    // What entity was affected
    EntityID   uuid.UUID      `gorm:"type:uuid;index"`          // ID of the affected entity
    EntityType string         `gorm:"type:varchar(100);index"`  // Type of entity (e.g., "assignment", "user", "course")
    
    // Additional context
    Details    string         `gorm:"type:text"`                // Human-readable description
    IPAddress  string         `gorm:"type:varchar(45)"`         // IPv4 or IPv6 address
    UserAgent  string         `gorm:"type:text"`                // Browser/device info
    Metadata   datatypes.JSON `gorm:"type:jsonb"`               // Additional structured data
    
    // Timestamps
    CreatedAt  time.Time      `gorm:"index"`
    
    // Relationships
    User       User           `gorm:"foreignKey:UserID"`
}

// Predefined actions for consistency
const (
    ActionLogin               = "user_login"
    ActionLogout              = "user_logout"
    ActionPasswordChange      = "password_change"
    ActionProfileUpdate       = "profile_update"
    
    ActionCourseCreate        = "course_create"
    ActionCourseUpdate        = "course_update"
    ActionCourseDelete        = "course_delete"
    ActionCourseEnroll        = "course_enroll"
    ActionCourseUnenroll      = "course_unenroll"
    
    ActionAssignmentCreate    = "assignment_create"
    ActionAssignmentUpdate    = "assignment_update"
    ActionAssignmentDelete    = "assignment_delete"
    ActionAssignmentSubmit    = "assignment_submitted"
    ActionAssignmentGrade     = "assignment_graded"
    
    ActionMaterialCreate      = "material_create"
    ActionMaterialUpdate      = "material_update"
    ActionMaterialDelete      = "material_delete"
    
    ActionPaymentSuccess      = "payment_success"
    ActionPaymentFailed       = "payment_failed"
    ActionSubscriptionStart   = "subscription_start"
    ActionSubscriptionEnd     = "subscription_end"
    
    ActionSystemBackup        = "system_backup"
    ActionSystemRestore       = "system_restore"
    ActionSystemMaintenance   = "system_maintenance"
)
package models

import (
	"time"
	"github.com/google/uuid"

)
type LiveClass struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Course Relationship
    CourseID           uuid.UUID      `gorm:"type:uuid;not null;index"`
    ChapterID          *uuid.UUID     `gorm:"type:uuid;index"` // Optional chapter association
    ModuleID           *uuid.UUID     `gorm:"type:uuid;index"` // Optional module association
    
    // Class Identification
    Title              string         `gorm:"type:varchar(255);not null"`
    Description        string         `gorm:"type:text"`
    Slug               string         `gorm:"type:varchar(300);uniqueIndex;not null"`
    ShortDescription   string         `gorm:"type:varchar(500)"`
    
    // Scheduling
    StartTime          time.Time      `gorm:"not null;index"` // Scheduled start time
    EndTime            time.Time      `gorm:"not null;index"` // Scheduled end time
    Duration           int            `gorm:"default:60"` // Minutes
    Timezone           string         `gorm:"type:varchar(50);default:'UTC'"`
    RecurrenceRule     string         `gorm:"type:varchar(200)"` // RRULE for recurring classes
    RecurrenceType     string         `gorm:"type:varchar(20);default:'none';check:recurrence_type IN ('none', 'daily', 'weekly', 'monthly')"`
    RecurrenceEndDate  *time.Time     // End date for recurring series
    
    // Host Information
    TutorID            uuid.UUID      `gorm:"type:uuid;not null;index"` // Primary instructor
    CoTutorIDs         []uuid.UUID    `gorm:"type:uuid[]"` // Additional instructors
    HostNotes          string         `gorm:"type:text"` // Private notes for host
    
    // Capacity & Access
    MaxAttendees       int            `gorm:"default:100"` // Maximum participants
    MinAttendees       int            `gorm:"default:1"`   // Minimum to not cancel
    WaitlistEnabled    bool           `gorm:"default:true"`
    WaitlistCapacity   int            `gorm:"default:20"` // Additional waitlist spots
    AccessLevel        string         `gorm:"type:varchar(20);default:'enrolled';check:access_level IN ('enrolled', 'premium', 'invite_only', 'public')"`
    RequiresApproval   bool           `gorm:"default:false"` // Manual approval needed
    
    // Meeting Configuration
    Platform           string         `gorm:"type:varchar(50);default:'zoom';check:platform IN ('zoom', 'teams', 'google_meet', 'custom', 'bigbluebutton', 'jitsi')"`
    MeetingID          string         `gorm:"type:varchar(200)"` // Platform meeting ID
    MeetingURL         string         `gorm:"type:varchar(500)"` // Join URL
    MeetingPassword    string         `gorm:"type:varchar(100)"` // Join password
    
    // Preparation & Materials
    Agenda             string         `gorm:"type:text"` // Class agenda
    LearningObjectives []string       `gorm:"type:text[]"` // Session objectives
    Prerequisites      []string       `gorm:"type:text[]"` // Required preparation
    
    // Technical Requirements
    TechRequirements   []string       `gorm:"type:text[]"` // Software/hardware needed
    RecommendedSetup   string         `gorm:"type:text"` // Setup instructions
    TestURL            string         `gorm:"type:varchar(500)"` // Test meeting URL
    

    
    // Recording Settings
    RecordAutomatically bool          `gorm:"default:false"`
    RecordingStorage   string         `gorm:"type:varchar(50);default:'platform';check:recording_storage IN ('platform', 's3', 'gcs', 'local')"`
    AutoPublishRecordings bool        `gorm:"default:false"`
    RecordingRetention int            `gorm:"default:30"` // Days to keep recordings
    
    // Relationships
    // Course             Course         `gorm:"foreignKey:CourseID"`
    // Chapter            Chapter        `gorm:"foreignKey:ChapterID"`
    // Tutor              User           `gorm:"foreignKey:TutorID"`

    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
    ApprovedAt         *time.Time
    CancelledAt        *time.Time
}

// TableName specifies the table name
func (LiveClass) TableName() string {
    return "live_classes"
}
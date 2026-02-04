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
    TopicID            *uuid.UUID     `gorm:"type:uuid;index"` // Optional topic association
    LessonID          *uuid.UUID     `gorm:"type:uuid;index"` // Optional lesson association
    // Class Identification
    Title              string         `gorm:"type:varchar(255);not null"`
    Description        string         `gorm:"type:text"`
    Slug               string         `gorm:"type:varchar(300);uniqueIndex;not null"`
    
    // Scheduling
    StartTime          time.Time      `gorm:"not null;index"` // Scheduled start time
    EndTime            time.Time      `gorm:"not null;index"` // Scheduled end time
    Duration           int            `gorm:"default:60"` // Minutes
    Timezone           string         `gorm:"type:varchar(50);default:'UTC'"`
  
    // Host Information
    TutorID            uuid.UUID      `gorm:"type:uuid;not null;index"` // Primary instructor
    HostNotes          string         `gorm:"type:text"` // Private notes for host
    
    // Capacity & Access
    MaxAttendees       int            `gorm:"default:100"` // Maximum participants
    MinAttendees       int            `gorm:"default:1"`   // Minimum to not cancel
    WaitlistEnabled    bool           `gorm:"default:true"`
    WaitlistCapacity   int            `gorm:"default:20"` // Additional waitlist spots
    AccessLevel        string         `gorm:"type:varchar(20);default:'enrolled';check:access_level IN ('enrolled', 'premium', 'invite_only', 'public')"`
    
    // Meeting Configuration
    Platform           string         `gorm:"type:varchar(50);default:'zoom';check:platform IN ('zoom', 'teams', 'google_meet', 'custom', 'bigbluebutton', 'jitsi')"`
    MeetingID          string         `gorm:"type:varchar(200)"` // Platform meeting ID
    MeetingURL         string         `gorm:"type:varchar(500)"` // Join URL
    MeetingPassword    string         `gorm:"type:varchar(100)"` // Join password
    
    // Preparation & Materials
    Agenda             string         `gorm:"type:text"` // Class agenda

    RecommendedSetup   string         `gorm:"type:text"` // Setup instructions
    TestURL            string         `gorm:"type:varchar(500)"` // Test meeting URL
    

    
    // Recording Settings
    RecordAutomatically bool          `gorm:"default:false"`
    RecordingStorage   string         `gorm:"type:varchar(50);default:'platform';check:recording_storage IN ('platform', 's3', 'gcs', 'local')"`
    AutoPublishRecordings bool        `gorm:"default:false"`
    
    // Relationships
    Course             Course         `gorm:"foreignKey:CourseID"`
    Chapter            Chapter        `gorm:"foreignKey:ChapterID"`
    Tutor              User           `gorm:"foreignKey:TutorID"`

    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
    IsCancelled        *bool      `json:"is_cancelled" default:"false"` // Set to true to cancel

}


// LiveClassInput - for creating live classes
type LiveClassInput struct {
    // Required fields
    CourseID         uuid.UUID   `json:"course_id" binding:"required"`
    Title            string      `json:"title" binding:"required,min=3,max=255"`
    StartTime        time.Time   `json:"start_time" binding:"required"`
    EndTime          time.Time   `json:"end_time" binding:"required"`
    TutorID          uuid.UUID   `json:"tutor_id" binding:"required"`
    
    // Optional relationships
    ChapterID        *uuid.UUID  `json:"chapter_id"`
    TopicID          *uuid.UUID  `json:"topic_id"`
    LessonID         *uuid.UUID  `json:"lesson_id"`
    
    // Optional details
    Description      string      `json:"description" binding:"max=2000"`
    Duration         int         `json:"duration" binding:"omitempty"` 
    Timezone         string      `json:"timezone" binding:"omitempty,timezone"`

    // Capacity & Access
    MaxAttendees     int         `json:"max_attendees" binding:"min=1,max=1000"`
    MinAttendees     int         `json:"min_attendees" binding:"min=1,max=1000"`
    WaitlistEnabled  bool        `json:"waitlist_enabled"`
    WaitlistCapacity int         `json:"waitlist_capacity" binding:"min=0,max=100"`
    AccessLevel      string      `json:"access_level" binding:"omitempty,oneof=enrolled premium invite_only public"`
    RequiresApproval bool        `json:"requires_approval"`
    
    // Meeting Platform
    Platform         string      `json:"platform" binding:"omitempty,oneof=zoom teams google_meet custom bigbluebutton jitsi"`
    
    // Preparation & Materials
    Agenda           string      `json:"agenda" binding:"max=5000"`
    RecommendedSetup string      `json:"recommended_setup" binding:"max=2000"`
    
    // Recording
    RecordAutomatically bool     `json:"record_automatically"`
    RecordingStorage   string    `json:"recording_storage" binding:"omitempty,oneof=platform s3 gcs local"`
    AutoPublishRecordings bool   `json:"auto_publish_recordings"`
    
    // Host notes
    HostNotes        string      `json:"host_notes" binding:"max=2000"`
}

// LiveClassResponse - for API responses
type LiveClassResponse struct {
    ID                 uuid.UUID      `json:"id"`
    CourseID           uuid.UUID      `json:"course_id"`
    CourseName         string         `json:"course_name,omitempty"`
    ChapterID          *uuid.UUID     `json:"chapter_id,omitempty"`
    TopicID            *uuid.UUID     `json:"topic_id,omitempty"`
    LessonID           *uuid.UUID     `json:"lesson_id,omitempty"`
    Title              string         `json:"title"`
    Description        string         `json:"description"`
    Slug               string         `json:"slug"`
    StartTime          time.Time      `json:"start_time"`
    EndTime            time.Time      `json:"end_time"`
    Duration           int            `json:"duration"`
    Timezone           string         `json:"timezone"`
    TutorID            uuid.UUID      `json:"tutor_id"`
    TutorName          string         `json:"tutor_name,omitempty"`
    MaxAttendees       int            `json:"max_attendees"`
    MinAttendees       int            `json:"min_attendees"`
    WaitlistEnabled    bool           `json:"waitlist_enabled"`
    WaitlistCapacity   int            `json:"waitlist_capacity"`
    AccessLevel        string         `json:"access_level"`
    RequiresApproval   bool           `json:"requires_approval"`
    Platform           string         `json:"platform"`
    MeetingID          string         `json:"meeting_id,omitempty"`
    MeetingURL         string         `json:"meeting_url,omitempty"`
    MeetingPassword    string         `json:"meeting_password,omitempty"`
    Agenda             string         `json:"agenda,omitempty"`
    LearningObjectives []string       `json:"learning_objectives,omitempty"`
    Prerequisites      []string       `json:"prerequisites,omitempty"`
    TechRequirements   []string       `json:"tech_requirements,omitempty"`
    RecommendedSetup   string         `json:"recommended_setup,omitempty"`
    TestURL            string         `json:"test_url,omitempty"`
    RecordAutomatically bool          `json:"record_automatically"`
    RecordingStorage   string         `json:"recording_storage"`
    AutoPublishRecordings bool        `json:"auto_publish_recordings"`
    HostNotes          string         `json:"host_notes,omitempty"`
    Status             string         `json:"status"` // scheduled, ongoing, completed, cancelled
    CreatedAt          time.Time      `json:"created_at"`
    UpdatedAt          time.Time      `json:"updated_at"`
    
    // Computed fields
    AvailableSeats     int            `json:"available_seats,omitempty"`
    TotalEnrolled      int            `json:"total_enrolled,omitempty"`
    IsUpcoming         bool           `json:"is_upcoming"`
    IsLiveNow          bool           `json:"is_live_now"`
}


// LiveClassUpdateInput - for updating live classes
type LiveClassUpdateInput struct {
    // Basic info (only updatable before start)
    Title            *string    `json:"title" binding:"omitempty,min=3,max=255"`
    Description      *string    `json:"description"`
    
    // Rescheduling (only before start)
    StartTime        *string    `json:"start_time"` // String for parsing
    EndTime          *string    `json:"end_time"`   // String for parsing
    Duration         *int       `json:"duration" binding:"omitempty"`
    Timezone         *string    `json:"timezone" binding:"omitempty,timezone"`
    
    // Relationships (only before start)
    ChapterID        *uuid.UUID `json:"chapter_id"`
    TopicID          *uuid.UUID `json:"topic_id"`
    LessonID         *uuid.UUID `json:"lesson_id"`
    TutorID          *uuid.UUID `json:"tutor_id"`
    
    // Capacity (only before start)
    MaxAttendees     *int       `json:"max_attendees" binding:"omitempty,min=1,max=1000"`
    MinAttendees     *int       `json:"min_attendees" binding:"omitempty,min=1,max=1000"`
    WaitlistEnabled  *bool      `json:"waitlist_enabled"`
    WaitlistCapacity *int       `json:"waitlist_capacity" binding:"omitempty,min=0,max=100"`
    
    // Access control (only before start)
    AccessLevel      *string    `json:"access_level" binding:"omitempty,oneof=enrolled premium invite_only public"`
    
    // Meeting platform (only before start)
    Platform         *string    `json:"platform" binding:"omitempty,oneof=zoom teams google_meet custom bigbluebutton jitsi"`
    
    // Content (can update anytime)
    Agenda           *string    `json:"agenda"`
    RecommendedSetup *string    `json:"recommended_setup"`
    HostNotes        *string    `json:"host_notes"`
    
    // Recording (can update anytime before class ends)
    RecordAutomatically *bool   `json:"record_automatically"`
    RecordingStorage   *string  `json:"recording_storage" binding:"omitempty,oneof=platform s3 gcs local"`
    AutoPublishRecordings *bool `json:"auto_publish_recordings"`
    
    // Cancellation (can update anytime)
    IsCancelled        *bool      `json:"is_cancelled" default:"false"` // Set to true to cancel

}

// TableName specifies the table name
func (LiveClass) TableName() string {
    return "live_classes"
}
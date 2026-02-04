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

    RecommendedSetup   string         `gorm:"type:text"` // Setup instructions
    TestURL            string         `gorm:"type:varchar(500)"` // Test meeting URL
    

    
    // Recording Settings
    RecordAutomatically bool          `gorm:"default:false"`
    RecordingStorage   string         `gorm:"type:varchar(50);default:'platform';check:recording_storage IN ('platform', 's3', 'gcs', 'local')"`
    AutoPublishRecordings bool        `gorm:"default:false"`
    RecordingRetention int            `gorm:"default:30"` // Days to keep recordings
    
    // Relationships
    Course             Course         `gorm:"foreignKey:CourseID"`
    Chapter            Chapter        `gorm:"foreignKey:ChapterID"`
    Tutor              User           `gorm:"foreignKey:TutorID"`

    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
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
    ShortDescription string      `json:"short_description" binding:"max=500"`
    Duration         int         `json:"duration" binding:"min=15,max=480"` // 15 min to 8 hours
    Timezone         string      `json:"timezone" binding:"omitempty,timezone"`
    
    // Recurrence
    RecurrenceType   string      `json:"recurrence_type" binding:"omitempty,oneof=none daily weekly monthly"`
    RecurrenceRule   string      `json:"recurrence_rule" binding:"omitempty,max=200"`
    RecurrenceEndDate *time.Time `json:"recurrence_end_date"`
    
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
    LearningObjectives []string  `json:"learning_objectives"`
    Prerequisites    []string    `json:"prerequisites"`
    TechRequirements []string    `json:"tech_requirements"`
    RecommendedSetup string      `json:"recommended_setup" binding:"max=2000"`
    
    // Recording
    RecordAutomatically bool     `json:"record_automatically"`
    RecordingStorage   string    `json:"recording_storage" binding:"omitempty,oneof=platform s3 gcs local"`
    AutoPublishRecordings bool   `json:"auto_publish_recordings"`
    RecordingRetention int       `json:"recording_retention" binding:"min=1,max=365"`
    
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
    ShortDescription   string         `json:"short_description"`
    StartTime          time.Time      `json:"start_time"`
    EndTime            time.Time      `json:"end_time"`
    Duration           int            `json:"duration"`
    Timezone           string         `json:"timezone"`
    RecurrenceType     string         `json:"recurrence_type"`
    RecurrenceRule     string         `json:"recurrence_rule,omitempty"`
    RecurrenceEndDate  *time.Time     `json:"recurrence_end_date,omitempty"`
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
    RecordingRetention int            `json:"recording_retention"`
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

// TableName specifies the table name
func (LiveClass) TableName() string {
    return "live_classes"
}
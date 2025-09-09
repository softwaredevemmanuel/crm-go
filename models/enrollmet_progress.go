package models

import (
	"time"
	"github.com/google/uuid"


)

type EnrollmentProgress struct {
    ID                 uuid.UUID      `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
    
    // Core Relationships
    EnrollmentID       uuid.UUID      `gorm:"type:uuid;not null;index"`
    UserID             uuid.UUID      `gorm:"type:uuid;not null;index"`
    CourseID           uuid.UUID      `gorm:"type:uuid;not null;index"`
    
    // Progress Tracking
    OverallProgress    float64        `gorm:"type:decimal(5,2);default:0;check:overall_progress >= 0 AND overall_progress <= 100"` // Overall completion percentage
    CurrentChapterID   *uuid.UUID     `gorm:"type:uuid;index"`        // Current chapter student is on
    CurrentLessonID    *uuid.UUID     `gorm:"type:uuid;index"`        // Current lesson student is on
    LastActivityAt     time.Time      `gorm:"default:CURRENT_TIMESTAMP;index"` // Last time student was active
    
    // Time Tracking
    TotalTimeSpent     int            `gorm:"default:0"`              // Total seconds spent in course
    EstimatedTimeLeft  int            `gorm:"default:0"`              // Estimated seconds to complete
    FirstAccessAt      *time.Time     `gorm:"index"`                  // First time student accessed
    LastCompletedAt    *time.Time     `gorm:"index"`                  // Last completion activity
    
    // Completion Metrics
    CompletedChapters  int            `gorm:"default:0"`              // Number of chapters completed
    TotalChapters      int            `gorm:"default:0"`              // Total chapters in course
    CompletedLessons   int            `gorm:"default:0"`              // Number of lessons completed
    TotalLessons       int            `gorm:"default:0"`              // Total lessons in course
    CompletedQuizzes   int            `gorm:"default:0"`              // Quizzes completed
    TotalQuizzes       int            `gorm:"default:0"`              // Total quizzes in course
    CompletedAssignments int          `gorm:"default:0"`              // Assignments submitted
    TotalAssignments   int            `gorm:"default:0"`              // Total assignments
    

    
    // Engagement Metrics
    VideoWatchTime     int            `gorm:"default:0"`              // Total seconds of video watched
    VideoCompletionRate float64       `gorm:"type:decimal(5,2);default:0"`    // % of videos completed
    ReadingTime        int            `gorm:"default:0"`              // Total seconds spent reading
    NotesTaken         int            `gorm:"default:0"`              // Number of notes taken
    BookmarksCreated   int            `gorm:"default:0"`              // Number of bookmarks
    DiscussionsParticipated int       `gorm:"default:0"`              // Forum discussions participated in
    
    // Progress Status
    Status             string         `gorm:"type:varchar(20);default:'in_progress';check:status IN ('not_started', 'in_progress', 'completed', 'behind_schedule', 'ahead_of_schedule', 'paused', 'dropped')"`
    CompletionStatus   string         `gorm:"type:varchar(20);default:'incomplete';check:completion_status IN ('incomplete', 'completed', 'passed', 'failed', 'certified')"`
    IsBehindSchedule   bool           `gorm:"default:false"`          // Based on expected pace
    IsAheadOfSchedule  bool           `gorm:"default:false"`          // Based on expected pace
    
    // Goal Tracking
    DailyGoal          int            `gorm:"default:0"`              // User's daily goal in minutes
    GoalAchievementRate float64       `gorm:"type:decimal(5,2);default:0"`    // % of goals achieved
    WeeklyTarget       float64        `gorm:"type:decimal(5,2);default:0"`    // Weekly progress target %
    TargetAchievement  float64        `gorm:"type:decimal(5,2);default:0"`    // % of targets achieved
    
    // Certification & Awards
    CertificateEarned  bool           `gorm:"default:false"`
    CertificateID      *uuid.UUID     `gorm:"type:uuid;index"`
    BadgesEarned       int            `gorm:"default:0"`              // Number of badges earned
    

    
    // Relationships
    Enrollment         Enrollment     `gorm:"foreignKey:EnrollmentID"`
    User               User           `gorm:"foreignKey:UserID"`
    Course             Course         `gorm:"foreignKey:CourseID"`
    Certificate        Certificate    `gorm:"foreignKey:CertificateID"`
    
    // Timestamps
    CreatedAt          time.Time
    UpdatedAt          time.Time
}

// TableName specifies the table name
func (EnrollmentProgress) TableName() string {
    return "enrollment_progress"
}
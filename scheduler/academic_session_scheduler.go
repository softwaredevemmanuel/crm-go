// scheduler/academic_session_scheduler.go
package scheduler

import (
	"log"
	"time"

	"gorm.io/gorm"

	"crm-go/models"
)

type AcademicSessionScheduler struct {
	db     *gorm.DB
	config SchedulerConfig
}

// SchedulerConfig holds the configuration for the scheduler
type SchedulerConfig struct {
	CheckInterval time.Duration
	TimeZone      string
	NotifyAdmins  bool
	AutoStart     bool
	AutoComplete  bool
}

// NewAcademicSessionScheduler creates a new scheduler instance
func NewAcademicSessionScheduler(db *gorm.DB) *AcademicSessionScheduler {
	return &AcademicSessionScheduler{
		db: db,
		config: SchedulerConfig{
			CheckInterval: 24 * time.Hour,
			TimeZone:      "Africa/Lagos",
			NotifyAdmins:  false,
			AutoStart:     true,
			AutoComplete:  true,
		},
	}
}

// NewAcademicSessionSchedulerWithConfig creates a scheduler with custom config
func NewAcademicSessionSchedulerWithConfig(db *gorm.DB, config SchedulerConfig) *AcademicSessionScheduler {
	return &AcademicSessionScheduler{
		db:     db,
		config: config,
	}
}

// GetDefaultConfig returns the default scheduler configuration
func GetDefaultConfig() SchedulerConfig {
	return SchedulerConfig{
		CheckInterval: 1 * time.Hour,
		TimeZone:      "Africa/Lagos",
		NotifyAdmins:  false,
		AutoStart:     true,
		AutoComplete:  true,
	}
}

// Start starts the academic session scheduler
func (s *AcademicSessionScheduler) Start() {
	log.Println("🚀 Starting Academic Session Scheduler...")
	log.Printf("📋 Config: Interval=%v, AutoStart=%v, AutoComplete=%v", 
		s.config.CheckInterval, s.config.AutoStart, s.config.AutoComplete)

	// Run immediately on start
	s.processAcademicSessions()

	// Run every hour
	ticker := time.NewTicker(s.config.CheckInterval)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			s.processAcademicSessions()
		}
	}()

	log.Println("✅ Academic Session Scheduler started successfully")
}

// processAcademicSessions handles all academic session updates
func (s *AcademicSessionScheduler) processAcademicSessions() {
	now := time.Now()
	log.Printf("🔄 Running scheduler check at %s", now.Format("2006-01-02 15:04:05"))

	// 1. Complete expired active sessions
	if s.config.AutoComplete {
		s.completeExpiredSessions(now)
	}

	// 2. Start pending sessions
	if s.config.AutoStart {
		s.startPendingSessions(now)
	}

	// 3. Update current session
	s.updateCurrentSession(now)
}

// completeExpiredSessions marks active sessions as completed when end_date passes
func (s *AcademicSessionScheduler) completeExpiredSessions(now time.Time) {
	result := s.db.Model(&models.AcademicSession{}).
		Where("end_date < ? AND status = ? AND deleted_at IS NULL", now, "active").
		Updates(map[string]interface{}{
			"status":     "completed",
			"is_current": false,
			"updated_at": now,
		})

	if result.Error != nil {
		log.Printf("❌ Failed to update expired academic sessions: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("✅ Completed %d expired academic session(s)", result.RowsAffected)
		
		// Log the sessions that were completed
		var sessions []models.AcademicSession
		s.db.Where("end_date < ? AND status = ? AND deleted_at IS NULL", now, "completed").
			Limit(int(result.RowsAffected)).
			Find(&sessions)
		
		for _, session := range sessions {
			log.Printf("   📅 Completed: %s (%s)", session.AcademicYear, session.Code)
		}
	}
}

// startPendingSessions activates sessions that are scheduled to start
func (s *AcademicSessionScheduler) startPendingSessions(now time.Time) {
	// Check if there's already an active session
	var activeCount int64
	if err := s.db.Model(&models.AcademicSession{}).
		Where("status = ? AND deleted_at IS NULL", "active").
		Count(&activeCount).Error; err != nil {
		log.Printf("❌ Failed to check active sessions: %v", err)
		return
	}

	// Don't start new sessions if there's already an active one
	if activeCount > 0 {
		return
	}

	// Find the next session to start
	var session models.AcademicSession
	if err := s.db.Where("start_date <= ? AND status = ? AND deleted_at IS NULL", 
		now, "pending").
		Order("start_date ASC").
		First(&session).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("❌ Failed to find pending sessions: %v", err)
		}
		return
	}

	// Update the session to active
	result := s.db.Model(&session).Updates(map[string]interface{}{
		"status":     "active",
		"is_current": true,
		"updated_at": now,
	})

	if result.Error != nil {
		log.Printf("❌ Failed to start pending session: %v", result.Error)
		return
	}

	if result.RowsAffected > 0 {
		log.Printf("✅ Started new academic session: %s (%s)", session.AcademicYear, session.Code)
		log.Printf("   📅 Period: %s to %s", 
			session.StartDate.Format("2006-01-02"), 
			session.EndDate.Format("2006-01-02"))
	}
}

// updateCurrentSession ensures only one session is marked as current
func (s *AcademicSessionScheduler) updateCurrentSession(now time.Time) {
	// Find the active session
	var activeSession models.AcademicSession
	if err := s.db.Where("status = ? AND deleted_at IS NULL", "active").
		First(&activeSession).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Printf("❌ Failed to find active session: %v", err)
		}
		return
	}

	// Check if there are any sessions incorrectly marked as current
	var currentSessions []models.AcademicSession
	if err := s.db.Where("is_current = ? AND id != ? AND deleted_at IS NULL", 
		true, activeSession.ID).
		Find(&currentSessions).Error; err != nil {
		log.Printf("❌ Failed to find duplicate current sessions: %v", err)
		return
	}

	// Fix duplicate current sessions
	if len(currentSessions) > 0 {
		for _, session := range currentSessions {
			s.db.Model(&session).Update("is_current", false)
		}
		log.Printf("✅ Fixed %d duplicate current session(s)", len(currentSessions))
	}
}

// GetCurrentSession returns the current academic session
func (s *AcademicSessionScheduler) GetCurrentSession() (*models.AcademicSession, error) {
	var session models.AcademicSession
	if err := s.db.Where("is_current = ? AND deleted_at IS NULL", true).
		First(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetActiveSessions returns all active academic sessions
func (s *AcademicSessionScheduler) GetActiveSessions() ([]models.AcademicSession, error) {
	var sessions []models.AcademicSession
	if err := s.db.Where("status = ? AND deleted_at IS NULL", "active").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetUpcomingSessions returns sessions that are scheduled to start
func (s *AcademicSessionScheduler) GetUpcomingSessions() ([]models.AcademicSession, error) {
	var sessions []models.AcademicSession
	if err := s.db.Where("status = ? AND deleted_at IS NULL", "pending").
		Order("start_date ASC").
		Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetSchedulerStatus returns the current status of the scheduler
func (s *AcademicSessionScheduler) GetSchedulerStatus() map[string]interface{} {
	status := make(map[string]interface{})
	
	var activeCount, pendingCount, completedCount int64
	
	s.db.Model(&models.AcademicSession{}).Where("status = ? AND deleted_at IS NULL", "active").Count(&activeCount)
	s.db.Model(&models.AcademicSession{}).Where("status = ? AND deleted_at IS NULL", "pending").Count(&pendingCount)
	s.db.Model(&models.AcademicSession{}).Where("status = ? AND deleted_at IS NULL", "completed").Count(&completedCount)
	
	status["active_sessions"] = activeCount
	status["pending_sessions"] = pendingCount
	status["completed_sessions"] = completedCount
	status["check_interval"] = s.config.CheckInterval.String()
	status["timezone"] = s.config.TimeZone
	status["auto_start"] = s.config.AutoStart
	status["auto_complete"] = s.config.AutoComplete
	
	return status
}
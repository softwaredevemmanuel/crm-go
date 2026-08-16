package seeds

import (
	"fmt"
	"log"
	"time"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func parseDate(date string) time.Time {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		panic(fmt.Sprintf("invalid date: %s", date))
	}

	return t
}

func SeedAcademicSessions() error {
	db := config.GetDB()

	adminID := uuid.MustParse(
		"fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd",
	)

	sessions := []models.AcademicSession{
		{
			ID:           uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			AcademicYear: "2026/2027",
			Code:         "2026-2027",
			StartDate:    parseDate("2026-09-01"),
			EndDate:      parseDate("2027-07-31"),
			Status:       "active",
			IsCurrent:    true,
			Description:  "2026/2027 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("686f42eb-e74c-49f3-9bd6-46b290fe0ef4"),
			AcademicYear: "2025/2026",
			Code:         "2025-2026",
			StartDate:    parseDate("2025-09-01"),
			EndDate:      parseDate("2026-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2025/2026 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			AcademicYear: "2024/2025",
			Code:         "2024-2025",
			StartDate:    parseDate("2024-09-01"),
			EndDate:      parseDate("2025-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2024/2025 Academic Session",
			CreatedBy:    adminID,
		},



	}

	for _, session := range sessions {

		// Check if the session already exists
		var existing models.AcademicSession

		result := db.Where(
			"id = ? OR academic_year = ? OR code = ?",
			session.ID,
			session.AcademicYear,
			session.Code,
		).First(&existing)

		if result.Error == nil {
			log.Printf(
				"⏭️ Academic Session already exists: %s",
				session.AcademicYear,
			)
			continue
		}

		// Create the session
		if err := db.Create(&session).Error; err != nil {
			return fmt.Errorf(
				"failed to seed academic session %s: %w",
				session.AcademicYear,
				err,
			)
		}

		log.Printf(
			"✅ Seeded Academic Session: %s",
			session.AcademicYear,
		)
	}

	return nil
}
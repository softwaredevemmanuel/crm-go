package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)



func SeedAcademicSessions() error {
	db := config.GetDB()

	adminID := uuid.MustParse(
		"fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd",
	)

	academic_sessions := []models.AcademicSession{
		{
			ID:           uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			AcademicYear: "2026/2027",
			Code:         "2026-2027",
			StartDate:    ParseDate("2026-09-01"),
			EndDate:      ParseDate("2027-07-31"),
			Status:       "active",
			IsCurrent:    true,
			Description:  "2026/2027 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("cb79e53b-bd13-431b-bd93-1ad9b1642d15"),
			AcademicYear: "2025/2026",
			Code:         "2025-2026",
			StartDate:    ParseDate("2025-09-01"),
			EndDate:      ParseDate("2026-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2026/2027 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("5c6c89bd-37a6-4afd-8874-65608ee75674"),
			AcademicYear: "2024/2025",
			Code:         "2024-2025",
			StartDate:    ParseDate("2024-09-01"),
			EndDate:      ParseDate("2025-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2026/2027 Academic Session",
			CreatedBy:    adminID,
		},
		

	}

	for _, academic_session := range academic_sessions {
		if err := db.Create(&academic_session).Error; err != nil {
			log.Printf("❌ Failed to seed academic session: %v", err)
		} else {
			log.Printf("✅ Seeded academic session: %s", academic_session.AcademicYear)
		}
	}

	return nil
}
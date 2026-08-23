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

	academic_sessions := []models.AcademicSession{
		{
			ID:           uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			AcademicYear: "2026/2027",
			Code:         "2026-2027",
			Term:         "first",
			StartDate:    parseDate("2026-09-01"),
			EndDate:      parseDate("2027-07-31"),
			Status:       "active",
			IsCurrent:    true,
			Description:  "2026/2027 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("cb79e53b-bd13-431b-bd93-1ad9b1642d15"),
			AcademicYear: "2026/2027",
			Code:         "2026-2027",
			Term:         "second",
			StartDate:    parseDate("2026-09-01"),
			EndDate:      parseDate("2027-07-31"),
			Status:       "completed",
			IsCurrent:    true,
			Description:  "2026/2027 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("5c6c89bd-37a6-4afd-8874-65608ee75674"),
			AcademicYear: "2026/2027",
			Code:         "2026-2027",
			Term:         "third",
			StartDate:    parseDate("2026-09-01"),
			EndDate:      parseDate("2027-07-31"),
			Status:       "completed",
			IsCurrent:    true,
			Description:  "2026/2027 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("686f42eb-e74c-49f3-9bd6-46b290fe0ef4"),
			AcademicYear: "2025/2026",
			Code:         "2025-2026",
			Term:         "first",
			StartDate:    parseDate("2025-09-01"),
			EndDate:      parseDate("2026-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2025/2026 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("dad56834-c36e-473c-82aa-d38ed4f7e44a"),
			AcademicYear: "2025/2026",
			Code:         "2025-2026",
			Term:         "second",
			StartDate:    parseDate("2025-09-01"),
			EndDate:      parseDate("2026-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2025/2026 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("f29c9e4e-ee97-46a3-a198-fa7a221dbd2d"),
			AcademicYear: "2025/2026",
			Code:         "2025-2026",
			Term:         "third",
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
			Term:         "first",
			StartDate:    parseDate("2024-09-01"),
			EndDate:      parseDate("2025-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2024/2025 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("af28892e-f6a5-4cea-8d12-4f609e9c8785"),
			AcademicYear: "2024/2025",
			Code:         "2024-2025",
			Term:         "second",
			StartDate:    parseDate("2024-09-01"),
			EndDate:      parseDate("2025-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2024/2025 Academic Session",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("47d7e25d-900e-4e86-91ef-d182cd742ab9"),
			AcademicYear: "2024/2025",
			Code:         "2024-2025",
			Term:         "third",
			StartDate:    parseDate("2024-09-01"),
			EndDate:      parseDate("2025-07-31"),
			Status:       "completed",
			IsCurrent:    false,
			Description:  "2024/2025 Academic Session",
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
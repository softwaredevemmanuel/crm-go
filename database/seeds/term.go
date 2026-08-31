package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
)

func SeedTerm() error {
	db := config.GetDB()

	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")
	academic_session1 := uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337")
	academic_session2 := uuid.MustParse("cb79e53b-bd13-431b-bd93-1ad9b1642d15")
	academic_session3 := uuid.MustParse("5c6c89bd-37a6-4afd-8874-65608ee75674")

	// prepare pointer dates since Term expects *time.Time for StartDate/EndDate
	// 2026/2027
firstStartDate2026 := ParseDate("2026-09-01")
firstEndDate2026 := ParseDate("2026-12-20")

secondStartDate2026 := ParseDate("2027-01-05")
secondEndDate2026 := ParseDate("2027-04-09")

thirdStartDate2026 := ParseDate("2027-04-26")
thirdEndDate2026 := ParseDate("2027-07-31")

// 2025/2026
firstStartDate2025 := ParseDate("2025-09-01")
firstEndDate2025 := ParseDate("2025-12-20")

secondStartDate2025 := ParseDate("2026-01-05")
secondEndDate2025 := ParseDate("2026-04-10")

thirdStartDate2025 := ParseDate("2026-04-27")
thirdEndDate2025 := ParseDate("2026-07-31")

// 2024/2025
firstStartDate2024 := ParseDate("2024-09-09")
firstEndDate2024 := ParseDate("2024-12-20")

secondStartDate2024 := ParseDate("2025-01-06")
secondEndDate2024 := ParseDate("2025-04-11")

thirdStartDate2024 := ParseDate("2025-04-28")
thirdEndDate2024 := ParseDate("2025-07-31")

	terms := []models.Term{
		{
			ID:                uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			AcademicSessionID: academic_session1,
			Name:              "First Term",
			Code:              "FT",
			TermNumber:        1,
			StartDate:         &firstStartDate2026,
			EndDate:           &firstEndDate2026,
			IsCurrent:         true,
			Status:            "active",
			Description:       "2026/2027 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("af28892e-f6a5-4cea-8d12-4f609e9c8785"),
			AcademicSessionID: academic_session1,
			Name:              "Second Term",
			Code:              "ST",
			TermNumber:        2,
			StartDate:         &secondStartDate2026,
			EndDate:           &secondEndDate2026,
			IsCurrent:         false,
			Status:            "inactive",
			Description:       "2026/2027 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("47d7e25d-900e-4e86-91ef-d182cd742ab9"),
			AcademicSessionID: academic_session1,
			Name:              "Third Term",
			Code:              "TT",
			TermNumber:        3,
			StartDate:         &thirdStartDate2026,
			EndDate:           &thirdEndDate2026,
			IsCurrent:         false,
			Status:            "inactive",
			Description:       "2026/2027 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("7f3a9c21-6d84-4b52-a1e7-93c8f05d2146"),
			AcademicSessionID: academic_session2,
			Name:              "First Term",
			Code:              "FT",
			TermNumber:        1,
			StartDate:         &firstStartDate2025,
			EndDate:           &firstEndDate2025,
			IsCurrent:         false,
			Status:            "completed",
			Description:       "2025/2026 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("b2e7f914-3a65-48dc-9f21-7c05e8a43619"),
			AcademicSessionID: academic_session2,
			Name:              "Second Term",
			Code:              "ST",
			TermNumber:        2,
			StartDate:         &secondStartDate2025,
			EndDate:           &secondEndDate2025,
			IsCurrent:         false,
			Status:            "completed",
			Description:       "2025/2026 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("d84c2a67-91f3-4e58-bb06-52a9d731c845"),
			AcademicSessionID: academic_session2,
			Name:              "Third Term",
			Code:              "TT",
			TermNumber:        3,
			StartDate:         &thirdStartDate2025,
			EndDate:           &thirdEndDate2025,
			IsCurrent:         false,
			Status:            "completed",
			Description:       "2025/2026 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("4a91e7c3-5d28-46fb-8c14-e963b7520a41"),
			AcademicSessionID: academic_session3,
			Name:              "First Term",
			Code:              "FT",
			TermNumber:        1,
			StartDate:         &firstStartDate2024,
			EndDate:           &firstEndDate2024,
			IsCurrent:         false,
			Status:            "completed",
			Description:       "2024/2025 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("c6f2d819-74a3-4b65-ae90-31d7c584926f"),
			AcademicSessionID: academic_session3,
			Name:              "Second Term",
			Code:              "ST",
			TermNumber:        2,
			StartDate:         &secondStartDate2024,
			EndDate:           &secondEndDate2024,
			IsCurrent:         false,
			Status:            "completed",
			Description:       "2024/2025 Academic Session",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("91b5e3a7-c462-4f18-8d29-75a0c634e9b2"),
			AcademicSessionID: academic_session3,
			Name:              "Third Term",
			Code:              "TT",
			TermNumber:        3,
			StartDate:         &thirdStartDate2024,
			EndDate:           &thirdEndDate2024,
			IsCurrent:         false,
			Status:            "completed",
			Description:       "2024/2025 Academic Session",
			CreatedBy:         adminID,
		},
	}

	for _, term := range terms {
		if err := db.Create(&term).Error; err != nil {
			log.Printf("❌ Failed to seed term: %v", err)
		} else {
			log.Printf("✅ Seeded academic term: %s", term.Name)
		}
	}

	return nil
}

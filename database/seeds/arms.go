package seeds

import (
	"fmt"
	"log"

	"crm-go/config"
	"crm-go/models"
	"gorm.io/gorm/clause"
	"github.com/google/uuid"
)

func SeedArms() error {
	db := config.GetDB()

	// Grade IDs
	jss1ID := uuid.MustParse("0f6d9ab4-2b7b-41e8-b823-3ba45103a899")
	jss2ID := uuid.MustParse("23064952-17b6-4b8d-9889-b44e35e46faf")
	jss3ID := uuid.MustParse("2dd8fe9c-3c2c-4668-bb6b-79fcd9e79bbe")

	ss1ID := uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046")
	ss2ID := uuid.MustParse("3a61e8d2-4f79-48b3-bc25-7e9d1a6f5082")
	ss3ID := uuid.MustParse("c84b1597-2e63-4a8f-91d4-6b7c3e205fa9")

	arms := []models.Arm{

		// =========================
		// JSS 1
		// =========================
		{
			ID:          uuid.MustParse("d7315f84-2a69-4c07-b913-568e1f427a35"),
			GradeID:     jss1ID,
			Name:        "A",
			Code:        "JSS1-A",
			Capacity:    100,
			Status:      "active",
			Description: "A class section of JSS1",
		},

		// =========================
		// JSS 2
		// =========================
		{
			ID:          uuid.MustParse("18b4e963-7d52-4fa1-a806-329c5e7184bd"),
			GradeID:     jss2ID,
			Name:        "Alpha",
			Code:        "A",
			Capacity:    15,
			Status:      "active",
			Description: "A class section of JSS2",
		},

		// =========================
		// JSS 3
		// =========================
		{
			ID:          uuid.MustParse("c85a2f17-6e43-49d0-b721-954a8c3f106e"),
			GradeID:     jss3ID,
			Name:        "Alpha",
			Code:        "A",
			Capacity:    15,
			Status:      "active",
			Description: "A class section of JSS3",
		},

		// =========================
		// SS 1
		// =========================
		{
			ID:          uuid.MustParse("73d9b526-1a84-4e60-8c35-617f2a948bcd"),
			GradeID:     ss1ID,
			Name:        "Alpha",
			Code:        "A",
			Capacity:    15,
			Status:      "active",
			Description: "A class section of SS1",
		},

		// =========================
		// SS 2
		// =========================
		{
			ID:          uuid.MustParse("f6c1e807-3b95-4a72-a514-826d9e4530bc"),
			GradeID:     ss2ID,
			Name:        "Alpha",
			Code:        "A",
			Capacity:    15,
			Status:      "active",
			Description: "A class section of SS2",
		},

		// =========================
		// SS 3
		// =========================
		{
			ID:          uuid.MustParse("2a57d914-8c36-4eb0-b751-649f1c382da5"),
			GradeID:     ss3ID,
			Name:        "Alpha",
			Code:        "A",
			Capacity:    15,
			Status:      "active",
			Description: "A class section of SS3",
		},

		
	}

for _, arm := range arms {
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"grade_id",
			"name",
			"code",
			"capacity",
			"status",
			"description",
		}),
	}).Create(&arm)

	if result.Error != nil {
		return fmt.Errorf(
			"failed to seed arm %s (%s): %w",
			arm.Name,
			arm.Code,
			result.Error,
		)
	}

	log.Printf(
		"✅ Seeded/updated arm: %s - %s",
		arm.Code,
		arm.Description,
	)
}

	return nil
}
package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func SeedSchemesOfWork() {
	db := config.GetDB()

	// Replace these with IDs from your existing seed data.
	biologyID := uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d")
	ss1ID := uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046")
	ss2ID := uuid.MustParse("3a61e8d2-4f79-48b3-bc25-7e9d1a6f5082")
	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	schemes := []models.SchemeOfWork{
		{
			ID:          uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146"),
			SubjectID:   biologyID,
			GradeID:     ss1ID,
			Term:        "first",
			Title:       "Biology Scheme of Work - SS1 First Term",
			Description: "First term Biology scheme of work for Senior Secondary School 1.",
			Status:      "published",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("c4e7b912-3a56-4d80-b135-628f9a1e4732"),
			SubjectID:   biologyID,
			GradeID:     ss1ID,
			Term:        "second",
			Title:       "Biology Scheme of Work - SS1 Second Term",
			Description: "Second term Biology scheme of work for Senior Secondary School 1.",
			Status:      "draft",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("1d8b6f43-9275-4ca1-a850-316e7d5942bc"),
			SubjectID:   biologyID,
			GradeID:     ss1ID,
			Term:        "third",
			Title:       "Biology Scheme of Work - SS1 Third Term",
			Description: "Third term Biology scheme of work for Senior Secondary School 1.",
			Status:      "draft",
			CreatedBy:   adminID,
		},

		{
			ID:          uuid.MustParse("7a2e9c51-4f83-46bd-b027-915c6a3e7481"),
			SubjectID:   biologyID,
			GradeID:     ss2ID,
			Term:        "first",
			Title:       "Biology Scheme of Work - SS2 First Term",
			Description: "First term Biology scheme of work for Senior Secondary School 2.",
			Status:      "published",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("b6d4f820-1c59-4e73-9a26-583b7d14e9c0"),
			SubjectID:   biologyID,
			GradeID:     ss2ID,
			Term:        "second",
			Title:       "Biology Scheme of Work - SS2 Second Term",
			Description: "Second term Biology scheme of work for Senior Secondary School 2.",
			Status:      "draft",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("3e91a745-6b28-4c50-b839-172d5f6a904e"),
			SubjectID:   biologyID,
			GradeID:     ss2ID,
			Term:        "third",
			Title:       "Biology Scheme of Work - SS2 Third Term",
			Description: "Third term Biology scheme of work for Senior Secondary School 2.",
			Status:      "draft",
			CreatedBy:   adminID,
		},
	}

	for _, scheme := range schemes {
		result := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"subject_id",
				"grade_id",
				"term",
				"title",
				"description",
				"status",
				"created_by",
			}),
		}).Create(&scheme)

		if result.Error != nil {
			log.Printf("❌ Failed to seed scheme %s: %v", scheme.Title, result.Error)
		} else {
			log.Printf("✅ Seeded/updated Scheme of Work: %s", scheme.Title)
		}
	}
}

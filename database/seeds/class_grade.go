package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedClassGrades() {
	db := config.GetDB()


	classGrades := []models.ClassGrade{
		{
			ID:           uuid.MustParse("0f6d9ab4-2b7b-41e8-b823-3ba45103a899"),
			Name:         "Junior Secondary School 1",
			Code:         "JSS1",
			Level:        7,
			Description:  "The first year of junior secondary education, focusing on building strong foundations in core subjects and essential learning skills.",
			Status:       "active",
		},
		{
			ID:           uuid.MustParse("23064952-17b6-4b8d-9889-b44e35e46faf"),
			Name:         "Junior Secondary School 2",
			Code:         "JSS2",
			Level:        8,
			Description:  "The second year of junior secondary education, where students strengthen their knowledge, develop deeper understanding, and improve their academic skills.",
			Status:       "active",
		},
		{
			ID:           uuid.MustParse("2dd8fe9c-3c2c-4668-bb6b-79fcd9e79bbe"),
			Name:         "Junior Secondary School 3",
			Code:         "JSS3",
			Level:        9,
			Description:  "The final year of junior secondary education, preparing students for senior secondary school through advanced learning and academic assessment.",
			Status:       "active",
		},
		{
			ID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Name:         "Senior Secondary School 1",
			Code:         "SS1",
			Level:        10,
			Description:  "The first year of senior secondary education, introducing students to more specialized subjects and preparing them for advanced academic study.",
			Status:       "active",
		},
		{
			ID:           uuid.MustParse("3a61e8d2-4f79-48b3-bc25-7e9d1a6f5082"),
			Name:         "Senior Secondary School 2",
			Code:         "SS2",
			Level:        11,
			Description:  "The second year of senior secondary education, where students deepen their subject knowledge and develop skills needed for examinations and future careers.",
			Status:       "active",
		},
		{
			ID:           uuid.MustParse("c84b1597-2e63-4a8f-91d4-6b7c3e205fa9"),
			Name:         "Senior Secondary School 3",
			Code:         "SS3",
			Level:        12,
			Description:  "The final year of senior secondary education, focused on completing the secondary curriculum and preparing students for major examinations and further education.",
			Status:       "active",
		},
		
		
	}

	for _, classGrade := range classGrades {
		if err := db.Create(&classGrade).Error; err != nil {
			log.Printf(
				"❌ Failed to seed class grade %s: %v",
				classGrade.Code,
				err,
			)
			continue
		}

		log.Printf(
			"✅ Seeded class grade: %s (%s)",
			classGrade.Name,
			classGrade.Code,
		)
	}
}
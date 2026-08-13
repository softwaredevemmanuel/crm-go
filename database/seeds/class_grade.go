package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedClassGrades() {
	db := config.GetDB()

	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	classGrades := []models.ClassGrade{
		{
			ID:           uuid.MustParse("0f6d9ab4-2b7b-41e8-b823-3ba45103a899"),
			Name:         "Junior Secondary School 1",
			Code:         "JSS1",
			Level:        1,
			Description:  "The first year of junior secondary education, focusing on building strong foundations in core subjects and essential learning skills.",
			AcademicYear: "2026/2027",
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("23064952-17b6-4b8d-9889-b44e35e46faf"),
			Name:         "Junior Secondary School 2",
			Code:         "JSS2",
			Level:        2,
			Description:  "The second year of junior secondary education, where students strengthen their knowledge, develop deeper understanding, and improve their academic skills.",
			AcademicYear: "2026/2027",
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("2dd8fe9c-3c2c-4668-bb6b-79fcd9e79bbe"),
			Name:         "Junior Secondary School 3",
			Code:         "JSS3",
			Level:        3,
			Description:  "The final year of junior secondary education, preparing students for senior secondary school through advanced learning and academic assessment.",
			AcademicYear: "2026/2027",
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("c37648bb-a2ea-4512-96b6-43d337c693da"),
			Name:         "Senior Secondary School 1",
			Code:         "SS1",
			Level:        4,
			Description:  "The first year of senior secondary education, introducing students to more specialized subjects and preparing them for advanced academic study.",
			AcademicYear: "2026/2027",
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("c91a8b19-add8-4841-a7ac-44c04cb6d7c4"),
			Name:         "Senior Secondary School 2",
			Code:         "SS2",
			Level:        5,
			Description:  "The second year of senior secondary education, where students deepen their subject knowledge and develop skills needed for examinations and future careers.",
			AcademicYear: "2026/2027",
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("f19aa519-6d40-46cd-977c-49ad96ca25ff"),
			Name:         "Senior Secondary School 3",
			Code:         "SS3",
			Level:        6,
			Description:  "The final year of senior secondary education, focused on completing the secondary curriculum and preparing students for major examinations and further education.",
			AcademicYear: "2026/2027",
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
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
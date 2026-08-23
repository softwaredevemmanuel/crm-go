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
	session6 := uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337")
	session5 := uuid.MustParse("686f42eb-e74c-49f3-9bd6-46b290fe0ef4")
	session4 := uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8")

	classGrades := []models.ClassGrade{
		{
			ID:           uuid.MustParse("0f6d9ab4-2b7b-41e8-b823-3ba45103a899"),
			Name:         "Junior Secondary School 1",
			Code:         "JSS1",
			Level:        7,
			Description:  "The first year of junior secondary education, focusing on building strong foundations in core subjects and essential learning skills.",
			AcademicSessionID:  session6,
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("23064952-17b6-4b8d-9889-b44e35e46faf"),
			Name:         "Junior Secondary School 2",
			Code:         "JSS2",
			Level:        8,
			Description:  "The second year of junior secondary education, where students strengthen their knowledge, develop deeper understanding, and improve their academic skills.",
			AcademicSessionID:  session6,
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("2dd8fe9c-3c2c-4668-bb6b-79fcd9e79bbe"),
			Name:         "Junior Secondary School 3",
			Code:         "JSS3",
			Level:        9,
			Description:  "The final year of junior secondary education, preparing students for senior secondary school through advanced learning and academic assessment.",
			AcademicSessionID:  session6,
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Name:         "Senior Secondary School 1",
			Code:         "SS1",
			Level:        10,
			Description:  "The first year of senior secondary education, introducing students to more specialized subjects and preparing them for advanced academic study.",
			AcademicSessionID:  session6,
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("3a61e8d2-4f79-48b3-bc25-7e9d1a6f5082"),
			Name:         "Senior Secondary School 2",
			Code:         "SS2",
			Level:        11,
			Description:  "The second year of senior secondary education, where students deepen their subject knowledge and develop skills needed for examinations and future careers.",
			AcademicSessionID:  session6,
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("c84b1597-2e63-4a8f-91d4-6b7c3e205fa9"),
			Name:         "Senior Secondary School 3",
			Code:         "SS3",
			Level:        12,
			Description:  "The final year of senior secondary education, focused on completing the secondary curriculum and preparing students for major examinations and further education.",
			AcademicSessionID:  session6,
			Capacity:     15,
			Status:       "active",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("7d2f6a83-91c4-45be-a738-5e1d9c4f206b"),
			Name:         "Junior Secondary School 1",
			Code:         "JSS1",
			Level:        7,
			Description:  "The first year of junior secondary education, focusing on building strong foundations in core subjects and essential learning skills.",
			AcademicSessionID:  session5,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("e5a93c71-8b26-4fd4-a159-3c7e2b905648"),
			Name:         "Junior Secondary School 2",
			Code:         "JSS2",
			Level:        8,
			Description:  "The second year of junior secondary education, where students strengthen their knowledge, develop deeper understanding, and improve their academic skills.",
			AcademicSessionID:  session5,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("1c7f4e92-a563-49bd-b281-6e3d8a5079c4"),
			Name:         "Junior Secondary School 3",
			Code:         "JSS3",
			Level:        9,
			Description:  "The final year of junior secondary education, preparing students for senior secondary school through advanced learning and academic assessment.",
			AcademicSessionID:  session5,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("b6d21f85-3a94-4ce7-97b2-8f1c5e6043ad"),
			Name:         "Senior Secondary School 1",
			Code:         "SS1",
			Level:        10,
			Description:  "The first year of senior secondary education, introducing students to more specialized subjects and preparing them for advanced academic study.",
			AcademicSessionID:  session5,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("4e8b3c16-d752-46a9-a1f4-9c6e2b5078d3"),
			Name:         "Senior Secondary School 2",
			Code:         "SS2",
			Level:        11,
			Description:  "The second year of senior secondary education, where students deepen their subject knowledge and develop skills needed for examinations and future careers.",
			AcademicSessionID:  session5,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("a2f95d63-7c41-4eb8-b592-1d8f6a3047ce"),
			Name:         "Senior Secondary School 3",
			Code:         "SS3",
			Level:        12,
			Description:  "The final year of senior secondary education, focused on completing the secondary curriculum and preparing students for major examinations and further education.",
			AcademicSessionID:  session5,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("6b3e9a27-c815-4f64-82d1-7a5c9e3046bf"),
			Name:         "Junior Secondary School 1",
			Code:         "JSS1",
			Level:        7,
			Description:  "The first year of junior secondary education, focusing on building strong foundations in core subjects and essential learning skills.",
			AcademicSessionID:  session4,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("f1c8a742-3e95-4d61-8b2a-7c5d91e4a630"),
			Name:         "Junior Secondary School 2",
			Code:         "JSS2",
			Level:        8,
			Description:  "The second year of junior secondary education, where students strengthen their knowledge, develop deeper understanding, and improve their academic skills.",
			AcademicSessionID:  session4,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("2b6d9f13-8a47-4cbe-91d5-e73a5f208164"),
			Name:         "Junior Secondary School 3",
			Code:         "JSS3",
			Level:        9,
			Description:  "The final year of junior secondary education, preparing students for senior secondary school through advanced learning and academic assessment.",
			AcademicSessionID:  session4,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("8e43c1b7-5d29-4fa6-a812-3c9e74b15d02"),
			Name:         "Senior Secondary School 1",
			Code:         "SS1",
			Level:        10,
			Description:  "The first year of senior secondary education, introducing students to more specialized subjects and preparing them for advanced academic study.",
			AcademicSessionID:  session4,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("4a9f2e61-c837-4bd4-92a1-6e5d3f7184bc"),
			Name:         "Senior Secondary School 2",
			Code:         "SS2",
			Level:        11,
			Description:  "The second year of senior secondary education, where students deepen their subject knowledge and develop skills needed for examinations and future careers.",
			AcademicSessionID:  session4,
			Capacity:     15,
			Status:       "completed",
			CreatedBy:    adminID,
		},
		{
			ID:           uuid.MustParse("d75b8c24-1f96-4ea3-b742-9c1e6a53d8f0"),
			Name:         "Senior Secondary School 3",
			Code:         "SS3",
			Level:        12,
			Description:  "The final year of senior secondary education, focused on completing the secondary curriculum and preparing students for major examinations and further education.",
			AcademicSessionID:  session4,
			Capacity:     15,
			Status:       "completed",
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
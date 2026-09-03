package seeds

import (
	"log"
	"time"
	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func SeedTests() error {
	db := config.GetDB()

	testDate := time.Date(2026, time.October, 24, 0, 0, 0, 0, time.UTC)
	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	tests := []models.Test{
		{
			ID:                uuid.MustParse("40d75889-aebf-43fb-909a-56c5adb48203"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          60,
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			TestDate:          &testDate,
			TestType:          "class_test",
			Title:             "Micro-organisms and Sexually Transmitted Infections Test",
			TotalMarks:        50,
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("d431c757-0bdb-4667-8750-f39837fc4b33"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          60,
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			TestDate:          &time.Time{},
			TestType:          "class_test",
			Title:             "Plant and Animal Health Test",
			TotalMarks:        50,
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("c23c32dd-68ed-48ed-b511-d146689fdddc"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          90,
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			TestDate:          &time.Time{},
			TestType:          "mid_term",
			Title:             "Mid-Term Biology and Agricultural Science Test",
			TotalMarks:        60,
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("fc282253-851e-438f-9e30-331d251643d0"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          60,
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			TestDate:          &time.Time{},
			TestType:          "class_test",
			Title:             "Food Production and Preservation Test",
			TotalMarks:        50,
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("16816378-c7f5-4be0-bebd-04918293a734"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          120,
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			TestDate:          &time.Time{},
			TestType:          "examination",
			Title:             "First Term Examination",
			TotalMarks:        100,
			CreatedBy:         adminID,
		},
	}

for _, test := range tests {
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"academic_session_id",
			"class_id",
			"duration",
			"status",
			"subject_id",
			"term_id",
			"test_date",
			"test_type",
			"title",
			"total_marks",
			"created_by",
		
		}),
	}).Create(&test)

	if result.Error != nil {
		log.Printf(
			"❌ Failed to seed test: %v",
			result.Error,
		)
	} else {
		log.Printf(
			"✅ Seeded/updated test: %s",
			test.Title,
		)
	}
}
return nil
}

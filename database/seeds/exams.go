package seeds

import (
	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
	"log"
	"time"
)

func SeedExams() {
	db := config.GetDB()

	examDate1 := time.Date(2026, time.November, 23, 0, 0, 0, 0, time.UTC)

	examDate3 := time.Date(2026, time.November, 25, 0, 0, 0, 0, time.UTC)

	examDate5 := time.Date(2026, time.November, 27, 0, 0, 0, 0, time.UTC)

	exams := []models.Exam{
		{
			ID:                uuid.MustParse("7427ab2d-614e-43b1-a702-a9cbcde51c75"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          120,
			ExamDate:          &examDate1,
			ExamType:          "theory",
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			Title:             "Biology First Term Examination",
			TotalMarks:        100,
			CreatedBy:         uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd"),
		},
		{
			ID:                uuid.MustParse("3de5fc79-2f0a-4f1f-be2e-2213b59ebac8"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          90,
			ExamDate:          &examDate3,
			ExamType:          "practical",
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			Title:             "Biology Practical Examination",
			TotalMarks:        50,
			CreatedBy:         uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd"),
		},
		{
			ID:                uuid.MustParse("13576954-2ac5-4145-8427-02c37fcebb95"),
			AcademicSessionID: uuid.MustParse("9b13e80c-18bc-4327-8efa-c96315a27337"),
			ClassID:           uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046"),
			Duration:          120,
			ExamDate:          &examDate5,
			ExamType:          "final",
			Status:            "scheduled",
			SubjectID:         uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			TermID:            uuid.MustParse("9dfc8c91-4df1-4b0b-8492-cd729a1115d8"),
			Title:             "First Term Final Examination",
			TotalMarks:        100,
			CreatedBy:         uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd"),
		},
	}

	for _, exam := range exams {
		if err := db.Create(&exam).Error; err != nil {
			log.Printf("❌ Failed to seed exam: %v", err)
		} else {
			log.Printf("✅ Seeded exam: %s", exam.Title)
		}
	}
}

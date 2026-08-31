package seeds

import (
	"fmt"
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedSchemesOfWork() error {
	db := config.GetDB()

	// Existing IDs
	academicSessionID := uuid.MustParse(
		"9b13e80c-18bc-4327-8efa-c96315a27337",
	)

	// Replace with your actual SS1 class ID
	classID := uuid.MustParse(
		"9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046",
	)

	// Replace with your actual Biology subject ID
	subjectID := uuid.MustParse(
		"a61583ba-fb14-4bd4-8054-80ea69da954d",
	)

	// Your First Term ID
	termID := uuid.MustParse(
		"9dfc8c91-4df1-4b0b-8492-cd729a1115d8",
	)
	adminID := uuid.MustParse(
		"fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd",
	)
	schemes := []models.SchemeOfWork{
		{
			ID:                uuid.MustParse("a13e7c92-5b46-4d81-9f20-6382e5741ab9"),
			AcademicSessionID: academicSessionID,
			ClassID:           classID,
			SubjectID:         subjectID,
			TermID:            termID,
			Title:             "Biology SS1 First Term Scheme of Work",
			Description:       "Scheme of work for Biology SS1 covering the major topics and learning activities for the first term.",
			Status:            "active",
			CreatedBy:         adminID,
		},
		{
			ID:                uuid.MustParse("b27f8d41-6c93-4ea5-a102-7593d6842fc1"),
			AcademicSessionID: academicSessionID,
			ClassID:           classID,
			SubjectID:         subjectID,
			TermID:            termID,
			Title:             "Biology SS1 First Term Practical Scheme",
			Description:       "Practical activities and experiments for Biology SS1 first term.",
			Status:            "active",
			CreatedBy:         adminID,
		},
	}

	for _, scheme := range schemes {
		if err := db.Create(&scheme).Error; err != nil {
			return fmt.Errorf(
				"failed to seed scheme of work %s: %w",
				scheme.Title,
				err,
			)
		}

		log.Printf("✅ Seeded scheme of work: %s", scheme.Title)
	}

	return nil
}

package seeds

import (
	"fmt"
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedModules() error {
	db := config.GetDB()
		subjectID := uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d")


	modules := []models.Module{
		{
			ID:          uuid.MustParse("2e7a1c45-8b93-4d16-a5f2-701c9e384b27"),
			Code:        "BIO-MICRO",
			Name:        "Micro-organisms",
			Description: "Study of microorganisms, their characteristics, importance and culturing.",
			Sequence:    1,
			SubjectID:   subjectID,
		},
		{
			ID:          uuid.MustParse("6b4f92d8-1e35-47ac-b609-83d72f514e96"),
			Code:        "BIO-STI",
			Name:        "Sexually Transmitted Infections",
			Description: "Study of sexually transmitted infections, their causes, effects, prevention and control.",
			Sequence:    2,
			SubjectID:   subjectID,
		},
		{
			ID:          uuid.MustParse("91c5e7a3-4d82-46bf-a019-72e8c5361db4"),
			Code:        "BIO-HEALTH",
			Name:        "Towards Better Health",
			Description: "Study of methods of maintaining good health and controlling diseases.",
			Sequence:    3,
			SubjectID:   subjectID,
		},
		{
			ID:          uuid.MustParse("3f8a62d1-7c45-4e90-b236-15d9a7840cef"),
			Code:        "BIO-AGRIC",
			Name:        "Biology and Agriculture",
			Description: "Study of the relevance and applications of biology in agriculture.",
			Sequence:    4,
			SubjectID:   subjectID,
		},
		{
			ID:          uuid.MustParse("a74e19c6-52d3-4b80-9f61-38c7d205ea94"),
			Code:        "BIO-PESTS",
			Name:        "Pests and Diseases",
			Description: "Study of pests and diseases affecting plants and animals.",
			Sequence:    5,
			SubjectID:   subjectID,
		},
		{
			ID:          uuid.MustParse("c8295f14-63ab-4e72-b508-91d6a3472fc0"),
			Code:        "BIO-FOOD",
			Name:        "Food Storage and Production",
			Description: "Study of methods of food production, preservation and storage.",
			Sequence:    6,
			SubjectID:   subjectID,
		},
	}

	for _, module := range modules {
		if err := db.Create(&module).Error; err != nil {
			return fmt.Errorf(
				"failed to seed module %s: %w",
				module.Name,
				err,
			)
		}

		log.Printf("✅ Seeded module: %s", module.Name)
	}

	return nil
}
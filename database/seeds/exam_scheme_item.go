package seeds

import (
	"log"
	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedExamSchemeItem() {
	db := config.GetDB()

	testSchemeItems := []models.ExamSchemeItem{
		{
			ExamID:                uuid.MustParse("7427ab2d-614e-43b1-a702-a9cbcde51c75"),
			SchemeOfWorkItemID: uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146"),
	
		},
		{
			ExamID:                uuid.MustParse("3de5fc79-2f0a-4f1f-be2e-2213b59ebac8"),
			SchemeOfWorkItemID: uuid.MustParse("7a2e9c51-4f83-46bd-b027-915c6a3e7481"),
	
		},
		{
			ExamID:               uuid.MustParse("13576954-2ac5-4145-8427-02c37fcebb95"),
			SchemeOfWorkItemID: uuid.MustParse("3e91a745-6b28-4c50-b839-172d5f6a904e"),
	
		},
	
		
	}

	for _, test := range testSchemeItems {
		if err := db.Create(&test).Error; err != nil {
			log.Printf("❌ Failed to seed test: %v", err)
		} else {
			log.Printf("✅ Seeded test scheme item")
		}
	}
}

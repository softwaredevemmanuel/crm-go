package seeds

import (
	"log"
	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedTestSchemeItem() {
	db := config.GetDB()

	testSchemeItems := []models.TestSchemeItem{
		{
			TestID:                uuid.MustParse("40d75889-aebf-43fb-909a-56c5adb48203"),
			SchemeOfWorkItemID: uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146"),
	
		},
		{
			TestID:                uuid.MustParse("d431c757-0bdb-4667-8750-f39837fc4b33"),
			SchemeOfWorkItemID: uuid.MustParse("7a2e9c51-4f83-46bd-b027-915c6a3e7481"),
	
		},
		{
			TestID:               uuid.MustParse("c23c32dd-68ed-48ed-b511-d146689fdddc"),
			SchemeOfWorkItemID: uuid.MustParse("3e91a745-6b28-4c50-b839-172d5f6a904e"),
	
		},
		{
			TestID:                uuid.MustParse("fc282253-851e-438f-9e30-331d251643d0"),
			SchemeOfWorkItemID: uuid.MustParse("f2c8d631-9a47-45be-a103-764e2b9158cd"),
	
		},
		{
			TestID:                uuid.MustParse("16816378-c7f5-4be0-bebd-04918293a734"),
			SchemeOfWorkItemID: uuid.MustParse("a9c3f815-2e64-4b71-8d30-547a1f96c2be"),
	
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

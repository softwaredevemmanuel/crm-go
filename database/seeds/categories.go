package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedCategories() {
	db := config.GetDB()
	categoryID1, err1 := uuid.Parse("283eaf1d-bf89-40fe-b363-b119e5815107")
	categoryID2, err2 := uuid.Parse("6109ae19-2a10-43b3-9daf-94c62f133724")
	categoryID3, err3 := uuid.Parse("84bd048a-f146-4191-99ea-5f45ebd551a3")
	categoryID4, err4 := uuid.Parse("d6b4bc52-c5c5-4af1-96f2-8daf6d977984")
	categoryID5, err5 := uuid.Parse("ec532f97-9de7-4fb4-8b93-dff24ca4acbd")

	if err1 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err5)
	}

	categories := []models.Category{
		{
		ID: 	  categoryID1,
		Name:     "Coding",
		},
		{
		ID: 	  categoryID2,
		Name:     "Design",
		},
		{
		ID: 	  categoryID3,
		Name:     "Backend Development",
		},
		{
		ID: 	  categoryID4,
		Name:     "Frontend Development",
		},
		{
		ID: 	  categoryID5,
		Name:     "Programming",
		},	
		
	}

	

	for _, category := range categories {
		if err := db.Create(&category).Error; err != nil {
			log.Printf("❌ Failed to seed Category: %v", err)
		} else {
			log.Printf("✅ Seeded Category: %v", &category.Name)
		}
	}
}


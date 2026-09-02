package seeds

import (
	"log"

	"github.com/google/uuid"
	"crm-go/models"
	"crm-go/config"

)

func SeedModules() {
	db := config.GetDB()

	schemeID := uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146")


	modules := []models.Module{
		{
			ID:             uuid.MustParse("f2c8d631-9a47-45be-a103-764e2b9158cd"),
			SchemeOfWorkID: schemeID,
			Title:          "Introduction to Biology",
			Description:    "Introduction to Biology, its meaning, branches, importance, and applications in everyday life.",
			ModuleOrder:    1,
		},
		{
			ID:             uuid.MustParse("5b7e4a29-83d1-4f60-b925-318c6e7a5042"),
			SchemeOfWorkID: schemeID,
			Title:          "The Cell",
			Description:    "Study of the cell as the basic structural and functional unit of living organisms.",
			ModuleOrder:    2,
		},
		{
			ID:             uuid.MustParse("a9c3f815-2e64-4b71-8d30-547a1f96c2be"),
			SchemeOfWorkID: schemeID,
			Title:          "Organization of Life",
			Description:    "Study of the levels of organization in living organisms from cells to organisms.",
			ModuleOrder:    3,
		},
		{
			ID:             uuid.MustParse("6d18b953-7c42-4ea0-a615-829f3b4c7061"),
			SchemeOfWorkID: schemeID,
			Title:          "Nutrition",
			Description:    "Study of nutrients, modes of nutrition, food substances, and the importance of balanced diets.",
			ModuleOrder:    4,
		},
		{
			ID:             uuid.MustParse("e5a7c240-1d83-4f96-b528-639e2a7154cd"),
			SchemeOfWorkID: schemeID,
			Title:          "Reproduction",
			Description:    "Introduction to reproduction in living organisms and the different methods of reproduction.",
			ModuleOrder:    5,
		},
		{
			ID:             uuid.MustParse("42f9b681-5c37-4ad2-8e10-753c6b9a214f"),
			SchemeOfWorkID: schemeID,
			Title:          "Ecology",
			Description:    "Introduction to ecology, ecosystems, environmental factors, and relationships between organisms.",
			ModuleOrder:    6,
		},
		{
			ID:             uuid.MustParse("9c6e2a57-4b91-46d3-a820-185f7c3946be"),
			SchemeOfWorkID: schemeID,
			Title:          "Classification of Living Organisms",
			Description:    "Study of the classification of living organisms and the characteristics used in classification.",
			ModuleOrder:    7,
		},
	}

	for _, module := range modules {
		if err := db.Create(&module).Error; err != nil {
			log.Printf("❌ Failed to seed module %s: %v", module.Title, err)
		} else {
			log.Printf("✅ Seeded Module: %s", module.Title)
		}
	}
}
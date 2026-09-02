package seeds

import (
	"log"

	"crm-go/config"
	"github.com/google/uuid"
	"crm-go/models"
)

func SeedTopics() {
	db := config.GetDB()
	moduleID1 := uuid.MustParse("f2c8d631-9a47-45be-a103-764e2b9158cd")
	moduleID2 := uuid.MustParse("5b7e4a29-83d1-4f60-b925-318c6e7a5042")
	moduleID3 := uuid.MustParse("a9c3f815-2e64-4b71-8d30-547a1f96c2be")
	moduleID4 := uuid.MustParse("6d18b953-7c42-4ea0-a615-829f3b4c7061")
	moduleID5 := uuid.MustParse("e5a7c240-1d83-4f96-b528-639e2a7154cd")
	moduleID6 := uuid.MustParse("42f9b681-5c37-4ad2-8e10-753c6b9a214f")
	moduleID7 := uuid.MustParse("9c6e2a57-4b91-46d3-a820-185f7c3946be")

	Topics := []models.Topic{
		{
			ID:          uuid.MustParse("5e7594e8-2b9a-48ac-9a83-565661971b3e"),
			ModuleID:    moduleID1,
			Title:       "Meaning of Biology",
			Description: "Introduction to Biology and the study of living organisms.",
			TopicOrder:  1,
		},
		{
			ID:          uuid.MustParse("825a9320-b960-4bbd-b884-54aef0dd6783"),
			ModuleID:    moduleID1,
			Title:       "Branches of Biology",
			Description: "Study of the major branches of Biology and their areas of specialization.",
			TopicOrder:  2,
		},
		{
			ID:          uuid.MustParse("5fde34d0-1382-41d8-bb6a-d757b33f320c"),
			ModuleID:    moduleID1,
			Title:       "Importance of Biology",
			Description: "The importance and applications of Biology in agriculture, medicine, environment, and everyday life.",
			TopicOrder:  3,
		},
		{
			ID:          uuid.MustParse("8e3f7f43-5373-4cc4-aaf3-98abed49668b"),
			ModuleID:    moduleID2,
			Title:       "Cell Theory",
			Description: "Study of the principles and development of the cell theory.",
			TopicOrder:  4,
		},
		{
			ID:          uuid.MustParse("45a2de63-5b9e-4405-9220-211d48421268"),
			ModuleID:    moduleID2,
			Title:       "Cell Structure",
			Description: "Study of the structures and components of plant and animal cells.",
			TopicOrder:  5,
		},
		{
			ID:          uuid.MustParse("d4fdb49d-2c35-45c4-a96a-1976149dc526"),
			ModuleID:    moduleID2,
			Title:       "Cell Organelles",
			Description: "Study of major cell organelles and their functions.",
			TopicOrder:  6,
		},
		{
			ID:          uuid.MustParse("ce9ac158-3c81-42b2-bee2-dc161b409724"),
			ModuleID:    moduleID2,
			Title:       "Plant and Animal Cells",
			Description: "Comparison of the structures and functions of plant and animal cells.",
			TopicOrder:  7,
		},
		{
			ID:          uuid.MustParse("8ba7b662-e85c-47fb-bfa4-e6a1808d1787"),
			ModuleID:    moduleID2,
			Title:       "Cell Division",
			Description: "Introduction to cell division and its importance in living organisms.",
			TopicOrder:  8,
		},
		{
			ID:          uuid.MustParse("1f9aa811-91a3-41c5-9dfc-9965e26147c4"),
			ModuleID:    moduleID3,
			Title:       "Levels of Organization",
			Description: "Study of the levels of organization from cells to tissues, organs, systems, and organisms.",
			TopicOrder:  9,
		},
		{
			ID:          uuid.MustParse("73b829cd-02bd-4527-9274-07d1307510ab"),
			ModuleID:    moduleID3,
			Title:       "Tissues",
			Description: "Study of tissues in plants and animals and their functions.",
			TopicOrder:  10,
		},
		{
			ID:          uuid.MustParse("d1abc70a-f053-4b51-aa82-ee60d6a787d3"),
			ModuleID:    moduleID3,
			Title:       "Organs and Organ Systems",
			Description: "Study of organs and organ systems and how they work together.",
			TopicOrder:  11,
		},
		{
			ID:          uuid.MustParse("d293e26c-2aa3-4bab-8a6e-e7497161d55e"),
			ModuleID:    moduleID4,
			Title:       "Classes of Food",
			Description: "Study of carbohydrates, proteins, fats, vitamins, minerals, water, and their functions.",
			TopicOrder:  12,
		},
		{
			ID:          uuid.MustParse("8632ca6f-79f5-456d-8cc1-d73cff707d0c"),
			ModuleID:    moduleID4,
			Title:       "Food Tests",
			Description: "Practical tests for identifying major food substances.",
			TopicOrder:  13,
		},
		{
			ID:          uuid.MustParse("9262a8f7-b0a7-409b-9b4f-c5e0a8dcd285"),
			ModuleID:    moduleID4,
			Title:       "Balanced Diet",
			Description: "Study of balanced diets and the importance of adequate nutrition.",
			TopicOrder:  14,
		},
		{
			ID:          uuid.MustParse("caf5e77e-ead6-468a-853f-27b5133e9fbe"),
			ModuleID:    moduleID4,
			Title:       "Modes of Nutrition",
			Description: "Study of autotrophic and heterotrophic modes of nutrition.",
			TopicOrder:  15,
		},
		{
			ID:          uuid.MustParse("daf2c83a-91c4-4b01-88a6-18f03126174b"),
			ModuleID:    moduleID5,
			Title:       "Meaning of Reproduction",
			Description: "Introduction to reproduction and its importance for the continuity of life.",
			TopicOrder:  16,
		},
		{
			ID:          uuid.MustParse("6820228b-1c40-434d-83d2-4e3fb862af48"),
			ModuleID:    moduleID5,
			Title:       "Asexual Reproduction",
			Description: "Study of reproduction involving a single parent and its different forms.",
			TopicOrder:  17,
		},
		{
			ID:          uuid.MustParse("e978a030-0166-4d0a-889c-a604575dd8ed"),
			ModuleID:    moduleID5,
			Title:       "Sexual Reproduction",
			Description: "Study of sexual reproduction and the involvement of male and female gametes.",
			TopicOrder:  18,
		},
		{
			ID:          uuid.MustParse("ccc69310-7aad-4303-8490-8527feee0fa7"),
			ModuleID:    moduleID5,
			Title:       "Reproduction in Flowering Plants",
			Description: "Study of the reproductive structures and processes in flowering plants.",
			TopicOrder:  19,
		},
		{
			ID:          uuid.MustParse("53f9dccf-54d8-4174-b05e-3393e7cfd076"),
			ModuleID:    moduleID6,
			Title:       "Meaning of Ecology",
			Description: "Introduction to ecology and the study of relationships between organisms and their environment.",
			TopicOrder:  20,
		},
		{
			ID:          uuid.MustParse("ffae1027-7b9c-4778-b336-4b3cbbe7064e"),
			ModuleID:    moduleID6,
			Title:       "Ecosystem",
			Description: "Study of ecosystems and their biotic and abiotic components.",
			TopicOrder:  21,
		},
		{
			ID:          uuid.MustParse("620a5bae-3087-4082-9bb5-a7334ebedf9d"),
			ModuleID:    moduleID6,
			Title:       "Food Chains and Food Webs",
			Description: "Study of feeding relationships among organisms in an ecosystem.",
			TopicOrder:  22,
		},
		{
			ID:          uuid.MustParse("0e316b17-89a9-4929-8b4a-7017630c0ef7"),
			ModuleID:    moduleID6,
			Title:       "Environmental Factors",
			Description: "Study of biotic and abiotic factors affecting living organisms.",
			TopicOrder:  23,
		},
		{
			ID:          uuid.MustParse("17a486e0-91a6-404a-9705-b595a3cb6fc5"),
			ModuleID:    moduleID7,
			Title:       "Meaning of Classification",
			Description: "Introduction to the classification of living organisms.",
			TopicOrder:  24,
		},
		{
			ID:          uuid.MustParse("af3a3497-a4b1-4db0-a918-5e31d231270a"),
			ModuleID:    moduleID7,
			Title:       "Importance of Classification",
			Description: "Study of the importance of classifying living organisms.",
			TopicOrder:  25,
		},
		{
			ID:          uuid.MustParse("053e2dea-c4a7-402c-beba-69d9ed1a8e43"),
			ModuleID:    moduleID7,
			Title:       "Kingdoms of Living Organisms",
			Description: "Introduction to the major kingdoms used in biological classification.",
			TopicOrder:  26,
		},
		{
			ID:          uuid.MustParse("8e58b2e2-bd1c-402f-8504-f81b8575b156"),
			ModuleID:    moduleID7,
			Title:       "Binomial Nomenclature",
			Description: "Study of the system used to give organisms scientific names.",
			TopicOrder:  27,
		},
	}

	for _, topic := range Topics {
		if err := db.Create(&topic).Error; err != nil {
			log.Printf("❌ Failed to seed topic %s: %v", topic.Title, err)
		} else {
			log.Printf("✅ Seeded Topic: %s", topic.Title)
		}
	}
}

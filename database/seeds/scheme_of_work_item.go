package seeds

import (
	"fmt"
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedSchemeOfWorkItems() error {
	db := config.GetDB()

	// Biology SS1 First Term Scheme of Work
	schemeOfWorkID := uuid.MustParse(
		"a13e7c92-5b46-4d81-9f20-6382e5741ab9",
	)

	// Biology Modules
	microOrganismsModuleID := uuid.MustParse(
		"2e7a1c45-8b93-4d16-a5f2-701c9e384b27",
	)

	stiModuleID := uuid.MustParse(
		"6b4f92d8-1e35-47ac-b609-83d72f514e96",
	)

	healthModuleID := uuid.MustParse(
		"91c5e7a3-4d82-46bf-a019-72e8c5361db4",
	)


	agricultureModuleID := uuid.MustParse(
		"3f8a62d1-7c45-4e90-b236-15d9a7840cef",
	)

	pestsModuleID := uuid.MustParse(
		"a74e19c6-52d3-4b80-9f61-38c7d205ea94",
	)

	foodModuleID := uuid.MustParse(
		"c8295f14-63ab-4e72-b508-91d6a3472fc0",
	)

	items := []models.SchemeOfWorkItem{

		// WEEK 1
		{
			ID:             uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       microOrganismsModuleID,
			Sequence:       1,
			WeekStart:      1,
			WeekEnd:        1,
			Topic:          "Micro-organisms Around Us",
			Subtopic:       "Concept of Micro-organisms",
			Content: `Introduction to micro-organisms. Definition and characteristics of microorganisms.
			Distribution of microorganisms in air, water, soil and living organisms.
			Beneficial and harmful microorganisms.`,
		},

			{
			ID:             uuid.MustParse("c4e7b912-3a56-4d80-b135-628f9a1e4732"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       microOrganismsModuleID,
			Sequence:       2,
			WeekStart:      2,
			WeekEnd:        2,
			Topic:          "Micro-organisms Around Us",
			Subtopic:       "Culturing and Identification of Micro-organisms",
			Content: `Concept of culturing microorganisms.
			Culture media, pure and mixed cultures, agar and laboratory equipment.
			Methods of identifying microorganisms and observation of microbial colonies.`,
		},
		// WEEK 2
		{
			ID:             uuid.MustParse("1d8b6f43-9275-4ca1-a850-316e7d5942bc"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       stiModuleID,
			Sequence:       3,
			WeekStart:      3,
			WeekEnd:        3,
			Topic:          "Sexually Transmitted Infections",
			Subtopic:       "Causes, Effects and Prevention of STIs",
			Content: `Meaning of sexually transmitted infections.
			Common examples, modes of transmission, symptoms, effects,
			prevention and control of sexually transmitted infections.`,
		},
	

		// WEEK 3
		{
			ID:             uuid.MustParse("7a2e9c51-4f83-46bd-b027-915c6a3e7481"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       healthModuleID,
			Sequence:       4,
			WeekStart:      4,
			WeekEnd:        4,
			Topic:          "Micro-organisms in Action",
			Subtopic:       "Bacteria, Viruses, Protozoa and Algae",
			Content: `Study of bacteria, viruses, protozoa and algae.
			Structure, classification, reproduction, nutrition and importance of microorganisms.
			Beneficial and harmful effects of microorganisms.`,
		},

		// WEEK 4
		{
			ID:             uuid.MustParse("b6d4f820-1c59-4e73-9a26-583b7d14e9c0"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       healthModuleID,
			Sequence:       5,
			WeekStart:      5,
			WeekEnd:        5,
			Topic:          "Sexually Transmitted Infections",
			Subtopic:       "Causes, Effects and Prevention of STIs",
			Content: `Meaning of sexually transmitted infections.
			Common examples, modes of transmission, symptoms, effects, prevention and control.
			Importance of personal hygiene, responsible behaviour and medical treatment.`,
		},

		// WEEK 5
		{
			ID:             uuid.MustParse("3e91a745-6b28-4c50-b839-172d5f6a904e"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       agricultureModuleID,
			Sequence:       6,
			WeekStart:      6,
			WeekEnd:        6,
			Topic:          "Towards Better Health",
			Subtopic:       "Control of Harmful Micro-organisms",
			Content: `Methods of controlling harmful microorganisms.
			Use of antibiotics, antiseptics and disinfectants.
			Effects of temperature, salting, dehydration, sanitation, immunization and balanced diet.`,
		},

		// WEEK 6
		{
			ID:             uuid.MustParse("f2c8d631-9a47-45be-a103-764e2b9158cd"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       agricultureModuleID,
			Sequence:       7,
			WeekStart:      7,
			WeekEnd:        7,
			Topic:          "Relevance of Biology to Agriculture",
			Subtopic:       "Classification of Plants",
			Content: `Importance of biology in agriculture.
			Classification of plants based on botanical features, agricultural importance,
			life cycle, size, growth habit, leaf characteristics, temperature and habitat.`,
		},

		// WEEK 7
		{
			ID:             uuid.MustParse("5b7e4a29-83d1-4f60-b925-318c6e7a5042"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       pestsModuleID,
			Sequence:       8,
			WeekStart:      8,
			WeekEnd:        8,
			Topic:          "Pests and Diseases of Plants",
			Subtopic:       "Plant Pests and Diseases",
			Content: `Common pests of agricultural plants.
			Types of plant diseases, causes, symptoms, effects and methods of prevention and control.`,
		},

		// WEEK 8
		{
			ID:             uuid.MustParse("a9c3f815-2e64-4b71-8d30-547a1f96c2be"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       pestsModuleID,
			Sequence:       9,
			WeekStart:      9,
			WeekEnd:        9,
			Topic:          "Pests and Diseases of Animals",
			Subtopic:       "Animal Pests and Diseases",
			Content: `Common pests and diseases affecting farm animals.
			Causes, symptoms, effects, prevention and control of animal diseases.`,
		},

		// WEEK 9
		{
			ID:             uuid.MustParse("6d18b953-7c42-4ea0-a615-829f3b4c7061"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       pestsModuleID,
			Sequence:       10,
			WeekStart:      10,
			WeekEnd:        10,
			Topic:          "Pests and Diseases of Animals",
			Subtopic:       "Control of Animal Diseases",
			Content: `Methods of preventing and controlling animal diseases.
			Importance of hygiene, vaccination, quarantine and proper animal husbandry.`,
		},

		// WEEK 10
		{
			ID:             uuid.MustParse("e5a7c240-1d83-4f96-b528-639e2a7154cd"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       foodModuleID,
			Sequence:       11,
			WeekStart:      11,
			WeekEnd:        11,
			Topic:          "Food Storage and Production",
			Subtopic:       "Food Production",
			Content: `Methods of food production.
			Importance of agriculture and biological principles in food production.
			Factors affecting food production.`,
		},

		// WEEK 11
		{
			ID:             uuid.MustParse("42f9b681-5c37-4ad2-8e10-753c6b9a214f"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       foodModuleID,
			Sequence:       12,
			WeekStart:      12,
			WeekEnd:        12,
			Topic:          "Food Storage and Production",
			Subtopic:       "Food Preservation and Storage",
			Content: `Methods of preserving and storing food.
			Drying, salting, smoking, refrigeration, freezing, canning and other methods.
			Importance of food preservation in preventing spoilage.`,
		},

		// WEEK 12
		{
			ID:             uuid.MustParse("9c6e2a57-4b91-46d3-a820-185f7c3946be"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:       foodModuleID,
			Sequence:       13,
			WeekStart:      13,
			WeekEnd:        13,
			Topic:          "Revision",
			Subtopic:       "Revision of Term's Work",
			Content: `Revision of all topics covered during the term.
			Review of microorganisms, health, agriculture, pests, diseases,
			food production, preservation and storage.`,
		},
	}

	for _, item := range items {
		if err := db.Create(&item).Error; err != nil {
			return fmt.Errorf(
				"failed to seed scheme of work item %s: %w",
				item.Topic,
				err,
			)
		}

		log.Printf(
			"✅ Seeded scheme item: Week %d - %s",
			item.WeekStart,
			item.Topic,
		)
	}

	return nil
}

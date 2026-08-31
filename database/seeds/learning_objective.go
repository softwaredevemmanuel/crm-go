package seeds

import (
	"fmt"
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedLearningObjectives() error {
	db := config.GetDB()

	schemeOfWorkItemID1:= uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146")
	schemeOfWorkItemID2:= uuid.MustParse("c4e7b912-3a56-4d80-b135-628f9a1e4732")
	schemeOfWorkItemID3:= uuid.MustParse("1d8b6f43-9275-4ca1-a850-316e7d5942bc")
	schemeOfWorkItemID4:= uuid.MustParse("7a2e9c51-4f83-46bd-b027-915c6a3e7481")
	schemeOfWorkItemID5:= uuid.MustParse("b6d4f820-1c59-4e73-9a26-583b7d14e9c0")
	schemeOfWorkItemID6:= uuid.MustParse("3e91a745-6b28-4c50-b839-172d5f6a904e")
	schemeOfWorkItemID7:= uuid.MustParse("f2c8d631-9a47-45be-a103-764e2b9158cd")
	schemeOfWorkItemID8:= uuid.MustParse("5b7e4a29-83d1-4f60-b925-318c6e7a5042")
	schemeOfWorkItemID9:= uuid.MustParse("a9c3f815-2e64-4b71-8d30-547a1f96c2be")
	schemeOfWorkItemID10:= uuid.MustParse("6d18b953-7c42-4ea0-a615-829f3b4c7061")
	schemeOfWorkItemID11:= uuid.MustParse("e5a7c240-1d83-4f96-b528-639e2a7154cd")
	schemeOfWorkItemID12:= uuid.MustParse("42f9b681-5c37-4ad2-8e10-753c6b9a214f")
	
	objectives := []models.LearningObjective{

		// =========================================================
		// WEEK 1 - Micro-organisms Around Us
		// Scheme Item:
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-1001-4a01-8101-100000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID1,
			Objective:          "Define micro-organisms and state their general characteristics.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-1002-4a02-8102-100000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID1,
			Objective:          "Identify places where micro-organisms are commonly found.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-1003-4a03-8103-100000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID1,
			Objective:          "Differentiate between beneficial and harmful micro-organisms.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 2 - Culturing and Identification
		// Scheme Item:
		// a1b2c3d4-2222-4a22-8222-222222222222
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-2001-4b01-8201-200000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID2,
			Objective:          "Explain the concept of culturing micro-organisms.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-2002-4b02-8202-200000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID2,
			Objective:          "Describe the purpose of culture media in growing micro-organisms.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-2003-4b03-8203-200000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID2,
			Objective:          "Identify basic laboratory equipment used for culturing micro-organisms.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 3 - Sexually Transmitted Infections
		// Scheme Item:
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-3001-4c01-8301-300000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID3,
			Objective:          "Define sexually transmitted infections and give common examples.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-3002-4c02-8302-300000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID3,
			Objective:          "Describe the major modes of transmission of sexually transmitted infections.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-3003-4c03-8303-300000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID3,
			Objective:          "Explain methods of preventing and controlling sexually transmitted infections.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 4 - Micro-organisms in Action
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-4001-4d01-8401-400000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID4,
			Objective:          "Identify the major groups of micro-organisms including bacteria, viruses, protozoa and algae.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-4002-4d02-8402-400000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID4,
			Objective:          "Describe the basic structure and characteristics of major groups of micro-organisms.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-4003-4d03-8403-400000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID4,
			Objective:          "Explain the beneficial and harmful effects of micro-organisms.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 5 - Causes, Effects and Prevention of STIs
		// Scheme Item:
		// a1b2c3d4-4444-4a44-8444-444444444444
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-5001-4e01-8501-500000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID5,
			Objective:          "Identify common sexually transmitted infections and their causative organisms.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-5002-4e02-8502-500000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID5,
			Objective:          "Describe the symptoms and effects of common sexually transmitted infections.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-5003-4e03-8503-500000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID5,
			Objective:          "Explain the importance of personal hygiene, responsible behaviour and medical treatment.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 6 - Control of Harmful Micro-organisms
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-6001-4f01-8601-600000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID6,
			Objective:          "Explain different methods used to control harmful micro-organisms.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-6002-4f02-8602-600000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID6,
			Objective:          "Differentiate between antibiotics, antiseptics and disinfectants.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-6003-4f03-8603-600000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID6,
			Objective:          "Explain how sanitation, immunization and proper nutrition help control disease.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 7 - Classification of Plants
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-7001-4a11-8701-700000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID7,
			Objective:          "Explain the relevance of biology to agriculture.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-7002-4a12-8702-700000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID7,
			Objective:          "Classify plants based on their botanical and agricultural characteristics.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-7003-4a13-8703-700000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID7,
			Objective:          "Classify plants according to life cycle, growth habit and habitat.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 8 - Plant Pests and Diseases
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-8001-4b11-8801-800000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID7,
			Objective:          "Identify common pests that affect agricultural plants.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-8002-4b12-8802-800000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID7,
			Objective:          "Describe common plant diseases, their causes and symptoms.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-8003-4b13-8803-800000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID7,
			Objective:          "Explain methods of preventing and controlling plant pests and diseases.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 9 - Animal Pests and Diseases
		// Scheme Item:
		// a1b2c3d4-8888-4a88-8888-888888888888
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-9001-4c11-8901-900000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID8,
			Objective:          "Identify common pests and diseases affecting farm animals.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-9002-4c12-8902-900000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID8,
			Objective:          "Describe the causes and symptoms of common animal diseases.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-9003-4c13-8903-900000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID8,
			Objective:          "Explain the effects of animal pests and diseases on agricultural production.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 10 - Control of Animal Diseases
		// =========================================================

		{
			ID:                 uuid.MustParse("c1d2e3f4-1001-4d11-8101-010000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID9,
			Objective:          "Explain methods used to prevent and control animal diseases.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-1002-4d12-8102-010000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID9,
			Objective:          "Explain the importance of hygiene and vaccination in animal health.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("c1d2e3f4-1003-4d13-8103-010000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID9,
			Objective:          "Explain the importance of quarantine and proper animal husbandry.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 11 - Food Production
		// =========================================================

		{
			ID:                 uuid.MustParse("d1e2f3a4-1101-4e11-8101-110000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID10,
			Objective:          "Describe different methods of food production.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("d1e2f3a4-1102-4e12-8102-110000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID10,
			Objective:          "Explain the importance of agriculture in food production.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("d1e2f3a4-1103-4e13-8103-110000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID10,
			Objective:          "Identify factors that affect food production.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 12 - Food Preservation and Storage
		// =========================================================

		{
			ID:                 uuid.MustParse("d1e2f3a4-1201-4f11-8201-120000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID11,
			Objective:          "Explain the importance of food preservation and storage.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("d1e2f3a4-1202-4f12-8202-120000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID11,
			Objective:          "Describe methods of food preservation including drying, salting, smoking and refrigeration.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("d1e2f3a4-1203-4f13-8203-120000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID11,
			Objective:          "Explain how proper food storage prevents spoilage and food wastage.",
			Sequence:           3,
		},

		// =========================================================
		// WEEK 13 - Revision
		// =========================================================

		{
			ID:                 uuid.MustParse("d1e2f3a4-1301-4a21-8301-130000000001"),
			SchemeOfWorkItemID: schemeOfWorkItemID12,
			Objective:          "Review the major concepts covered during the term.",
			Sequence:           1,
		},
		{
			ID:                 uuid.MustParse("d1e2f3a4-1302-4a22-8302-130000000002"),
			SchemeOfWorkItemID: schemeOfWorkItemID12,
			Objective:          "Recall key concepts relating to microorganisms, health, agriculture, pests and diseases.",
			Sequence:           2,
		},
		{
			ID:                 uuid.MustParse("d1e2f3a4-1303-4a23-8303-130000000003"),
			SchemeOfWorkItemID: schemeOfWorkItemID12,
			Objective:          "Apply knowledge from the term's work to answer revision questions.",
			Sequence:           3,
		},
	}

	for _, objective := range objectives {
		if err := db.Create(&objective).Error; err != nil {
			return fmt.Errorf(
				"failed to seed learning objective %s: %w",
				objective.Objective,
				err,
			)
		}

		log.Printf(
			"✅ Seeded learning objective: %s",
			objective.Objective,
		)
	}

	return nil
}
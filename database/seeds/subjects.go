package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedSubjects() {
	db := config.GetDB()

	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")
	LanguagesDepartment := uuid.MustParse("4a547934-2c83-4705-8616-e1e817b66bb2")
	MathematicsDepartment := uuid.MustParse("78a57227-d318-4264-b373-c0ae96f7e83b")
	SciencesDepartment := uuid.MustParse("842dfd3f-6344-4c80-89a9-72ce9bb3ab15")
	SocialSciencesDepartment := uuid.MustParse("8bc3ebc8-79b7-4092-931f-d235e3bd7d1f")
	ArtsDepartment := uuid.MustParse("94d593bf-3ef4-4c10-9a82-ac32603a4350")
	TechnologyDepartment := uuid.MustParse("a9dddddb-28b8-4264-8bef-675fe5cc4e5c")
	BusinessStudiesDepartment := uuid.MustParse("b323135b-b1bc-4d71-9fa0-3d6affd8ad8b")
	AgriculturalSciencesDepartment := uuid.MustParse("d6ec9227-b333-4c58-9829-060005320746")
	HomeEconomicsDepartment := uuid.MustParse("eaa85465-26f5-430d-a4d8-71150b1352e7")
	ReligiousStudiesDepartment := uuid.MustParse("f7329d7c-1fce-4737-b117-7976729877de")
	PhysicalEducationDepartment := uuid.MustParse("ca1d4208-32d5-4dbf-9a5d-fbc8089dc341")

	subjects := []models.Subject{
		{
			ID:          uuid.MustParse("018d9b13-9aa5-4687-b91b-5a6db922b187"),
			Name:        "English Language",
			Code:        "ENG",
			Description: "Develops students' reading, writing, speaking, listening, and communication skills.",
			DepartmentID:  LanguagesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("031b69a7-afd5-4899-85d4-68fba93bc7a1"),
			Name:        "Mathematics",
			Code:        "MTH",
			Description: "Develops logical reasoning, problem-solving, numerical, and analytical skills.",
			DepartmentID:  MathematicsDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("0b1ea262-d457-444f-bf5d-16b36ce530ed"),
			Name:        "Social Studies",
			Code:        "SST",
			Description: "Explores society, culture, relationships, citizenship, and responsible social behaviour.",
			DepartmentID:  SocialSciencesDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("12a59d94-d8e1-4d48-acea-888d1eaac29a"),
			Name:        "Basic Science",
			Code:        "BSC",
			Description: "Introduces students to fundamental concepts in biology, chemistry, physics, and the environment.",
			DepartmentID:  SciencesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("17f496d4-1789-4bc9-836c-bf4b29d9a8a5"),
			Name:        "Basic Technology",
			Code:        "BTE",
			Description: "Introduces students to technology, technical drawing, materials, tools, and practical skills.",
			DepartmentID:  TechnologyDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("1e85ba4a-50d9-42f6-8c7a-7797a0514dfc"),
			Name:        "Civic Education",
			Code:        "CIV",
			Description: "Teaches citizenship, rights, responsibilities, democracy, and national values.",
			DepartmentID:  SocialSciencesDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("30344e8c-4396-47a6-a4ba-fef3be7d5df5"),
			Name:        "Computer Studies",
			Code:        "CST",
			Description: "Introduces students to computers, digital literacy, applications, and information technology.",
			DepartmentID:  TechnologyDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("45886b22-627c-4659-a697-803917de2872"),
			Name:        "Agricultural Science",
			Code:        "AGR",
			Description: "Introduces agricultural practices, farming, livestock, crops, and agricultural resources.",
			DepartmentID:  AgriculturalSciencesDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("4913623c-e696-4a81-9eee-7d4dda9c724c"),
			Name:        "Christian Religious Studies",
			Code:        "CRS",
			Description: "Studies Christian beliefs, values, teachings, morality, and biblical principles.",
			DepartmentID:  ReligiousStudiesDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("495aea99-2318-4e61-a911-7710347a86ca"),
			Name:        "Islamic Religious Studies",
			Code:        "IRS",
			Description: "Studies Islamic beliefs, practices, values, history, and moral principles.",
			DepartmentID:  ReligiousStudiesDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("4dfa28f7-181a-4ee7-812d-c839ef6ed181"),
			Name:        "Business Studies",
			Code:        "BST",
			Description: "Introduces students to commerce, office practice, accounting, and entrepreneurship.",
			DepartmentID:  BusinessStudiesDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("4fa32f2c-8a64-416b-9e81-73ff1bd84f76"),
			Name:        "Home Economics",
			Code:        "HEC",
			Description: "Develops practical knowledge of food, clothing, family living, and home management.",
			DepartmentID:  HomeEconomicsDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("51985d2a-562a-4def-9d8c-e201f2aae249"),
			Name:        "Cultural & Creative Arts",
			Code:        "CCA",
			Description: "Develops creativity through visual arts, crafts, music, drama, and cultural expression.",
			DepartmentID:  ArtsDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("73111766-effc-4c53-9130-b81e857d386d"),
			Name:        "Physical & Health Education",
			Code:        "PHE",
			Description: "Promotes physical fitness, healthy living, sports, and personal wellbeing.",
			DepartmentID:  PhysicalEducationDepartment,
			Credits:     1,
			Status:      "inactive",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("88c83da0-1845-4d1e-a534-455cbaff102e"),
			Name:        "French",
			Code:        "FRE",
			Description: "Develops basic communication skills and cultural awareness in the French language.",
			DepartmentID:  LanguagesDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("944c6ea6-b13c-434f-8501-ad4ba27eaff6"),
			Name:        "Literature in English",
			Code:        "LIT",
			Description: "Develops appreciation and understanding of poetry, prose, drama, and literary expression.",
			DepartmentID:  LanguagesDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("9c8be49c-058c-4b67-b486-fa2edd7f2faf"),
			Name:        "Geography",
			Code:        "GEO",
			Description: "Studies the physical environment, people, resources, climate, and geographical processes..",
			DepartmentID:  SocialSciencesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("9df5f176-89d4-43a4-87c2-f0a5cda20fd6"),
			Name:        "Economics",
			Code:        "ECO",
			Description: "Introduces the study of resources, production, consumption, markets, and economic decision-making.",
			DepartmentID:  SocialSciencesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("a544ba6f-4c81-470f-bd8e-b4af872df2b0"),
			Name:        "Government",
			Code:        "GOV",
			Description: "Examines political systems, governance, citizenship, institutions, and public administration.",
			DepartmentID:  SocialSciencesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("a61583ba-fb14-4bd4-8054-80ea69da954d"),
			Name:        "Biology",
			Code:        "BIO",
			Description: "Studies living organisms, their structures, functions, environments, and interactions.",
			DepartmentID:  SciencesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("a6612241-a522-472c-96a8-dce231d6657a"),
			Name:        "Chemistry",
			Code:        "CHE",
			Description: "Studies matter, chemical reactions, elements, compounds, and laboratory processes.",
			DepartmentID:  SciencesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("ab036e8b-fef4-48e1-be0c-f64fe810baa4"),
			Name:        "Physics",
			Code:        "PHY",
			Description: "Studies matter, energy, forces, motion, electricity, waves, and physical phenomena.",
			DepartmentID:  SciencesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("c7d39dbd-2910-43bf-aedb-2a742c995241"),
			Name:        "Further Mathematics",
			Code:        "FMT",
			Description: "Provides advanced mathematical concepts, reasoning, and problem-solving skills.",
			DepartmentID:  MathematicsDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		},
		{
			ID:          uuid.MustParse("cc0adc77-519d-4294-8a12-a0f527f779e6"),
			Name:        "Technical Drawing",
			Code:        "TDR",
			Description: "Develops skills in technical graphics, design, drafting, and engineering representation.",
			DepartmentID:  TechnologyDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("d17fa857-8338-4082-8458-71af375f3536"),
			Name:        "Data Processing",
			Code:        "DAP",
			Description: "Develops skills in computer applications, data management, programming concepts, and digital systems.",
			DepartmentID:  TechnologyDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("d47ce2cd-71c0-4ea2-b6eb-b19c7070a911"),
			Name:        "Financial Accounting",
			Code:        "FAC",
			Description: "Teaches accounting principles, financial records, transactions, and preparation of accounts.",
			DepartmentID:  BusinessStudiesDepartment,
			Credits:     3,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("ddd5f0d6-03de-45cb-b4f3-53a0cbc57c4e"),
			Name:        "Commerce",
			Code:        "COM",
			Description: "Studies trade, business activities, distribution, markets, and commercial practices.",
			DepartmentID:  BusinessStudiesDepartment,
			Credits:     2,
			Status:      "active",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("e4830fe2-0b8e-4b84-966c-84bea410de4a"),
			Name:        "Marketing",
			Code:        "MKT",
			Description: "Introduces principles of promoting, pricing, distributing, and selling products and services.",
			DepartmentID:  BusinessStudiesDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("e6217a27-7199-4486-92c7-62d3361f9755"),
			Name:        "Food & Nutrition",
			Code:        "FNT",
			Description: "Covers nutrition, food preparation, meal planning, health, and food management.",
			DepartmentID:  HomeEconomicsDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("f8ca986d-6fb8-4a17-9c02-7825da6a3bf0"),
			Name:        "Visual Art",
			Code:        "VAT",
			Description: "Develops skills and knowledge in drawing, painting, design, sculpture, and visual creativity.",
			DepartmentID:  ArtsDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("f9339467-117e-4782-af0c-19fbb6bd0f8f"),
			Name:        "Music",
			Code:        "MUS",
			Description: "Develops knowledge and skills in musical theory, performance, composition, and appreciation.",
			DepartmentID:  ArtsDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		}, {
			ID:          uuid.MustParse("fd42da47-c663-4887-83ff-11e5cc63a5a7"),
			Name:        "Drama & Theatre Arts",
			Code:        "DTA",
			Description: "Develops performance, acting, storytelling, stagecraft, and appreciation of theatre.",
			DepartmentID:  ArtsDepartment,
			Credits:     2,
			Status:      "inactive",
			CreatedBy:   adminID,
		},
	}

	for _, subject := range subjects {
		if err := db.Create(&subject).Error; err != nil {
			log.Printf(
				"❌ Failed to seed class grade %s: %v",
				subject.Name,
				err,
			)
			continue
		}

		log.Printf(
			"✅ Seeded class grade: %s",
			subject.Name,
		)
	}
}

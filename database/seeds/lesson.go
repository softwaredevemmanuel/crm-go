package seeds

import (
	"encoding/json"
	"log"
	"gorm.io/datatypes"
	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedLessons() {
	db := config.GetDB()

	// Scheme of work id
	schemeOfWorkID := uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146")

	// Module IDs from SeedModules
	moduleID1 := uuid.MustParse("f2c8d631-9a47-45be-a103-764e2b9158cd")
	moduleID2 := uuid.MustParse("5b7e4a29-83d1-4f60-b925-318c6e7a5042")
	moduleID3 := uuid.MustParse("a9c3f815-2e64-4b71-8d30-547a1f96c2be")
	moduleID4 := uuid.MustParse("6d18b953-7c42-4ea0-a615-829f3b4c7061")
	moduleID5 := uuid.MustParse("e5a7c240-1d83-4f96-b528-639e2a7154cd")
	moduleID6 := uuid.MustParse("42f9b681-5c37-4ad2-8e10-753c6b9a214f")
	moduleID7 := uuid.MustParse("9c6e2a57-4b91-46d3-a820-185f7c3946be")

	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	// Helper function to convert map to datatypes.JSON
	toJSON := func(data map[string]interface{}) datatypes.JSON {
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			return datatypes.JSON([]byte("{}"))
		}
		return datatypes.JSON(jsonData)
	}

	lessons := []models.Lesson{

		// =========================================================
		// MODULE 1: INTRODUCTION TO BIOLOGY
		// =========================================================

		{
			ID:          uuid.MustParse("efcc2c97-f489-4de1-a7c1-7a7c691a8ff4"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID1,
			TopicID:     uuid.MustParse("5e7594e8-2b9a-48ac-9a83-565661971b3e"),
			LessonOrder: 1,
			Week:        1,
			Title:       "Meaning of Biology",
			Description: "Introduction to Biology and the study of living organisms.",
			Objectives:  "Students should be able to define Biology, explain what Biology studies, and identify examples of living organisms.",
			Activities:  "Teacher explains the meaning of Biology. Students identify living and non-living things around them.",
			Resources:   "Biology textbook, charts, pictures of living organisms.",
			Assessment:  "Define Biology and list five examples of living organisms.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"introduction": "Biology is the scientific study of life and living organisms.",
				"sections": []interface{}{
					map[string]interface{}{
						"title":   "Meaning of Biology",
						"content": "Biology is the branch of science concerned with the study of living organisms.",
					},
					map[string]interface{}{
						"title":   "Living Organisms",
						"content": "Living organisms include plants, animals, fungi, bacteria and other forms of life.",
					},
				},
				"key_terms": []interface{}{
					"Biology",
					"organism",
					"life",
				},
				"summary": "Biology is the scientific study of living organisms.",
			}),
		},

		{
			ID:          uuid.MustParse("21f8fbdc-bb30-491c-8f1b-c3d052033bda"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID1,
			TopicID:     uuid.MustParse("825a9320-b960-4bbd-b884-54aef0dd6783"),
			LessonOrder: 1,
			Week:        1,
			Title:       "Branches of Biology",
			Description: "Study of the major branches of Biology and their areas of specialization.",
			Objectives:  "Students should be able to identify major branches of Biology and explain what each branch studies.",
			Activities:  "Teacher introduces major branches. Students match branches with their areas of study.",
			Resources:   "Biology textbook, branch classification chart.",
			Assessment:  "List five branches of Biology and state what each studies.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"introduction": "Biology is divided into several specialized branches.",
				"sections": []interface{}{
					map[string]interface{}{
						"title":   "Zoology",
						"content": "Zoology is the study of animals.",
					},
					map[string]interface{}{
						"title":   "Botany",
						"content": "Botany is the study of plants.",
					},
					map[string]interface{}{
						"title":   "Microbiology",
						"content": "Microbiology is the study of microorganisms.",
					},
					map[string]interface{}{
						"title":   "Ecology",
						"content": "Ecology is the study of relationships between organisms and their environment.",
					},
				},
				"summary": "Different branches of Biology focus on different aspects of life.",
			}),
		},

		{
			ID:          uuid.MustParse("9d496443-0566-4504-b6b6-8634ac538bb4"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID1,
			TopicID:     uuid.MustParse("5fde34d0-1382-41d8-bb6a-d757b33f320c"),
			LessonOrder: 1,
			Week:        2,
			Title:       "Importance of Biology",
			Description: "The importance and applications of Biology in everyday life.",
			Objectives:  "Students should be able to explain the importance of Biology in medicine, agriculture and environmental management.",
			Activities:  "Class discussion on how Biology affects everyday life.",
			Resources:   "Textbook, pictures and charts.",
			Assessment:  "State five ways Biology is useful to society.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"introduction": "Biology has many applications in human life.",
				"sections": []interface{}{
					map[string]interface{}{
						"title":   "Medicine",
						"content": "Biology helps in understanding diseases, developing medicines and improving healthcare.",
					},
					map[string]interface{}{
						"title":   "Agriculture",
						"content": "Biology helps farmers understand crops, pests, diseases and animal production.",
					},
					map[string]interface{}{
						"title":   "Environment",
						"content": "Biology helps us understand and protect ecosystems and natural resources.",
					},
				},
				"summary": "Biology is important in medicine, agriculture, environmental management and many other fields.",
			}),
		},

		// =========================================================
		// MODULE 2: THE CELL
		// =========================================================

		{
			ID:          uuid.MustParse("9d895135-c5f3-478f-8945-7a353f8fe40b"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID2,
			TopicID:     uuid.MustParse("8e3f7f43-5373-4cc4-aaf3-98abed49668b"),
			LessonOrder: 1,
			Week:        3,
			Title:       "Cell Theory",
			Description: "Study of the principles and development of the cell theory.",
			Objectives:  "Students should be able to state the major principles of the cell theory.",
			Activities:  "Teacher explains the development of cell theory. Students discuss the three major principles.",
			Resources:   "Microscope, cell charts, Biology textbook.",
			Assessment:  "State the three main principles of cell theory.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"introduction": "The cell theory explains the importance of cells in living organisms.",
				"sections": []interface{}{
					map[string]interface{}{
						"title":   "First Principle",
						"content": "All living organisms are made up of one or more cells.",
					},
					map[string]interface{}{
						"title":   "Second Principle",
						"content": "The cell is the basic structural and functional unit of life.",
					},
					map[string]interface{}{
						"title":   "Third Principle",
						"content": "All cells arise from pre-existing cells.",
					},
				},
				"summary": "Cells are the basic units of life.",
			}),
		},

		{
			ID:          uuid.MustParse("b27fc0a5-51db-44c7-827c-35fde95d6d09"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID2,
			TopicID:     uuid.MustParse("45a2de63-5b9e-4405-9220-211d48421268"),
			LessonOrder: 1,
			Week:        3,
			Title:       "Cell Structure",
			Description: "Study of the structures and components of plant and animal cells.",
			Objectives:  "Students should be able to identify the major structures found in cells.",
			Activities:  "Students draw and label a typical plant and animal cell.",
			Resources:   "Cell diagrams, microscope, prepared slides.",
			Assessment:  "Draw and label a typical animal cell.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"introduction": "Cells contain different structures that perform specialized functions.",
				"sections": []interface{}{
					map[string]interface{}{
						"title":   "Cell Membrane",
						"content": "The cell membrane controls the movement of substances into and out of the cell.",
					},
					map[string]interface{}{
						"title":   "Cytoplasm",
						"content": "The cytoplasm is the region where many cellular activities occur.",
					},
					map[string]interface{}{
						"title":   "Nucleus",
						"content": "The nucleus contains genetic material and controls many activities of the cell.",
					},
				},
				"summary": "A cell contains specialized structures that work together to keep it alive.",
			}),
		},

		{
			ID:          uuid.MustParse("d86a44d8-eb91-4207-acca-eeab47b2b5b8"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID2,
			TopicID:     uuid.MustParse("d4fdb49d-2c35-45c4-a96a-1976149dc526"),
			LessonOrder: 1,
			Week:        4,
			Title:       "Cell Organelles",
			Description: "Study of major cell organelles and their functions.",
			Objectives:  "Students should be able to identify major cell organelles and state their functions.",
			Activities:  "Students study cell diagrams and match organelles to their functions.",
			Resources:   "Cell charts, textbook, microscope.",
			Assessment:  "State the functions of the nucleus, mitochondria, ribosomes and chloroplasts.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"introduction": "Cell organelles are specialized structures that perform specific functions.",
				"organelles": []interface{}{
					map[string]interface{}{
						"name":     "Nucleus",
						"function": "Controls cell activities and contains genetic material.",
					},
					map[string]interface{}{
						"name":     "Mitochondrion",
						"function": "Site of aerobic respiration and energy release.",
					},
					map[string]interface{}{
						"name":     "Ribosome",
						"function": "Site of protein synthesis.",
					},
					map[string]interface{}{
						"name":     "Chloroplast",
						"function": "Site of photosynthesis in plant cells.",
					},
				},
				"summary": "Different organelles perform different functions within cells.",
			}),
		},

		{
			ID:          uuid.MustParse("31930904-b18b-4622-a268-661d92192604"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID2,
			TopicID:     uuid.MustParse("ce9ac158-3c81-42b2-bee2-dc161b409724"),
			LessonOrder: 1,
			Week:        4,
			Title:       "Plant and Animal Cells",
			Description: "Comparison of the structures and functions of plant and animal cells.",
			Objectives:  "Students should be able to identify similarities and differences between plant and animal cells.",
			Activities:  "Students compare labeled diagrams of plant and animal cells.",
			Resources:   "Plant and animal cell charts.",
			Assessment:  "List four differences between plant and animal cells.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"comparison": []interface{}{
					map[string]interface{}{
						"feature": "Cell wall",
						"plant":   "Present",
						"animal":  "Absent",
					},
					map[string]interface{}{
						"feature": "Chloroplast",
						"plant":   "Present in photosynthetic cells",
						"animal":  "Absent",
					},
					map[string]interface{}{
						"feature": "Large permanent vacuole",
						"plant":   "Usually present",
						"animal":  "Usually absent or small",
					},
				},
				"summary": "Plant and animal cells share many structures but also have important differences.",
			}),
		},

		{
			ID:          uuid.MustParse("9ff98707-f88c-454b-bbbb-ed65200e6549"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID2,
			TopicID:     uuid.MustParse("8ba7b662-e85c-47fb-bfa4-e6a1808d1787"),
			LessonOrder: 1,
			Week:        5,
			Title:       "Cell Division",
			Description: "Introduction to cell division and its importance in living organisms.",
			Objectives:  "Students should be able to explain why cells divide and identify major types of cell division.",
			Activities:  "Teacher explains cell division using diagrams.",
			Resources:   "Cell division charts, textbook.",
			Assessment:  "Explain two reasons why cell division is important.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"introduction": "Cell division is the process by which cells produce new cells.",
				"types": []interface{}{
					"Mitosis",
					"Meiosis",
				},
				"importance": []interface{}{
					"Growth",
					"Repair",
					"Replacement of damaged cells",
					"Reproduction",
				},
				"summary": "Cell division is essential for growth, repair and reproduction.",
			}),
		},

		// =========================================================
		// MODULE 3: ORGANIZATION OF LIFE
		// =========================================================

		{
			ID:          uuid.MustParse("a31598f5-5bf5-4078-a548-8d3be5b2aba7"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID3,
			TopicID:     uuid.MustParse("1f9aa811-91a3-41c5-9dfc-9965e26147c4"),
			LessonOrder: 1,
			Week:        5,
			Title:       "Levels of Organization",
			Description: "Study of the levels of organization from cells to organisms.",
			Objectives:  "Students should be able to describe the major levels of biological organization.",
			Activities:  "Students arrange cards representing cells, tissues, organs and systems in the correct order.",
			Resources:   "Charts, cards and textbook.",
			Assessment:  "List the levels of organization in a multicellular organism.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"levels": []interface{}{
					"Cell",
					"Tissue",
					"Organ",
					"Organ system",
					"Organism",
				},
				"summary": "Cells combine to form tissues, tissues form organs, and organs work together in organ systems.",
			}),
		},

		{
			ID:          uuid.MustParse("5f95d639-d9a1-439e-9cc2-30f511ad44da"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID3,
			TopicID:     uuid.MustParse("73b829cd-02bd-4527-9274-07d1307510ab"),
			LessonOrder: 1,
			Week:        6,
			Title:       "Tissues",
			Description: "Study of tissues in plants and animals and their functions.",
			Objectives:  "Students should be able to define tissue and identify examples of plant and animal tissues.",
			Activities:  "Teacher displays tissue diagrams and explains their functions.",
			Resources:   "Microscope, tissue charts, textbook.",
			Assessment:  "Define tissue and give three examples.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "A tissue is a group of similar cells working together to perform a specific function.",
				"examples": []interface{}{
					"Muscle tissue",
					"Nervous tissue",
					"Epithelial tissue",
					"Xylem",
					"Phloem",
				},
				"summary": "Tissues are groups of specialized cells performing particular functions.",
			}),
		},

		{
			ID:          uuid.MustParse("d89641fd-10cc-44ef-9949-233308fbcc11"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID3,
			TopicID:     uuid.MustParse("d1abc70a-f053-4b51-aa82-ee60d6a787d3"),
			LessonOrder: 1,
			Week:        6,
			Title:       "Organs and Organ Systems",
			Description: "Study of organs and organ systems and how they work together.",
			Objectives:  "Students should be able to explain the relationship between tissues, organs and organ systems.",
			Activities:  "Students identify organs and the systems to which they belong.",
			Resources:   "Human body charts, textbook.",
			Assessment:  "Give three examples of organs and state the systems they belong to.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "An organ is made up of different tissues working together to perform a specific function.",
				"examples": []interface{}{
					map[string]interface{}{
						"organ":  "Heart",
						"system": "Circulatory system",
					},
					map[string]interface{}{
						"organ":  "Lung",
						"system": "Respiratory system",
					},
					map[string]interface{}{
						"organ":  "Stomach",
						"system": "Digestive system",
					},
				},
				"summary": "Organs contain tissues and work together as organ systems.",
			}),
		},

		// =========================================================
		// MODULE 4: NUTRITION
		// =========================================================

		{
			ID:          uuid.MustParse("45bc297b-ce57-4a78-82e7-6ef12a1ed518"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID4,
			TopicID:     uuid.MustParse("d293e26c-2aa3-4bab-8a6e-e7497161d55e"),
			LessonOrder: 1,
			Week:        7,
			Title:       "Classes of Food",
			Description: "Study of the major classes of food and their functions.",
			Objectives:  "Students should be able to identify the major classes of food and state their functions.",
			Activities:  "Students classify common foods according to their major nutrients.",
			Resources:   "Food samples, charts, textbook.",
			Assessment:  "List the major classes of food and give one example of each.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"classes": []interface{}{
					"Carbohydrates",
					"Proteins",
					"Fats and oils",
					"Vitamins",
					"Minerals",
					"Water",
				},
				"summary": "The body requires different classes of food for energy, growth, repair and regulation.",
			}),
		},

		{
			ID:          uuid.MustParse("186baf58-f12c-426c-b8f5-8486707de890"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID4,
			TopicID:     uuid.MustParse("8632ca6f-79f5-456d-8cc1-d73cff707d0c"),
			LessonOrder: 1,
			Week:        7,
			Title:       "Food Tests",
			Description: "Practical tests for identifying major food substances.",
			Objectives:  "Students should be able to describe tests for starch, protein, reducing sugar and fats.",
			Activities:  "Students perform simple food tests under teacher supervision.",
			Resources:   "Iodine solution, Benedict's solution, Biuret reagent, ethanol and food samples.",
			Assessment:  "Describe the test for starch and state the expected result.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"tests": []interface{}{
					map[string]interface{}{
						"food":   "Starch",
						"test":   "Add iodine solution",
						"result": "Blue-black colour indicates starch.",
					},
					map[string]interface{}{
						"food":   "Protein",
						"test":   "Biuret test",
						"result": "Purple colour indicates protein.",
					},
					map[string]interface{}{
						"food":   "Reducing sugar",
						"test":   "Benedict's test",
						"result": "On heating, a coloured precipitate may form.",
					},
				},
				"summary": "Chemical tests can be used to identify nutrients in food.",
			}),
		},

		{
			ID:          uuid.MustParse("f978a62a-f843-41cd-a61c-0e8684511994"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID4,
			TopicID:     uuid.MustParse("9262a8f7-b0a7-409b-9b4f-c5e0a8dcd285"),
			LessonOrder: 1,
			Week:        8,
			Title:       "Balanced Diet",
			Description: "Study of balanced diets and the importance of adequate nutrition.",
			Objectives:  "Students should be able to explain balanced diet and identify factors affecting dietary requirements.",
			Activities:  "Students design a sample balanced meal.",
			Resources:   "Food charts, nutrition guides, textbook.",
			Assessment:  "Define balanced diet and explain why it is important.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "A balanced diet contains all essential nutrients in the correct proportions.",
				"factors": []interface{}{
					"Age",
					"Sex",
					"Activity level",
					"Health condition",
				},
				"summary": "A balanced diet provides the nutrients needed for healthy growth and body function.",
			}),
		},

		{
			ID:          uuid.MustParse("3ac6dd07-612c-4203-861e-93c218509512"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID4,
			TopicID:     uuid.MustParse("caf5e77e-ead6-468a-853f-27b5133e9fbe"),
			LessonOrder: 1,
			Week:        8,
			Title:       "Modes of Nutrition",
			Description: "Study of autotrophic and heterotrophic modes of nutrition.",
			Objectives:  "Students should be able to distinguish between autotrophic and heterotrophic nutrition.",
			Activities:  "Students compare feeding methods in plants and animals.",
			Resources:   "Plant charts, textbook.",
			Assessment:  "Differentiate between autotrophic and heterotrophic nutrition.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"types": []interface{}{
					map[string]interface{}{
						"type":        "Autotrophic nutrition",
						"description": "Organisms manufacture their own food, usually through photosynthesis.",
					},
					map[string]interface{}{
						"type":        "Heterotrophic nutrition",
						"description": "Organisms obtain food from other organisms.",
					},
				},
				"summary": "Organisms obtain nutrients either by producing their own food or by consuming other organisms.",
			}),
		},

		// =========================================================
		// MODULE 5: REPRODUCTION
		// =========================================================

		{
			ID:          uuid.MustParse("8474ec8c-e72c-4f1c-8d86-02fc3a1c8007"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID5,
			TopicID:     uuid.MustParse("daf2c83a-91c4-4b01-88a6-18f03126174b"),
			LessonOrder: 1,
			Week:        9,
			Title:       "Meaning of Reproduction",
			Description: "Introduction to reproduction and its importance for continuity of life.",
			Objectives:  "Students should be able to define reproduction and explain its importance.",
			Activities:  "Class discussion about how organisms produce offspring.",
			Resources:   "Biology textbook, reproductive system charts.",
			Assessment:  "Define reproduction and state two reasons why it is important.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "Reproduction is the biological process by which organisms produce new individuals of the same species.",
				"importance": []interface{}{
					"Continuation of species",
					"Replacement of organisms",
					"Transfer of genetic information",
				},
				"summary": "Reproduction ensures the continuation of living organisms and their species.",
			}),
		},

		{
			ID:          uuid.MustParse("7f14af9f-18c5-40ff-8d61-bd9c1c1a0824"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID5,
			TopicID:     uuid.MustParse("6820228b-1c40-434d-83d2-4e3fb862af48"),
			LessonOrder: 1,
			Week:        9,
			Title:       "Asexual Reproduction",
			Description: "Study of reproduction involving a single parent.",
			Objectives:  "Students should be able to define asexual reproduction and identify its major forms.",
			Activities:  "Students study examples of binary fission, budding and vegetative propagation.",
			Resources:   "Charts, textbook, plant specimens.",
			Assessment:  "List three methods of asexual reproduction.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "Asexual reproduction involves the production of offspring from a single parent without the fusion of gametes.",
				"methods": []interface{}{
					"Binary fission",
					"Budding",
					"Spore formation",
					"Vegetative propagation",
				},
				"summary": "Asexual reproduction generally produces offspring that are genetically similar to the parent.",
			}),
		},

		{
			ID:          uuid.MustParse("1bae2bb0-2739-40a2-8747-910f7226348f"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID5,
			TopicID:     uuid.MustParse("e978a030-0166-4d0a-889c-a604575dd8ed"),
			LessonOrder: 1,
			Week:        10,
			Title:       "Sexual Reproduction",
			Description: "Study of sexual reproduction and the involvement of male and female gametes.",
			Objectives:  "Students should be able to explain sexual reproduction and define fertilization.",
			Activities:  "Teacher explains gametes and fertilization using diagrams.",
			Resources:   "Reproductive system charts, textbook.",
			Assessment:  "Explain the meaning of fertilization.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "Sexual reproduction involves the fusion of male and female gametes.",
				"key_terms": []interface{}{
					"Gamete",
					"Fertilization",
					"Zygote",
				},
				"summary": "Sexual reproduction involves the fusion of gametes to form a zygote.",
			}),
		},

		{
			ID:          uuid.MustParse("78214568-255b-48d8-b65b-053d7c869845"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID5,
			TopicID:     uuid.MustParse("ccc69310-7aad-4303-8490-8527feee0fa7"),
			LessonOrder: 1,
			Week:        10,
			Title:       "Reproduction in Flowering Plants",
			Description: "Study of reproductive structures and processes in flowering plants.",
			Objectives:  "Students should be able to identify the reproductive parts of a flower and explain their functions.",
			Activities:  "Students dissect and label a flower.",
			Resources:   "Fresh flowers, hand lens, diagrams and textbook.",
			Assessment:  "Draw and label the reproductive parts of a flower.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"male_part":   "Stamen",
				"female_part": "Carpel or pistil",
				"processes": []interface{}{
					"Pollination",
					"Fertilization",
					"Seed formation",
					"Fruit formation",
				},
				"summary": "Flowers contain reproductive structures that enable sexual reproduction in flowering plants.",
			}),
		},

		// =========================================================
		// MODULE 6: ECOLOGY
		// =========================================================

		{
			ID:          uuid.MustParse("6965d83a-e130-4632-a1eb-bfd8ce153665"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID6,
			TopicID:     uuid.MustParse("53f9dccf-54d8-4174-b05e-3393e7cfd076"),
			LessonOrder: 1,
			Week:        11,
			Title:       "Meaning of Ecology",
			Description: "Introduction to ecology and relationships between organisms and their environment.",
			Objectives:  "Students should be able to define ecology and explain its importance.",
			Activities:  "Students observe organisms and environmental factors around the school.",
			Resources:   "School environment, textbook, ecology charts.",
			Assessment:  "Define ecology.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "Ecology is the study of relationships between organisms and their environment.",
				"key_terms": []interface{}{
					"Organism",
					"Environment",
					"Habitat",
					"Population",
					"Community",
				},
				"summary": "Ecology examines how organisms interact with one another and with their environment.",
			}),
		},

		{
			ID:          uuid.MustParse("2637a9a4-c923-475f-85eb-59db84a76f79"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID6,
			TopicID:     uuid.MustParse("ffae1027-7b9c-4778-b336-4b3cbbe7064e"),
			LessonOrder: 1,
			Week:        11,
			Title:       "Ecosystem",
			Description: "Study of ecosystems and their biotic and abiotic components.",
			Objectives:  "Students should be able to define an ecosystem and identify its components.",
			Activities:  "Students identify biotic and abiotic components of the school ecosystem.",
			Resources:   "School environment, charts, textbook.",
			Assessment:  "List five biotic and five abiotic components of an ecosystem.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "An ecosystem consists of living organisms interacting with one another and with non-living components of their environment.",
				"components": map[string]interface{}{
					"biotic": []interface{}{
						"Plants",
						"Animals",
						"Microorganisms",
					},
					"abiotic": []interface{}{
						"Light",
						"Temperature",
						"Water",
						"Soil",
						"Air",
					},
				},
				"summary": "An ecosystem contains both living and non-living components.",
			}),
		},

		{
			ID:          uuid.MustParse("f1f4c04f-f614-4e26-b5ab-ef20ea05a3b8"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID6,
			TopicID:     uuid.MustParse("620a5bae-3087-4082-9bb5-a7334ebedf9d"),
			LessonOrder: 1,
			Week:        12,
			Title:       "Food Chains and Food Webs",
			Description: "Study of feeding relationships among organisms in an ecosystem.",
			Objectives:  "Students should be able to construct simple food chains and explain food webs.",
			Activities:  "Students construct food chains using local organisms.",
			Resources:   "Food chain charts, cards and textbook.",
			Assessment:  "Construct a food chain containing a producer and three consumers.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"food_chain_example": []interface{}{
					"Grass",
					"Grasshopper",
					"Frog",
					"Snake",
					"Hawk",
				},
				"key_terms": []interface{}{
					"Producer",
					"Consumer",
					"Predator",
					"Prey",
				},
				"summary": "Food chains show the transfer of energy from one organism to another.",
			}),
		},

		{
			ID:          uuid.MustParse("f302ba93-c00b-4f4e-a777-aae618716bcf"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID6,
			TopicID:     uuid.MustParse("0e316b17-89a9-4929-8b4a-7017630c0ef7"),
			LessonOrder: 1,
			Week:        12,
			Title:       "Environmental Factors",
			Description: "Study of biotic and abiotic factors affecting living organisms.",
			Objectives:  "Students should be able to distinguish between biotic and abiotic environmental factors.",
			Activities:  "Students identify environmental factors around the school.",
			Resources:   "School environment, thermometer, charts and textbook.",
			Assessment:  "List five abiotic factors that affect organisms.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"biotic_factors": []interface{}{
					"Competition",
					"Predation",
					"Disease",
					"Availability of food",
				},
				"abiotic_factors": []interface{}{
					"Temperature",
					"Light",
					"Water",
					"Soil",
					"Humidity",
				},
				"summary": "Environmental factors influence the survival, growth and distribution of organisms.",
			}),
		},

		// =========================================================
		// MODULE 7: CLASSIFICATION OF LIVING ORGANISMS
		// =========================================================

		{
			ID:          uuid.MustParse("ac3d0027-5499-493e-ad10-f7428b22c2f8"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID7,
			TopicID:     uuid.MustParse("17a486e0-91a6-404a-9705-b595a3cb6fc5"),
			LessonOrder: 1,
			Week:        13,
			Title:       "Meaning of Classification",
			Description: "Introduction to the classification of living organisms.",
			Objectives:  "Students should be able to define classification and explain why organisms are classified.",
			Activities:  "Students group objects according to common characteristics.",
			Resources:   "Classification charts, textbook.",
			Assessment:  "Define biological classification.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "Classification is the systematic arrangement of organisms into groups based on shared characteristics.",
				"importance": []interface{}{
					"Easy identification",
					"Study of relationships",
					"Organization of biological information",
				},
				"summary": "Classification organizes organisms into groups based on their similarities and differences.",
			}),
		},

		{
			ID:          uuid.MustParse("30404e8a-4c5f-47b1-9e9d-10e2a4a7ee5c"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID7,
			TopicID:     uuid.MustParse("af3a3497-a4b1-4db0-a918-5e31d231270a"),
			LessonOrder: 1,
			Week:        13,
			Title:       "Importance of Classification",
			Description: "Study of the importance of classifying living organisms.",
			Objectives:  "Students should be able to explain the importance of classification in Biology.",
			Activities:  "Students discuss problems that would arise without biological classification.",
			Resources:   "Biology textbook, classification charts.",
			Assessment:  "State four importance of biological classification.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"importance": []interface{}{
					"Helps identify organisms",
					"Provides a systematic way to study organisms",
					"Shows relationships among organisms",
					"Reduces confusion caused by local names",
				},
				"summary": "Classification makes the study and identification of organisms easier.",
			}),
		},

		{
			ID:          uuid.MustParse("195a2757-1dcf-4665-8128-9d847d07d65d"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID7,
			TopicID:     uuid.MustParse("053e2dea-c4a7-402c-beba-69d9ed1a8e43"),
			LessonOrder: 1,
			Week:        14,
			Title:       "Kingdoms of Living Organisms",
			Description: "Introduction to the major kingdoms used in biological classification.",
			Objectives:  "Students should be able to identify the major kingdoms and give examples of organisms in each.",
			Activities:  "Students classify familiar organisms into their appropriate kingdoms.",
			Resources:   "Classification chart, Biology textbook.",
			Assessment:  "List the major kingdoms and give one example of each.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"kingdoms": []interface{}{
					map[string]interface{}{
						"name":    "Monera",
						"example": "Bacteria",
					},
					map[string]interface{}{
						"name":    "Protista",
						"example": "Amoeba",
					},
					map[string]interface{}{
						"name":    "Fungi",
						"example": "Mushroom",
					},
					map[string]interface{}{
						"name":    "Plantae",
						"example": "Maize",
					},
					map[string]interface{}{
						"name":    "Animalia",
						"example": "Human",
					},
				},
				"summary": "Living organisms can be grouped into kingdoms based on their characteristics.",
			}),
		},

		{
			ID:          uuid.MustParse("62bcb354-f483-47f5-a269-fc218a6712bc"),
			SchemeOfWorkID: schemeOfWorkID,
			ModuleID:    moduleID7,
			TopicID:     uuid.MustParse("8e58b2e2-bd1c-402f-8504-f81b8575b156"),
			LessonOrder: 1,
			Week:        14,
			Title:       "Binomial Nomenclature",
			Description: "Study of the system used to give organisms scientific names.",
			Objectives:  "Students should be able to explain binomial nomenclature and correctly write scientific names.",
			Activities:  "Students practice writing scientific names using examples.",
			Resources:   "Classification textbook and charts.",
			Assessment:  "Explain the two parts of a scientific name.",
			CreatedBy:   adminID,
			Content: toJSON(map[string]interface{}{
				"definition": "Binomial nomenclature is a system of naming organisms using two names.",
				"parts": []interface{}{
					"Genus",
					"Species",
				},
				"examples": []interface{}{
					"Homō sapiens",
					"Zea mays",
				},
				"rules": []interface{}{
					"Genus begins with a capital letter.",
					"Species begins with a lowercase letter.",
					"The scientific name is usually italicized when typed.",
				},
				"summary": "Binomial nomenclature gives each organism a standardized two-part scientific name.",
			}),
		},
	}

	// ---------------------------------------------------------
	// Insert lessons
	// ---------------------------------------------------------

	for _, lesson := range lessons {
		if err := db.Create(&lesson).Error; err != nil {
			log.Printf("❌ Failed to seed lesson %s: %v", lesson.Title, err)
		} else {
			log.Printf("✅ Seeded Lesson: %s", lesson.Title)
		}
	}
}
package seeds

import (
	"log"
	"crm-go/config"
	"crm-go/models"
	"gorm.io/gorm/clause"
	"github.com/google/uuid"
)

func SeedLessons() error{
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
			Content: "<div>Biology is the scientific study of life and living organisms, including their structure, function, growth, evolution, distribution, and taxonomy. It encompasses various fields such as botany, zoology, microbiology, and ecology.</div>",
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
			Content: "<div>The branches of Biology include botany (study of plants), zoology (study of animals), microbiology (study of microorganisms), ecology (study of interactions between organisms and their environment), and genetics (study of heredity and variation).</div>",
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
			Content: "<div> The importance of Biology in everyday life cannot be overstated. It helps us understand the world around us and provides solutions to many of the challenges we face. </div>",
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
			Content: "<p>The cell theory is a fundamental concept in Biology that describes the properties of cells. It states that all living organisms are composed of cells, that the cell is the basic unit of life, and that all cells arise from pre-existing cells.</p>",
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
			Content: "<div>Cells are the basic structural and functional units of life. They contain various organelles that perform specific functions necessary for the cell's survival and operation.</div>",
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
			Content: "<div>Cell organelles are specialized structures within a cell that perform distinct processes. For example, the nucleus controls cell activities, mitochondria produce energy, ribosomes synthesize proteins, and chloroplasts conduct photosynthesis in plant cells.</div>",
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
			Content: "<div>Plant and animal cells share many common structures, such as the nucleus, cytoplasm, and cell membrane. However, plant cells have a cell wall, chloroplasts, and large central vacuoles, which are not found in animal cells.</div>",
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
			Content: "<div>Cell division is a fundamental process by which a parent cell divides into two or more daughter cells. It is essential for growth, repair, and reproduction in living organisms. The two main types of cell division are mitosis and meiosis.</div>",
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
			Content: "<div>The levels of organization in living organisms include cells, tissues, organs, organ systems, and the organism as a whole. Each level represents a higher degree of complexity and specialization.</div>",
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
			Content: "<div>Tissues are groups of similar cells that work together to perform a specific function. In animals, examples include epithelial tissue, muscle tissue, and nervous tissue. In plants, examples include xylem and phloem.</div>",
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
			Content: "<div>Organs are structures composed of different tissues that work together to perform specific functions. Organ systems are groups of organs that collaborate to carry out complex bodily functions. For example, the heart and blood vessels form the circulatory system.</div>",
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
			Content: "<div>Food can be classified into major classes based on the nutrients they provide. The main classes include carbohydrates, proteins, fats, vitamins, minerals, and water. Each class has specific functions in the body, such as providing energy, building tissues, and regulating body processes.</div>",
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
			Content: "<div>Food tests are simple chemical tests used to identify the presence of specific nutrients in food samples. For example, the iodine test is used for starch, Benedict's test for reducing sugars, Biuret test for proteins, and the ethanol emulsion test for fats.</div>",
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
			Content: "<div>A balanced diet is one that provides all the essential nutrients in the right proportions to maintain health and well-being. It includes a variety of foods from different food groups, ensuring that the body receives adequate energy, protein, vitamins, and minerals.</div>",
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
			Content: "<div>Autotrophic nutrition is the process by which organisms produce their own food, typically through photosynthesis, as seen in plants. Heterotrophic nutrition involves obtaining food from other organisms, as seen in animals and humans.</div>",
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
			Content: "<div>Reproduction is the biological process by which new individual organisms are produced. It is essential for the continuation of species and the transfer of genetic information from one generation to the next.</div>",
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
			Content: "<div>Asexual reproduction is a mode of reproduction that involves a single parent and results in offspring that are genetically identical to the parent. Major forms include binary fission, budding, and vegetative propagation.</div>",
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
			Content: "<div>Sexual reproduction involves the fusion of male and female gametes to produce offspring that are genetically different from their parents. Fertilization is the process by which the male gamete (sperm) fuses with the female gamete (egg) to form a zygote.</div>",
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
			Content: "<div>Flowering plants reproduce sexually through the production of flowers, which contain male and female reproductive organs. The male part is the stamen, which produces pollen, while the female part is the carpel, which contains the ovary, style, and stigma. Pollination and fertilization lead to the formation of seeds and fruits.</div>",
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
			Content: "<div>Ecology is the branch of Biology that studies the interactions between organisms and their environment. It examines how living organisms adapt to their surroundings, how they interact with each other, and how they are affected by environmental factors.</div>",
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
			Content: "<div>An ecosystem is a community of living organisms interacting with their physical environment. It includes both biotic components (living things) and abiotic components (non-living things) that influence the survival and growth of organisms.</div>",
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
			Content: "<div>A food chain is a linear sequence of organisms through which nutrients and energy pass as one organism eats another. A food web is a complex network of interconnected food chains that shows the feeding relationships in an ecosystem.</div>",
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
			Content: "<div>Environmental factors can be classified as biotic (living) or abiotic (non-living). Biotic factors include other organisms, while abiotic factors include temperature, light, water, and soil. Both types of factors influence the survival and distribution of organisms.</div>",
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
			Content: "<div>Classification is the process of grouping organisms based on shared characteristics. It helps scientists organize and understand the diversity of life, making it easier to study and communicate about different species.</div>",
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
			Content: "<div>Classification is important because it helps scientists organize and understand the diversity of life, facilitates communication about organisms, aids in identifying and naming species, and provides insights into evolutionary relationships.</div>",
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
			Content: "<div>The major kingdoms of living organisms include Animalia (animals), Plantae (plants), Fungi (fungi), Protista (protists), and Monera (bacteria). Each kingdom groups organisms based on shared characteristics and evolutionary relationships.</div>",
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
			Content: "<div>Binomial nomenclature is a system of naming organisms using two names: the genus name (capitalized) and the species name (lowercase). For example, the scientific name for humans is Homo sapi",
		},
	}


for _, lesson := range lessons {
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "id"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"scheme_of_work_id",
			"module_id",
			"topic_id",
			"lesson_order",
			"week",
			"title",
			"description",
			"objectives",
			"activities",
			"resources",
			"assessment",
			"created_by",
			"content",
		}),
	}).Create(&lesson)

	if result.Error != nil {
		log.Printf(
			"❌ Failed to seed lesson %s: %v",
			lesson.Title,
			result.Error,
		)
	} else {
		log.Printf(
			"✅ Seeded/updated Lesson: %s",
			lesson.Title,
		)
	}
}
return nil
}
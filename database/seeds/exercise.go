package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func SeedExercises() error {
	db := config.GetDB()
	adminID := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	exercise1 := uuid.MustParse("c23e2651-5438-4418-af00-c045965e89f4")
	exercise2 := uuid.MustParse("65e95115-54bf-4f5d-b873-9e2ce26e130d")
	exercise3 := uuid.MustParse("c97ef6d8-12e1-4995-b758-b12d1f735de9")
	exercise4 := uuid.MustParse("10866863-448e-4cc8-814c-fcfb4461d654")
	exercise5 := uuid.MustParse("44cc0821-0093-43fa-9d34-da65663a746d")
	exercise6 := uuid.MustParse("a05deb1e-dd4e-464a-90cc-048f4ebb87fa")
	exercise7 := uuid.MustParse("21be7c52-f017-4669-931e-2161ab69ebfc")
	exercise8 := uuid.MustParse("734f3c5e-5aa4-4b7d-a313-8f3bfebd6866")
	exercise9 := uuid.MustParse("e743abcc-2121-48f7-986f-16457b7f6fa2")
	exercise10 := uuid.MustParse("f7518f05-2e8d-452d-947b-374dd2c4e92e")
	exercise11 := uuid.MustParse("9f4a8cfb-c1c1-4e01-9996-c827ee1f10e4")
	exercise12 := uuid.MustParse("6125b941-fd19-46a6-be8b-774b5ae70a85")
	exercise13 := uuid.MustParse("a3637ec8-b2a5-4d38-8e2f-99a0f61272b3")

	exercises := []models.Exercise{
		{
			ID:       exercise1,
			LessonID: uuid.MustParse("e2a83449-fd98-4b24-8fcf-e51b4e17456b"),
			Title:    "Exercise on Introduction to Micro-organisms",
			Content: `
				1. Define micro-organisms. (2 marks)

				2. List five groups of micro-organisms. (5 marks)

				3. State two characteristics of bacteria. (2 marks)

				4. Give two examples of useful micro-organisms. (2 marks)

				5. Mention two harmful effects of micro-organisms. (2 marks)

				6. Explain why micro-organisms cannot usually be seen with the naked eye. (2 marks)`,
			Instructions: "Answer all questions. Write your answers clearly and give examples where required.",
			TotalMarks:   15,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise2,
			LessonID: uuid.MustParse("9e12aaf1-fdb9-430e-8f5e-e98c36ec3bf9"),
			Title:    "Exercise on Culturing and Identification of Micro-organisms",
			Content: `
				1. What is meant by culturing micro-organisms? (2 marks)

				2. State four conditions necessary for the growth of micro-organisms. (4 marks)

				3. What is a culture medium? (2 marks)

				4. Mention three methods that can be used to identify micro-organisms. (3 marks)

				5. State four safety precautions that should be observed when working with micro-organisms. (4 marks)`,
			Instructions: "Answer all questions. Use appropriate biological terms in your answers.",
			TotalMarks:   15,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise3,
			LessonID: uuid.MustParse("d837b61a-d747-4cb9-bca8-c165bc81e6fe"),
			Title:    "Exercise on Sexually Transmitted Infections",
			Content: `
				1. Define sexually transmitted infections. (2 marks)

				2. List five examples of sexually transmitted infections. (5 marks)

				3. State four ways through which STIs can be transmitted. (4 marks)

				4. Mention three signs or symptoms that may occur with STIs. (3 marks)

				5. State three ways of preventing STIs. (3 marks)

				6. Explain why early diagnosis and treatment of STIs are important. (3 marks)`,
			Instructions: "Answer all questions. Give clear and scientifically accurate answers.",
			TotalMarks:   20,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise4,
			LessonID: uuid.MustParse("caafaae2-7efa-4cd8-89fb-a893c35a88e2"),
			Title:    "Exercise on Micro-organisms in Action",
			Content: `
				1. State four useful activities of micro-organisms. (4 marks)

				2. Explain the role of micro-organisms in decomposition. (3 marks)

				3. What is fermentation? (2 marks)

				4. Mention three products that can be produced through microbial fermentation. (3 marks)

				5. State three harmful effects of micro-organisms. (3 marks)

				6. Explain one role of micro-organisms in agriculture. (2 marks)`,
			Instructions: "Answer all questions. Where applicable, support your answers with examples.",
			TotalMarks:   17,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise5,
			LessonID: uuid.MustParse("3f93af4b-9667-4772-8ada-b53875136304"),
			Title:    "Exercise on Causes, Effects and Prevention of STIs",
			Content: `
				1. State four risk factors associated with STIs. (4 marks)

				2. Mention four effects of untreated STIs. (4 marks)

				3. Explain three ways of preventing STIs. (6 marks)

				4. Why is testing important in the control of STIs? (2 marks)

				5. State two possible social effects of untreated STIs. (2 marks)`,
			Instructions: "Answer all questions. Explain your answers where necessary.",
			TotalMarks:   18,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise6,
			LessonID: uuid.MustParse("710f0501-9c6b-49f7-b8b9-5b1ef22f519f"),
			Title:    "Exercise on Control of Harmful Micro-organisms",
			Content: `
				1. What is sterilization? (2 marks)

				2. Differentiate between sterilization and disinfection. (4 marks)

				3. State five methods of controlling harmful micro-organisms. (5 marks)

				4. Explain how refrigeration helps to control micro-organisms. (2 marks)

				5. Mention three chemicals commonly used to control micro-organisms. (3 marks)

				6. State two importance of proper hygiene and sanitation. (2 marks)`,
			Instructions: "Answer all questions. Use examples from everyday life where appropriate.",
			TotalMarks:   18,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise7,
			LessonID: uuid.MustParse("90169597-ef87-4451-8d79-e6f3451b344d"),
			Title:    "Exercise on Classification of Plants",
			Content: `
				1. What is plant classification? (2 marks)

				2. State five major groups of plants. (5 marks)

				3. Give two characteristics of bryophytes. (2 marks)

				4. State three characteristics of flowering plants. (3 marks)

				5. Differentiate between monocotyledons and dicotyledons. (4 marks)

				6. Give two examples each of monocotyledonous and dicotyledonous plants. (4 marks)`,
			Instructions: "Answer all questions. Use examples to support your answers where necessary.",
			TotalMarks:   20,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise8,
			LessonID: uuid.MustParse("e6492333-05af-4d29-8e38-8f83d050304d"),
			Title:    "Exercise on Plant Pests and Diseases",
			Content: `1. Define a plant pest. (2 marks)

2. List five common plant pests. (5 marks)

3. State four signs of plant disease. (4 marks)

4. Mention three effects of plant diseases on crop production. (3 marks)

5. Differentiate between a plant pest and a plant disease. (3 marks)

6. Give two examples of crops and diseases that commonly affect them. (4 marks)`,
			Instructions: "Answer all questions. Clearly distinguish between pests and diseases.",
			TotalMarks:   21,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise9,
			LessonID: uuid.MustParse("ae02e07e-e0f9-4d58-aa29-1f7f1d35e382"),
			Title:    "Exercise on Animal Pests and Diseases",
			Content: `1. What is an animal pest? (2 marks)

2. List four common animal pests. (4 marks)

3. Name five diseases that affect farm animals. (5 marks)

4. State four signs that may indicate disease in livestock. (4 marks)

5. Differentiate between internal and external parasites. (4 marks)

6. State two effects of animal diseases on livestock production. (2 marks)`,
			Instructions: "Answer all questions. Give relevant examples from livestock production.",
			TotalMarks:   21,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise10,
			LessonID: uuid.MustParse("ebe29753-3c96-4a0b-a7c0-5539d4d63a0b"),
			Title:    "Exercise on Control of Animal Diseases",
			Content: `1. State five methods of controlling animal diseases. (5 marks)

2. What is vaccination? (2 marks)

3. Explain the importance of quarantine in livestock management. (3 marks)

4. State four good livestock management practices that help prevent diseases. (4 marks)

5. Why should sick animals be isolated from healthy animals? (2 marks)

6. State two roles of veterinary professionals in controlling animal diseases. (2 marks)`,
			Instructions: "Answer all questions. Explain your answers where required.",
			TotalMarks:   18,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise11,
			LessonID: uuid.MustParse("e4698004-cad7-4c38-979c-4d0e78101bc6"),
			Title:    "Exercise on Food Production",
			Content: `1. Define food production. (2 marks)

2. Identify five sources of food. (5 marks)

3. State four methods of crop production. (4 marks)

4. Mention four factors that affect food production. (4 marks)

5. Explain the importance of livestock production in food supply. (3 marks)

6. State two ways in which fisheries contribute to food production. (2 marks)`,
			Instructions: "Answer all questions. Use examples from your local environment where possible.",
			TotalMarks:   20,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise12,
			LessonID: uuid.MustParse("3c8f106c-0203-4c69-bfb2-f507b1fc5819"),
			Title:    "Exercise on Food Preservation and Storage",
			Content: `1. Define food preservation. (2 marks)

2. State six methods of preserving food. (6 marks)

3. Explain how drying preserves food. (3 marks)

4. Explain how refrigeration helps to preserve food. (3 marks)

5. Differentiate between food preservation and food storage. (4 marks)

6. State four good food storage practices. (4 marks)`,
			Instructions: "Answer all questions. Give practical examples where applicable.",
			TotalMarks:   22,
			CreatedBy:    adminID,
		},

		{
			ID:       exercise13,
			LessonID: uuid.MustParse("9983316d-d951-4ba5-8959-90faa942dd49"),
			Title:    "Revision Exercise on Term's Work",
			Content: `1. Define micro-organisms and list five groups of micro-organisms. (7 marks)

2. State four useful activities of micro-organisms. (4 marks)

3. Mention four sexually transmitted infections and state three methods of prevention. (7 marks)

4. State five methods of controlling harmful micro-organisms. (5 marks)

5. List five major groups of plants. (5 marks)

6. State four common plant pests and three effects of plant diseases. (7 marks)

7. Mention four animal diseases and four methods of controlling animal diseases. (8 marks)

8. State five methods of food preservation. (5 marks)

9. Explain the importance of proper food storage. (2 marks)`,
			Instructions: "Answer all questions. This is a comprehensive revision exercise covering the major topics studied during the term.",
			TotalMarks:   50,
			CreatedBy:    adminID,
		},
	}

	for _, exercise := range exercises {
		result := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"lesson_id",
				"title",
				"content",
				"instructions",
				"total_marks",
				"created_by",
			}),
		}).Create(&exercise)

		if result.Error != nil {
			log.Printf(
				"❌ Failed to seed exercise %s: %v",
				exercise.Title,
				result.Error,
			)
			continue
		}

		log.Printf(
			"✅ Seeded/updated exercise: %s",
			exercise.Title,
		)
	}
	return nil
}

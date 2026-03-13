package seeds

import (
	"log"
	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
)

func SeedObjectiveQuestions() {
	db := config.GetDB()


	// ✅ Validate UUIDs for GoLang Lessons
	questionID1, err01 := uuid.Parse("18547b5a-a3b3-473b-a02e-53c259cb60fb")

	if err01 != nil {
		log.Fatalf("❌ Invalid question 1 UUID: %v", err01)
	}
	questionOptionID1, err1 := uuid.Parse("8c9065ad-c6da-4b56-80c3-1ef2409b3e23")
	questionOptionID2, err2 := uuid.Parse("9861e609-7402-4852-9847-4b7f8f35bcf5")
	questionOptionID3, err3 := uuid.Parse("fd06b822-a339-4fe3-83fb-811d4f3fef62")
	questionOptionID4, err4 := uuid.Parse("309ecc66-4784-42ef-9d46-d7a3080dee81")


	if err1 != nil {
		log.Fatalf("❌ Invalid question option 1 UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid question option 2 UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid question option 3 UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid question option 4 UUID: %v", err4)
	}


	// ✅ Validate UUIDs for GoLang Lessons
	lessonID1, err1 := uuid.Parse("05d294c9-80f1-4203-8e78-d5f2c6d74c2c")



	if err1 != nil {
		log.Fatalf("❌ Invalid lesson 1 UUID: %v", err1)
	}


		// ✅ Validate UUIDs  Courses
	courseId1, err26 := uuid.Parse("27d8ae14-4311-4380-8397-057ad5043fd6") // Go Lang Course

	if err26 != nil {
		log.Fatalf("❌ Invalid Course 26 UUID: %v", err26)
	}


	// ✅ Validate UUIDs for Go Lang Course
	moduleID31, err31 := uuid.Parse("1fb7c5b1-a438-4a19-9eb7-19a9046d0124")

	
	if err31 != nil {
		log.Fatalf("❌ Invalid module 31 UUID: %v", err31)
	}




	createdByID, err2 := uuid.Parse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	if err2 != nil {
		log.Fatalf("❌ Invalid CreatedBy UUID: %v", err2)
	}



	questions := []models.ObjectiveQuestion{
		// Go Lang Question  Lesson: What is Go?
	{
		ID:			  	questionID1,
		QuestionText:    "What is Go?",
		QuestionType:    "multiple_choice",
		DifficultyLevel: "easy",
		Points:          1,
		ImageURL:        "https://wwww.imagequestion.com",
		VideoURL:        "https://wwww.videoquestion.com",
		CourseID:        courseId1,
		ModuleID:        moduleID31,
		LessonID:        lessonID1,
		TutorID:         createdByID,
		AnswerExplanation: "Go is an open-source programming language designed for simplicity, efficiency, and strong concurrency support.",
		SolutionSteps:     "1. Go was created at Google in 2007 by Robert Griesemer, Rob Pike, and Ken Thompson.\n2. It was designed to address issues of scalability and maintainability in large software systems.\n3. Go features a simple syntax, garbage collection, and built-in support for concurrent programming through goroutines and channels.",
		Hint:              "Think about a programming language developed by Google that emphasizes simplicity and performance.",
		IsApproved:        true,

		Options: []models.QuestionOption{
			{
				ID: questionOptionID1,
				QuestionID: questionID1,
				Explanation: "Go is a statically typed, compiled programming language designed for simplicity and efficiency.",
				IsCorrect:   true,
				OptionText:  "Go is a statically typed, compiled programming language designed for simplicity and efficiency.",
				SortOrder:   1,

			},
			{
				ID: questionOptionID2,
				QuestionID: questionID1,
				Explanation: "Go is a dynamically typed scripting language.",
				IsCorrect:   false,
				OptionText:  "Go is a dynamically typed scripting language.",
				SortOrder:   2,
			},
			{
				ID: questionOptionID3,
				QuestionID: questionID1,
				Explanation: "Go is a functional programming language.",
				IsCorrect:   false,
				OptionText:  "Go is a functional programming language.",
				SortOrder:   3,
			},
			{
				ID: questionOptionID4,
				QuestionID: questionID1,
				Explanation: "Go is a systems programming language.",
				IsCorrect:   false,
				OptionText:  "Go is a systems programming language.",
				SortOrder:   4,
			},
		},
	},
	{
		ID:			  	uuid.New(),
		QuestionText:    "Which of the following is a feature of Go?",
		QuestionType:    "multiple_choice",
		DifficultyLevel: "medium",
		Points:          2,
		ImageURL:        "https://wwww.imagequestion2.com",
		VideoURL:        "https://wwww.videoquestion2.com",
		CourseID:        courseId1,
		ModuleID:        moduleID31,
		LessonID:        lessonID1,
		TutorID:         createdByID,
		AnswerExplanation: "Go has built-in support for concurrent programming through goroutines and channels.",
		SolutionSteps:     "1. Go's concurrency model is based on goroutines, which are lightweight threads managed by the Go runtime.\n2. Channels are used to communicate between goroutines, allowing for safe and efficient concurrent programming.",
		Hint:              "Think about how Go handles concurrency and communication between threads.",
		IsApproved:        true,

		Options: []models.QuestionOption{
			{
				ID: uuid.New(),
				QuestionID: questionID1,
				Explanation: "Go has built-in support for concurrent programming through goroutines and channels.",
				IsCorrect:   true,
				OptionText:  "Go has built-in support for concurrent programming through goroutines and channels.",
				SortOrder:   1,
			},
			{
				ID: uuid.New(),
				QuestionID: questionID1,
				Explanation: "Go is a statically typed, compiled programming language.",
				IsCorrect:   false,
				OptionText:  "Go is a statically typed, compiled programming language.",
				SortOrder:   2,
			},
			{
				ID: uuid.New(),
				QuestionID: questionID1,
				Explanation: "Go is a dynamically typed scripting language.",
				IsCorrect:   false,
				OptionText:  "Go is a dynamically typed scripting language.",
				SortOrder:   3,
			},
			{
				ID: uuid.New(),
				QuestionID: questionID1,
				Explanation: "Go is a systems programming language.",
				IsCorrect:   false,
				OptionText:  "Go is a systems programming language.",
				SortOrder:   4,
			},
		},
	},
}


	for _, question := range questions {
		if err := db.Create(&question).Error; err != nil {
			log.Printf("❌ Failed to seed Question: %v", err)
		} else {
			log.Printf("✅ Seeded Question: %v", &question.QuestionText)
		}
	}
}

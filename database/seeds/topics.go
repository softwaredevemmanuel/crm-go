package seeds

import (
	"log"
	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedTopics() {
	db := config.GetDB()

	// ✅ Validate UUIDs for GoLang Lessons
	lessonID1, err1 := uuid.Parse("05d294c9-80f1-4203-8e78-d5f2c6d74c2c")
	lessonID2, err2 := uuid.Parse("00afe9fe-f686-4c04-8589-07e76beaebe4")
	lessonID3, err3 := uuid.Parse("5be66915-f2eb-49c2-b07a-41f467b1817a")
	lessonID4, err4 := uuid.Parse("813e6ca8-6513-4995-965f-3bb846f0b941")
	lessonID5, err5 := uuid.Parse("c6238448-c520-40ed-8a9b-54c1a5cd5c62")

	// ✅ Validate UUIDs for React Lessons
	lessonID6, err6 := uuid.Parse("5fbf5947-1904-40ef-b1c9-c1630b40e6b6")
	lessonID7, err7 := uuid.Parse("72223633-21be-4df6-9026-e4d7e09ba183")
	lessonID8, err8 := uuid.Parse("a778551d-fed1-47f1-a0b1-e4becd8cdc08")
	lessonID9, err9 := uuid.Parse("e9a14e98-761d-4d18-bd6a-db4e742413be")
	lessonID10, err10 := uuid.Parse("ff7edbb8-4515-49ba-b8b6-bac165bb61da")

	// ✅ Validate UUIDs for Python Lessons
	lessonID11, err11 := uuid.Parse("02770b59-bea7-4010-905f-02599cf25f33")
	lessonID12, err12 := uuid.Parse("47efa3d4-0b73-4c0a-8794-4542ba1847b6")
	lessonID13, err13 := uuid.Parse("648fcd8a-bc3a-4691-bd69-d365a73a463d")
	lessonID14, err14 := uuid.Parse("876e9404-448c-4b13-a3bf-9068ed888bcb")
	lessonID15, err15 := uuid.Parse("f3a4f79c-2a33-400a-b7f6-fea6f161d7f5")

	// ✅ Validate UUIDs for FastAPI Lessons
	lessonID16, err16 := uuid.Parse("9dbf26e9-0ae8-4a9e-90cc-178b25208f9d")
	lessonID17, err17 := uuid.Parse("d1c389a6-7396-45a4-a311-807e72b20d5c")
	lessonID18, err18 := uuid.Parse("da71896a-caf4-429f-b735-e7a847bb7069")
	lessonID19, err19 := uuid.Parse("deeefc60-f790-442d-afe8-51223f642411")
	lessonID20, err20 := uuid.Parse("fd84f5eb-c1cf-4d46-9b4b-118d43e294f8")

	// ✅ Validate UUIDs for Machine Learning Lessons
	lessonID21, err21 := uuid.Parse("49aff6ea-1ba0-410e-9563-a4053ee859db")
	lessonID22, err22 := uuid.Parse("a123f9a4-80ad-4336-ba52-beb76417a9c2")
	lessonID23, err23 := uuid.Parse("cd562cc9-37e1-4bc7-b7a5-8b711e411721")
	lessonID24, err24 := uuid.Parse("d9a7750e-42f4-4afc-8fb2-ded89ce0e2c8")
	lessonID25, err25 := uuid.Parse("ed8f20a1-0267-49d7-a8d4-6e232836423b")

	if err1 != nil {
		log.Fatalf("❌ Invalid lesson 1 UUID: %v", err1)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid lesson 3 UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid lesson 4 UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid lesson 5 UUID: %v", err5)
	}
	if err6 != nil {
		log.Fatalf("❌ Invalid lesson 6 UUID: %v", err6)
	}
	if err7 != nil {
		log.Fatalf("❌ Invalid lesson 7 UUID: %v", err7)
	}
	if err8 != nil {
		log.Fatalf("❌ Invalid lesson 8 UUID: %v", err8)
	}
	if err9 != nil {
		log.Fatalf("❌ Invalid lesson 9 UUID: %v", err9)
	}
	if err10 != nil {
		log.Fatalf("❌ Invalid lesson 10 UUID: %v", err10)
	}
	if err11 != nil {
		log.Fatalf("❌ Invalid lesson 11 UUID: %v", err11)
	}
	if err12 != nil {
		log.Fatalf("❌ Invalid lesson 12 UUID: %v", err12)
	}
	if err13 != nil {
		log.Fatalf("❌ Invalid lesson 13 UUID: %v", err13)
	}
	if err14 != nil {
		log.Fatalf("❌ Invalid lesson 14 UUID: %v", err14)
	}
	if err15 != nil {
		log.Fatalf("❌ Invalid lesson 15 UUID: %v", err15)
	}
	if err16 != nil {
		log.Fatalf("❌ Invalid lesson 16 UUID: %v", err16)
	}
	if err17 != nil {
		log.Fatalf("❌ Invalid lesson 17 UUID: %v", err17)
	}
	if err18 != nil {
		log.Fatalf("❌ Invalid lesson 18 UUID: %v", err18)
	}
	if err19 != nil {
		log.Fatalf("❌ Invalid lesson 19 UUID: %v", err19)
	}
	if err20 != nil {
		log.Fatalf("❌ Invalid lesson 20 UUID: %v", err20)
	}
	if err21 != nil {
		log.Fatalf("❌ Invalid lesson 21 UUID: %v", err21)
	}
	if err22 != nil {
		log.Fatalf("❌ Invalid lesson 22 UUID: %v", err22)
	}
	if err23 != nil {
		log.Fatalf("❌ Invalid lesson 23 UUID: %v", err23)
	}
	if err24 != nil {
		log.Fatalf("❌ Invalid lesson 24 UUID: %v", err24)
	}
	if err25 != nil {
		log.Fatalf("❌ Invalid lesson 25 UUID: %v", err25)
	}

		// ✅ Validate UUIDs  Courses
	courseId1, err26 := uuid.Parse("27d8ae14-4311-4380-8397-057ad5043fd6") // Go Lang Course
	courseId2, err27 := uuid.Parse("909b6026-30da-41f7-868f-42e6acba72c3") // React Course
	courseId3, err28 := uuid.Parse("b8ef3c14-d8ef-46fd-b63e-01b50cc9f227") // Python Course
	courseId4, err29 := uuid.Parse("c40c9c00-0779-490a-931b-e8dbd91549bf") // FastAPI Course
	courseId5, err30 := uuid.Parse("cdbe9a63-c659-4912-abb0-58dcb9d2f341") // Data Analysis Course

	if err26 != nil {
		log.Fatalf("❌ Invalid Course 26 UUID: %v", err26)
	}
	if err27 != nil {
		log.Fatalf("❌ Invalid Course 27 UUID: %v", err27)
	}
	if err28 != nil {
		log.Fatalf("❌ Invalid Course 28 UUID: %v", err28)
	}
	if err29 != nil {
		log.Fatalf("❌ Invalid Course 29 UUID: %v", err29)
	}
	if err30 != nil {
		log.Fatalf("❌ Invalid Course 30 UUID: %v", err30)
	}

	// ✅ Validate UUIDs for Go Lang Course
	moduleID31, err31 := uuid.Parse("1fb7c5b1-a438-4a19-9eb7-19a9046d0124")
	moduleID32, err32 := uuid.Parse("212021f5-c570-4898-b93d-129e4478497a")
	moduleID33, err33 := uuid.Parse("77af55f5-247c-44c7-b34d-67d118b95a26")
	moduleID34, err34 := uuid.Parse("b2df713a-0ca7-4fb3-9d43-d109b09d1491")
	moduleID35, err35 := uuid.Parse("ec980f1d-2ade-4e92-bb9b-912fd20c8f47")

	// ✅ Validate UUIDs for React
	moduleID36, err36 := uuid.Parse("1cee4f93-df85-4fb0-8557-c9773872e66a")
	moduleID37, err37 := uuid.Parse("79879e71-61fd-407d-8ccf-c3b4f813f9a0")
	moduleID38, err38 := uuid.Parse("7a4c17e4-3c69-433c-ad2d-f19e4ef01df3")
	moduleID39, err39 := uuid.Parse("bd26f9a3-4574-4fdb-9491-fb2be51927fb")
	moduleID40, err40 := uuid.Parse("bfe47ad2-49eb-4ed3-910f-3d22cb0e3e23")

	// ✅ Validate UUIDs for Python
	moduleID41, err41 := uuid.Parse("04936943-eaac-473c-a64d-3ac554465ea3")
	moduleID42, err42 := uuid.Parse("0beb9581-05ca-48f5-835c-4df7abefbb8f")
	moduleID43, err43 := uuid.Parse("353a4e83-e5b9-48f9-a7f7-d7ba351d6185")
	moduleID44, err44 := uuid.Parse("71c88341-9d44-4959-a272-5ea42a5f1503")
	moduleID45, err45 := uuid.Parse("96a5d5b4-abf0-46e3-b0c4-b35629966771")

	// ✅ Validate UUIDs for FastAPI
	moduleID46, err46 := uuid.Parse("2a2b4068-b6c5-4272-96c2-b0a76c7970eb")
	moduleID47, err47 := uuid.Parse("c94b72af-e65b-4203-91b8-c0fcdfaa44d5")
	moduleID48, err48 := uuid.Parse("dcd8dfc5-3e3b-4274-b442-6f0b1455268b")
	moduleID49, err49 := uuid.Parse("f2d65324-0e9d-42fe-8b50-26b67c1134d5")
	moduleID50, err50 := uuid.Parse("fe0d1f0d-3c07-4c96-a60a-80102aa1b6de")
	// ✅ Validate UUIDs for Machine Learning
	moduleID51, err51 := uuid.Parse("1f64e54d-0329-4cb3-8e4a-4199d7131723")
	moduleID52, err52 := uuid.Parse("490585ae-2c27-4295-be56-9ff5014e716c")
	moduleID53, err53 := uuid.Parse("89ab26e4-ff4d-484b-afcc-7c0a8c77ff7a")
	moduleID54, err54 := uuid.Parse("b9228452-d56a-4c98-9252-55280459326c")
	moduleID55, err55 := uuid.Parse("fb2aca67-4b51-40a2-897f-030b1b2b8af8")
	
	if err31 != nil {
		log.Fatalf("❌ Invalid module 31 UUID: %v", err31)
	}
	if err32 != nil {
		log.Fatalf("❌ Invalid module 32 UUID: %v", err32)
	}
	if err33 != nil {
		log.Fatalf("❌ Invalid module 33 UUID: %v", err33)
	}
	if err34 != nil {
		log.Fatalf("❌ Invalid module 34 UUID: %v", err34)
	}
	if err35 != nil {
		log.Fatalf("❌ Invalid module 35 UUID: %v", err35)
	}
	if err36 != nil {
		log.Fatalf("❌ Invalid module 36 UUID: %v", err36)
	}
	if err37 != nil {
		log.Fatalf("❌ Invalid module 37 UUID: %v", err37)
	}
	if err38 != nil {
		log.Fatalf("❌ Invalid module 38 UUID: %v", err38)
	}
	if err39 != nil {
		log.Fatalf("❌ Invalid module 39 UUID: %v", err39)
	}
	if err40 != nil {
		log.Fatalf("❌ Invalid module 40 UUID: %v", err40)
	}
	if err41 != nil {
		log.Fatalf("❌ Invalid module 41 UUID: %v", err41)
	}
	if err42 != nil {
		log.Fatalf("❌ Invalid module 42 UUID: %v", err42)
	}
	if err43 != nil {
		log.Fatalf("❌ Invalid module 43 UUID: %v", err43)
	}
	if err44 != nil {
		log.Fatalf("❌ Invalid module 44 UUID: %v", err44)
	}
	if err45 != nil {
		log.Fatalf("❌ Invalid module 45 UUID: %v", err45)
	}
	if err46 != nil {
		log.Fatalf("❌ Invalid module 46 UUID: %v", err46)
	}
	if err47 != nil {
		log.Fatalf("❌ Invalid module 47 UUID: %v", err47)
	}
	if err48 != nil {
		log.Fatalf("❌ Invalid module 48 UUID: %v", err48)
	}
	if err49 != nil {
		log.Fatalf("❌ Invalid module 49 UUID: %v", err49)
	}
	if err50 != nil {
		log.Fatalf("❌ Invalid module 50 UUID: %v", err50)
	}
	if err51 != nil {
		log.Fatalf("❌ Invalid module 51 UUID: %v", err51)
	}
	if err52 != nil {
		log.Fatalf("❌ Invalid module 52 UUID: %v", err52)
	}
	if err53 != nil {
		log.Fatalf("❌ Invalid module 53 UUID: %v", err53)
	}
	if err54 != nil {
		log.Fatalf("❌ Invalid module 54 UUID: %v", err54)
	}
	if err55 != nil {
		log.Fatalf("❌ Invalid module 55 UUID: %v", err55)	
	}



	createdByID, err2 := uuid.Parse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	if err2 != nil {
		log.Fatalf("❌ Invalid CreatedBy UUID: %v", err2)
	}



	lessons := []models.Topics{

		// Go Lang Course
		{
			ID:        topicID1,
			LessonID:        lessonID1,
			CourseID: courseId1,
			ModuleID:   moduleID31,
			TutorID:  createdByID,
			Title:     "What is Go?",
			Description:     "This lesson explains the history of Go, why it was created, and the key advantages it provides for modern software development. Students will learn about Go’s simplicity, performance, and built-in support for concurrency. By the end of this lesson, learners will understand where Go is commonly used, including backend services, cloud infrastructure, and distributed systems.",
			Order:      1,
		},
		{
			ID:        lessonID2,
			CourseID: courseId1,
			ModuleID:   moduleID32,
			TutorID:  createdByID,
			Title:     "Installing Go",
			Description:     "Students will learn how to download and install Go on their computer and verify the installation using the Go command line tool.",
			Order:      2,
		},
			{
			ID:        lessonID3,
			CourseID: courseId1,
			ModuleID:   moduleID33,
			TutorID:  createdByID,
			Title:     "Working with Variables",
			Description:     "Students will learn how to declare variables and use Go's built-in data types to store information. This lesson covers variable declaration, initialization, and the different data types available in Go. By the end of this lesson, learners will be able to create and manipulate variables to hold data in their Go programs.",
			Order:      3,
		},
		{
			ID:        lessonID4,
			CourseID: courseId1,
			ModuleID:   moduleID34,
			TutorID:  createdByID,
			Title:     "Conditional Statements",
			Description:     "Students will learn how to use conditional statements to make decisions within a Go program. This includes if statements, else if statements, and switch statements. By the end of this lesson, learners will be able to control the flow of their Go programs based on different conditions.",
			Order:      4,
		},
		{
			ID:        lessonID5,
			CourseID: courseId1,
			ModuleID:   moduleID35,
			TutorID:  createdByID,
			Title:     "Creating Functions",
			Description:     "Students will learn how to define functions, pass parameters, and return values in Go. This lesson covers the basics of function syntax and how to use functions to organize code and promote reusability. By the end of this lesson, learners will be able to create their own functions to perform specific tasks within their Go programs.",
			Order:      5,
		},

			// React Course
		{
			ID:        lessonID6,
			CourseID: courseId2,
			ModuleID:   moduleID36,
			TutorID:  createdByID,
			Title:     "What is React?",
			Description:     "This lesson explains the fundamentals of React, how it works, and why developers use it to build modern web applications.",
			Order:      1,
		},
		{
			ID:        lessonID7,
			CourseID: courseId2,
			ModuleID:   moduleID37,
			TutorID:  createdByID,
			Title:     "Creating Your First React App",
			Description:     "Students will learn how to set up a new React application and understand the basic folder structure of a React project.",
			Order:      2,
		},
			{
			ID:        lessonID8,
			CourseID: courseId2,
			ModuleID:   moduleID38,
			TutorID:  createdByID,
			Title:     "Creating Components",
			Description:     "Students will learn how to create functional components and structure their React applications using reusable UI pieces.",
			Order:      3,
		},
		{
			ID:        lessonID9,
			CourseID: courseId2,
			ModuleID:   moduleID39,
			TutorID:  createdByID,
			Title:     "Understanding Props",
			Description:     "Students will learn how to pass data from one component to another using props.",
			Order:      4,
		},
		{
			ID:        lessonID10,
			CourseID: courseId2,
			ModuleID:   moduleID40,
			TutorID:  createdByID,
			Title:     "React Event Handling",
			Description:     "Students will learn how to handle user actions like button clicks and input changes in React components.",
			Order:      5,
		},

		// Python Course
		{
			ID:        lessonID11,
			CourseID: courseId3,
			ModuleID:   moduleID41,
			TutorID:  createdByID,
			Title:     "What is Python?",
			Description:     "This lesson explains the basics of Python, its history, and why developers choose it for building applications.",
			Order:      1,
		},
		{
			ID:        lessonID12,
			CourseID: courseId3,
			ModuleID:   moduleID42,
			TutorID:  createdByID,
			Title:     "Installing Python",
			Description:     "Students will learn how to install Python on their computer and verify that it is working correctly.",
			Order:      2,
		},
			{
			ID:        lessonID13,
			CourseID: courseId3,
			ModuleID:   moduleID43,
			TutorID:  createdByID,
			Title:     "Working with Variables",
			Description:     "Students will learn how to create variables and assign values in Python programs.",
			Order:      3,
		},
		{
			ID:        lessonID14,
			CourseID: courseId3,
			ModuleID:   moduleID44,
			TutorID:  createdByID,
			Title:     "Conditional Statements",
			Description:     "Students will learn how to use if, elif, and else statements to control program logic.",
			Order:      4,
		},
		{
			ID:        lessonID15,
			CourseID: courseId3,
			ModuleID:   moduleID45,
			TutorID:  createdByID,
			Title:     "Creating Functions",
			Description:     "Students will learn how to define and use functions to make Python programs more organized and reusable.",
			Order:      5,
		},

			// Fast API Course
		{
			ID:        lessonID16,
			CourseID: courseId4,
			ModuleID:   moduleID46,
			TutorID:  createdByID,
			Title:     "What is FastAPI?",
			Description:     "Students will learn what FastAPI is, why it is popular for backend development, and how it compares with other Python frameworks.",
			Order:      1,
		},
		{
			ID:        lessonID17,
			CourseID: courseId4,
			ModuleID:   moduleID47,
			TutorID:  createdByID,
			Title:     "Installing FastAPI",
			Description:     "Students will learn how to install FastAPI and run their first FastAPI application.",
			Order:      2,
		},
			{
			ID:        lessonID18,
			CourseID: courseId4,
			ModuleID:   moduleID48,
			TutorID:  createdByID,
			Title:     "Building a Basic API Endpoint",
			Description:     "Students will learn how to create a simple API route that returns data.",
			Order:      3,
		},
		{
			ID:        lessonID19,
			CourseID: courseId4,
			ModuleID:   moduleID49,
			TutorID:  createdByID,
			Title:     "Working with Request Data",
			Description:     "Students will learn how to accept parameters and request data in FastAPI endpoints.",
			Order:      4,
		},
		{
			ID:        lessonID20,
			CourseID: courseId4,
			ModuleID:   moduleID50,
			TutorID:  createdByID,
			Title:     "Creating Data Models",
			Description:     "Students will learn how to define request and response models using Pydantic.",
			Order:      5,
		},
		
			// Data Science Course
		{
			ID:        lessonID21,
			CourseID: courseId5,
			ModuleID:   moduleID51,
			TutorID:  createdByID,
			Title:     "What is Data Analysis?",
			Description:     "Students will learn what data analysis is, why it matters in business and research, and the steps involved in analyzing data effectively.",
			Order:      1,
		},
		{
			ID:        lessonID22,
			CourseID: courseId5,
			ModuleID:   moduleID52,
			TutorID:  createdByID,
			Title:     "Collecting and Cleaning Data",
			Description:     "Students will learn how to import datasets, handle missing values, and clean data for accurate analysis.",
			Order:      2,
		},
			{
			ID:        lessonID23,
			CourseID: courseId5,
			ModuleID:   moduleID53,
			TutorID:  createdByID,
			Title:     "Exploring and Summarizing Data",
			Description:     "Students will learn to calculate basic statistics and use charts to visualize data trends and distributions.",
			Order:      3,
		},
		{
			ID:        lessonID24,
			CourseID: courseId5,
			ModuleID:   moduleID54,
			TutorID:  createdByID,
			Title:     "Using Python Libraries for Data Analysis",
			Description:     "Students will learn how to manipulate datasets using Pandas and perform numerical operations using NumPy.",
			Order:      4,
		},
		{
			ID:        lessonID25,
			CourseID: courseId5,
			ModuleID:   moduleID55,
			TutorID:  createdByID,
			Title:     "Creating Visualizations",
			Description:     "Students will learn how to create charts and graphs using Python libraries to communicate findings effectively.",
			Order:      5,
		},
		
	
		
	}

	for _, lesson := range lessons {
		if err := db.Create(&lesson).Error; err != nil {
			log.Printf("❌ Failed to seed Lesson: %v", err)
		} else {
			log.Printf("✅ Seeded Lesson: %v", &lesson.Title)
		}
	}
}

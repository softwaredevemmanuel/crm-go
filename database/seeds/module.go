package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedModules() {
	db := config.GetDB()
	// ✅ Validate UUIDs for Go Lang
	moduleID1, err1 := uuid.Parse("1fb7c5b1-a438-4a19-9eb7-19a9046d0124")
	moduleID2, err2 := uuid.Parse("212021f5-c570-4898-b93d-129e4478497a")
	moduleID3, err3 := uuid.Parse("77af55f5-247c-44c7-b34d-67d118b95a26")
	moduleID4, err4 := uuid.Parse("b2df713a-0ca7-4fb3-9d43-d109b09d1491")
	moduleID5, err5 := uuid.Parse("ec980f1d-2ade-4e92-bb9b-912fd20c8f47")

	// ✅ Validate UUIDs for React
	moduleID6, err6 := uuid.Parse("1cee4f93-df85-4fb0-8557-c9773872e66a")
	moduleID7, err7 := uuid.Parse("79879e71-61fd-407d-8ccf-c3b4f813f9a0")
	moduleID8, err8 := uuid.Parse("7a4c17e4-3c69-433c-ad2d-f19e4ef01df3")
	moduleID9, err9 := uuid.Parse("bd26f9a3-4574-4fdb-9491-fb2be51927fb")
	moduleID10, err10 := uuid.Parse("bfe47ad2-49eb-4ed3-910f-3d22cb0e3e23")

	// ✅ Validate UUIDs for Python
	moduleID11, err11 := uuid.Parse("04936943-eaac-473c-a64d-3ac554465ea3")
	moduleID12, err12 := uuid.Parse("0beb9581-05ca-48f5-835c-4df7abefbb8f")
	moduleID13, err13 := uuid.Parse("353a4e83-e5b9-48f9-a7f7-d7ba351d6185")
	moduleID14, err14 := uuid.Parse("71c88341-9d44-4959-a272-5ea42a5f1503")
	moduleID15, err15 := uuid.Parse("96a5d5b4-abf0-46e3-b0c4-b35629966771")

	// ✅ Validate UUIDs for Django Web Development
	moduleID16, err16 := uuid.Parse("2a2b4068-b6c5-4272-96c2-b0a76c7970eb")
	moduleID17, err17 := uuid.Parse("c94b72af-e65b-4203-91b8-c0fcdfaa44d5")
	moduleID18, err18 := uuid.Parse("dcd8dfc5-3e3b-4274-b442-6f0b1455268b")
	moduleID19, err19 := uuid.Parse("f2d65324-0e9d-42fe-8b50-26b67c1134d5")
	moduleID20, err20 := uuid.Parse("fe0d1f0d-3c07-4c96-a60a-80102aa1b6de")

	// ✅ Validate UUIDs for Machine Learning
	moduleID21, err21 := uuid.Parse("1f64e54d-0329-4cb3-8e4a-4199d7131723")
	moduleID22, err22 := uuid.Parse("490585ae-2c27-4295-be56-9ff5014e716c")
	moduleID23, err23 := uuid.Parse("89ab26e4-ff4d-484b-afcc-7c0a8c77ff7a")
	moduleID24, err24 := uuid.Parse("b9228452-d56a-4c98-9252-55280459326c")
	moduleID25, err25 := uuid.Parse("fb2aca67-4b51-40a2-897f-030b1b2b8af8")

	if err1 != nil {
		log.Fatalf("❌ Invalid module 1 UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid module 2 UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid module 3 UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid module 4 UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid module 5 UUID: %v", err5)
	}
	if err6 != nil {
		log.Fatalf("❌ Invalid Course 6 UUID: %v", err6)
	}
	if err7 != nil {
		log.Fatalf("❌ Invalid Course 7 UUID: %v", err7)
	}
	if err8 != nil {
		log.Fatalf("❌ Invalid Course 8 UUID: %v", err8)
	}
	if err9 != nil {
		log.Fatalf("❌ Invalid Course 9 UUID: %v", err9)
	}
	if err10 != nil {
		log.Fatalf("❌ Invalid Course 10 UUID: %v", err10)
	}
	if err11 != nil {
		log.Fatalf("❌ Invalid module 11 UUID: %v", err11)
	}
	if err12 != nil {
		log.Fatalf("❌ Invalid module 12 UUID: %v", err12)
	}
	if err13 != nil {
		log.Fatalf("❌ Invalid module 13 UUID: %v", err13)
	}
	if err14 != nil {
		log.Fatalf("❌ Invalid module 14 UUID: %v", err14)
	}
	if err15 != nil {
		log.Fatalf("❌ Invalid module 15 UUID: %v", err15)
	}
	if err16 != nil {
		log.Fatalf("❌ Invalid module 16 UUID: %v", err16)
	}
	if err17 != nil {
		log.Fatalf("❌ Invalid module 17 UUID: %v", err17)
	}
	if err18 != nil {
		log.Fatalf("❌ Invalid module 18 UUID: %v", err18)
	}
	if err19 != nil {
		log.Fatalf("❌ Invalid module 19 UUID: %v", err19)
	}
	if err20 != nil {
		log.Fatalf("❌ Invalid module 20 UUID: %v", err20)
	}
	if err21 != nil {
		log.Fatalf("❌ Invalid module 21 UUID: %v", err21)
	}
	if err22 != nil {
		log.Fatalf("❌ Invalid module 22 UUID: %v", err22)
	}
	if err23 != nil {
		log.Fatalf("❌ Invalid module 23 UUID: %v", err23)
	}
	if err24 != nil {
		log.Fatalf("❌ Invalid module 24 UUID: %v", err24)
	}

	if err25 != nil {
		log.Fatalf("❌ Invalid module 25 UUID: %v", err25)
	}

	// ✅ Validate UUIDs  Courses
	courseId1, err26 := uuid.Parse("27d8ae14-4311-4380-8397-057ad5043fd6") // Go Lang Course
	courseId2, err27 := uuid.Parse("909b6026-30da-41f7-868f-42e6acba72c3") // React Course
	courseId3, err28 := uuid.Parse("b8ef3c14-d8ef-46fd-b63e-01b50cc9f227") // Python Course
	courseId4, err29 := uuid.Parse("c40c9c00-0779-490a-931b-e8dbd91549bf") // Django Web Development Course
	courseId5, err30 := uuid.Parse("cdbe9a63-c659-4912-abb0-58dcb9d2f341") // Machine Learning Course

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

	modules := []models.Module{
		{
			ID:            moduleID1,
			ModuleNumber:  1,
			CourseID:      courseId1,
			Title:         "Introduction to Go",
			Slug:          "introduction-to-go",
			Description:   "This module introduces the Go programming language, its purpose, and why it is widely used for building scalable backend services and cloud applications.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID2,
			ModuleNumber:  2,
			CourseID:      courseId1,
			Title:         "Setting Up Go Environment",
			Slug:          "setting-up-go-environment",
			Description:   "This module teaches learners how to install Go and prepare their development environment for writing Go programs.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID3,
			ModuleNumber:  3,
			CourseID:      courseId1,
			Title:         "Concurrent Programming in Go",
			Slug:          "concurrent-programming-in-go",
			Description:   "This module introduces how Go stores and manages data using variables and built-in data types. It covers basic data types, composite types, and how to work with them in Go.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID4,
			ModuleNumber:  4,
			CourseID:      courseId1,
			Title:         "Control Flow in Go",
			Slug:          "control-flow-in-go",
			Description:   "This module focuses on controlling the execution flow of programs using conditions and loops. It covers if statements, switch cases, and various loop constructs in Go.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID5,
			ModuleNumber:  5,
			CourseID:      courseId1,
			Title:         "Functions in Go",
			Slug:          "functions-in-go",
			Description:   "This module introduces functions and demonstrates how they help organize reusable blocks of code. It covers function declaration, parameters, return values, and how to use functions effectively in Go.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		// React Modules
		{
			ID:            moduleID6,
			ModuleNumber:  1,
			CourseID:      courseId2,
			Title:         "Introduction to React",
			Slug:          "introduction-to-react",
			Description:   "This module introduces React, a popular JavaScript library used for building user interfaces. Students will learn why React is widely used for creating fast and interactive web applications.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID7,
			ModuleNumber:  2,
			CourseID:      courseId2,
			Title:         "Setting Up a React Project",
			Slug:          "setting-up-react-project",
			Description:   "This module teaches students how to create and configure a React project using modern development tools. It covers setting up the project structure, installing dependencies, and running the development server.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID8,
			ModuleNumber:  3,
			CourseID:      courseId2,
			Title:         "React Components and State Management",
			Slug:          "react-components-and-state-management",
			Description:   "This module explains how React applications are built using reusable components. Students will learn about React's component-based architecture, how to create functional and class components, and how to manage state and props effectively.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID9,
			ModuleNumber:  4,
			CourseID:      courseId2,
			Title:         "React Props and State Management",
			Slug:          "react-props-and-state-management",
			Description:   "This module teaches how data is passed and managed inside React components. It covers the concept of props for passing data from parent to child components and state for managing data within a component. Students will learn how to use props and state to create dynamic and interactive user interfaces.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID10,
			ModuleNumber:  5,
			CourseID:      courseId2,
			Title:         "Handling Events in React",
			Slug:          "handling-events-in-react",
			Description:   "This module focuses on how React handles user interactions such as clicks and form submissions.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		// Python Modules
		{
			ID:            moduleID11,
			ModuleNumber:  1,
			CourseID:      courseId3,
			Title:         "Introduction to Python",
			Slug:          "introduction-to-python",
			Description:   "This module introduces Python, one of the most popular and beginner-friendly programming languages. Students will learn what Python is, where it is used, and why it is widely adopted for web development, automation, and data science.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID12,
			ModuleNumber:  2,
			CourseID:      courseId3,
			Title:         "Setting Up Python Environment",
			Slug:          "setting-up-python-environment",
			Description:   "This module teaches learners how to install Python and prepare their computer for writing and running Python programs.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID13,
			ModuleNumber:  3,
			CourseID:      courseId3,
			Title:         "Variables and Data Types in Python",
			Slug:          "variables-and-data-types-in-python",
			Description:   "This module introduces how Python stores information using variables and different data types. It covers basic data types, collections, and how to work with them in Python.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID14,
			ModuleNumber:  4,
			CourseID:      courseId3,
			Title:         "Control Flow in Python",
			Slug:          "control-flow-in-python",
			Description:   "This module explains how Python programs make decisions and repeat tasks using conditions and loops. It covers if statements, for and while loops, and how to use them to control the flow of a Python program.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID15,
			ModuleNumber:  5,
			CourseID:      courseId3,
			Title:         "Functions in Python",
			Slug:          "functions-in-python",
			Description:   "This module introduces functions and shows how they help organize code into reusable blocks. It covers function definition, parameters, return values, and how to use functions effectively in Python.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		// Python FastAPI Modules
		{
			ID:            moduleID16,
			ModuleNumber:  1,
			CourseID:      courseId4,
			Title:         "Introduction to FastAPI",
			Slug:          "introduction-to-fastapi",
			Description:   "This module introduces FastAPI, a modern Python web framework used for building high-performance APIs quickly and efficiently. Students will learn about FastAPI's key features such as automatic data validation, interactive API documentation, and asynchronous programming support. By the end of this module, learners will understand why FastAPI is a great choice for building APIs and how it compares to other Python web frameworks.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID17,
			ModuleNumber:  2,
			CourseID:      courseId4,
			Title:         "Setting Up FastAPI Development Environment",
			Slug:          "setting-up-fastapi-development-environment",
			Description:   "This module teaches learners how to install FastAPI and set up a development environment to build APIs. It covers creating a virtual environment, installing FastAPI and its dependencies, and running a simple FastAPI application. By the end of this module, students will have a working FastAPI development environment ready for building APIs.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID18,
			ModuleNumber:  3,
			CourseID:      courseId4,
			Title:         "Creating Your First API Endpoint",
			Slug:          "creating-your-first-api-endpoint",
			Description:   "This module teaches how to create simple API endpoints using FastAPI. It covers defining routes, handling HTTP methods, and returning responses. Students will learn how to create GET, POST, PUT, and DELETE endpoints to perform basic CRUD operations. By the end of this module, learners will be able to build simple APIs using FastAPI.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID19,
			ModuleNumber:  4,
			CourseID:      courseId4,
			Title:         "Request and Response Handling",
			Slug:          "fastapi-request-response",
			Description:   "This module explains how FastAPI handles incoming requests and sends responses back to clients. It covers request parsing, response formatting, and how to use FastAPI's features to create robust APIs. Students will learn about request bodies, query parameters, path parameters, and how to return different types of responses. By the end of this module, learners will understand the request-response cycle in FastAPI and how to build APIs that effectively handle client interactions.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID20,
			ModuleNumber:  5,
			CourseID:      courseId4,
			Title:         "Data Validation with Pydantic in FastAPI",
			Slug:          "fastapi-data-validation-with-pydantic",
			Description:   "This module introduces Pydantic models and how FastAPI uses them for validating request data. It covers defining Pydantic models, using them to validate incoming data, and how to handle validation errors. Students will learn how to create complex data models with Pydantic and leverage FastAPI's automatic validation features to ensure that API endpoints receive valid data. By the end of this module, learners will be able to implement robust data validation in their FastAPI applications using Pydantic.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		// Data Analysis Modules
		{
			ID:            moduleID21,
			ModuleNumber:  1,
			CourseID:      courseId5,
			Title:         "Introduction to Data Analysis",
			Slug:          "introduction-to-data-analysis",
			Description:   "This module introduces the fundamentals of data analysis, including its importance, applications, and basic workflow for extracting insights from data.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID22,
			ModuleNumber:  2,
			CourseID:      courseId5,
			Title:         "Data Collection and Cleaning",
			Slug:          "data-collection-and-cleaning",
			Description:   "This module focuses on gathering data from various sources and preparing it for analysis by cleaning and transforming it.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID23,
			ModuleNumber:  3,
			CourseID:      courseId5,
			Title:         "Exploratory Data Analysis (EDA)",
			Slug:          "exploratory-data-analysis",
			Description:   "This module teaches how to explore datasets, understand patterns, and summarize key information using descriptive statistics and visualizations.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID24,
			ModuleNumber:  4,
			CourseID:      courseId5,
			Title:         "Data Analysis with Python",
			Slug:          "data-analysis-with-python",
			Description:   "This module introduces Python tools and libraries for analyzing datasets effectively, including Pandas, NumPy, and Matplotlib.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
		{
			ID:            moduleID25,
			ModuleNumber:  5,
			CourseID:      courseId5,
			Title:         "Data Visualization",
			Slug:          "data-visualization",
			Description:   "This module focuses on presenting insights using charts and graphs to make data easy to understand and interpret.",
			Status:        "published",
			EstimatedTime: 30,
			IsFree:        true,
		},
	}

	for _, module := range modules {
		if err := db.Create(&module).Error; err != nil {
			log.Printf("❌ Failed to seed Modules: %v", err)
		} else {
			log.Printf("✅ Seeded Module: %v", &module.Title)
		}
	}
}

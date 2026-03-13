package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedCourseMaterials() {
	db := config.GetDB()

	// ✅ Validate UUIDs for Course Material
	courseMaterial1, err1 := uuid.Parse("daf69939-f822-4171-b2ca-6929006355c2") // Go Lang Basics
	courseMaterial2, err2 := uuid.Parse("23b8b636-2fb7-47a5-878e-0d1476129121") // React Fundamentals
	courseMaterial3, err3 := uuid.Parse("35ef32ad-870d-43ac-8c44-2e4ef4f40808") // python
	courseMaterial4, err4 := uuid.Parse("a15ed5ca-89b2-449d-9e31-2091b7102969") // Django Web Development
	courseMaterial5, err5 := uuid.Parse("42ca53f0-4ac8-4ea5-91ae-5c995622523e") // Machine Learning Algorithms
	if err1 != nil {
		log.Fatalf("❌ Invalid course material 1 UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid course material 2 UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid course material 3 UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid course material 4 UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid course material 5 UUID: %v", err5)
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

	modules := []models.CourseMaterial{
		{
			ID:            courseMaterial1,
			CourseID:      courseId1,
			Title:         "Go Lang Basics",
			Description:   "Introduction to Go programming language, covering syntax, data types, and basic constructs.",
			Type:          "document",
			FileURL:       "https://example.com/go-lang-basics.pdf",
			Status:        "published",
		},
		{
			ID:            courseMaterial2,
			CourseID:      courseId2,
			Title:         "React Fundamentals",
			Description:   "Learn the fundamentals of React, including components, state management, and hooks.",
			Type:          "video",
			FileURL:       "https://example.com/react-fundamentals.mp4",
			Status:        "published",
		},
		{
			ID:            courseMaterial3,
			CourseID:      courseId3,
			Title:         "Python for Data Science",
			Description:   "Explore Python libraries and tools for data science, including NumPy, pandas, and Matplotlib.",
			Type:          "document",
			FileURL:       "https://example.com/python-data-science.pdf",
			Status:        "published",
		},
		{
			ID:            courseMaterial4,
			CourseID:      courseId4,
			Title:         "Django Web Development",
			Description:   "Build web applications using Django, covering models, views, templates, and REST APIs.",
			Type:          "video",
			FileURL:       "https://example.com/django-web-development.mp4",
			Status:        "published",
		},
		{
			ID:            courseMaterial5,
			CourseID:      courseId5,
			Title:         "Machine Learning Algorithms",
			Description:   "An overview of popular machine learning algorithms, including regression, classification, and clustering.",
			Type:          "document",
			FileURL:       "https://example.com/machine-learning-algorithms.pdf",
			Status:        "published",
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

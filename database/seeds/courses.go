package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedCourses() {
	db := config.GetDB()
	tutorID, err := uuid.Parse("5a853260-31fc-44ee-9d69-bb2a2957ba48")

	courseID1, err1 := uuid.Parse("27d8ae14-4311-4380-8397-057ad5043fd6")
	courseID2, err2 := uuid.Parse("909b6026-30da-41f7-868f-42e6acba72c3")
	courseID3, err3 := uuid.Parse("b8ef3c14-d8ef-46fd-b63e-01b50cc9f227")
	courseID4, err4 := uuid.Parse("c40c9c00-0779-490a-931b-e8dbd91549bf")
	courseID5, err5 := uuid.Parse("cdbe9a63-c659-4912-abb0-58dcb9d2f341")

	if err != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err)
	}	
	if err1 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err5)
	}

	courses := []models.Course{
		{
			ID:          courseID1,
			Title:       "Introduction to Golang",
			Description: "Learn the basics of the Go programming language.",
			Image:       "https://example.com/images/golang.png",
			VideoURL:    "https://example.com/videos/golang-intro.mp4",
			TutorID:     tutorID, 
		},
		{
			ID:          courseID2,
			Title:       "Advanced React",
			Description: "Deep dive into advanced patterns in React.js.",
			Image:       "https://example.com/images/react.png",
			VideoURL:    "https://example.com/videos/react-advanced.mp4",
			TutorID:     tutorID,
		},
		{
			ID:          courseID3,
			Title:       "Python for Data Science",
			Description: "Learn Python programming with a focus on data science applications.",
			Image:       "https://example.com/images/python.png",
			VideoURL:    "https://example.com/videos/python-data-science.mp4",
			TutorID:     tutorID,
		},
		{
			ID:          courseID4,
			Title:       "Web Development with Django",
			Description: "Build robust web applications using the Django framework.",
			Image:       "https://example.com/images/django.png",
			VideoURL:    "https://example.com/videos/django-web-dev.mp4",
			TutorID:     tutorID,
		},
		{
			ID:          courseID5,
			Title:       "Machine Learning Basics",
			Description: "An introduction to machine learning concepts and techniques.",
			Image:       "https://example.com/images/machine-learning.png",
			VideoURL:    "https://example.com/videos/machine-learning-basics.mp4",
			TutorID:     tutorID,
		},
		
	}

	for _, course := range courses {
		if err := db.Create(&course).Error; err != nil {
			log.Printf("❌ Failed to seed course: %v", err)
		} else {
			log.Printf("✅ Seeded course: %s", course.Title)
		}
	}
}


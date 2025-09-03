package seeds

import (
	"log"
	"time"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedCourses() {
	db := config.GetDB()
	tutorID1, err1 := uuid.Parse("5a853260-31fc-44ee-9d69-bb2a2957ba48")
	tutorID2, err2 := uuid.Parse("5a853260-31fc-44ee-9d69-bb2a2957ba48")
	if err1 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err1)
	}
	courses := []models.Course{
		{
			ID:          uuid.New(),
			Title:       "Introduction to Golang",
			Description: "Learn the basics of the Go programming language.",
			Image:       "https://example.com/images/golang.png",
			VideoURL:    "https://example.com/videos/golang-intro.mp4",
			TutorID:     tutorID1, // replace with a valid tutor UUID if you already have tutors
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Advanced React",
			Description: "Deep dive into advanced patterns in React.js.",
			Image:       "https://example.com/images/react.png",
			VideoURL:    "https://example.com/videos/react-advanced.mp4",
			TutorID:     tutorID2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          uuid.New(),
			Title:       "Beginners React",
			Description: "Deep dive into biggerners patterns in React.js.",
			Image:       "https://example.com/images/react.png",
			VideoURL:    "https://example.com/videos/react-advanced.mp4",
			TutorID:     tutorID2,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
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


package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
	"gorm.io/datatypes"
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
			ID:               courseID1,
			Title:            "Introduction to Golang",
			Description:      "Learn the basics of the Go programming language.",
			Image:            "https://example.com/images/golang.png",
			VideoURL:         "https://example.com/videos/golang-intro.mp4",
			TutorID:          tutorID,
			LearningOutcomes: datatypes.JSON([]byte(`["Understand Go syntax", "Work with Go routines", "Build a simple web server"]`)),
			Requirements:     datatypes.JSON([]byte(`["Basic programming knowledge", "Familiarity with command line", "Willingness to learn", "No prior Go experience required"]`)),
		},
		{
			ID:               courseID2,
			Title:            "Advanced React",
			Description:      "Deep dive into advanced patterns in React.js.",
			Image:            "https://example.com/images/react.png",
			VideoURL:         "https://example.com/videos/react-advanced.mp4",
			TutorID:          tutorID,
			LearningOutcomes: datatypes.JSON([]byte(`["Understand advanced React patterns", "State management with Redux", "Performance optimization"]`)),
			Requirements:     datatypes.JSON([]byte(`["Basic React knowledge", "Familiarity with JavaScript ES6+", "Understanding of web development concepts", "Experience with building React applications"]`)),
		},
		{
			ID:               courseID3,
			Title:            "Python for Data Science",
			Description:      "Learn Python programming with a focus on data science applications.",
			Image:            "https://example.com/images/python.png",
			VideoURL:         "https://example.com/videos/python-data-science.mp4",
			TutorID:          tutorID,
			LearningOutcomes: datatypes.JSON([]byte(`["Understand Python basics", "Data manipulation with Pandas", "Data visualization with Matplotlib"]`)),
			Requirements:     datatypes.JSON([]byte(`["Basic programming knowledge", "Familiarity with command line", "Willingness to learn", "No prior Python experience required"]`)),
		},
		{
			ID:          courseID4,
			Title:       "Building APIs with FastAPI",
			Description: "Learn how to develop high-performance APIs using Python FastAPI framework.",
			Image:       "https://example.com/images/fastapi.png",
			VideoURL:    "https://example.com/videos/fastapi-web-dev.mp4",
			TutorID:     tutorID,
			LearningOutcomes: datatypes.JSON([]byte(`[ "Set up FastAPI project", "Create API endpoints", "Validate requests with Pydantic", "Handle path and query parameters", "Deploy FastAPI applications"]`)),
			Requirements: datatypes.JSON([]byte(`[ "Basic Python knowledge", "Familiarity with web development concepts", "Understanding of REST APIs", "No prior FastAPI experience required"]`)),
		},
		{
			ID:          courseID5,
			Title:       "Data Analysis with Python",
			Description: "Learn to analyze and visualize data effectively using Python libraries.",
			Image:       "https://example.com/images/data-analysis.png",
			VideoURL:    "https://example.com/videos/data-analysis-basics.mp4",
			TutorID:     tutorID,
			LearningOutcomes: datatypes.JSON([]byte(`[ "Collect and clean datasets", "Perform exploratory data analysis", "Use Python libraries like Pandas and Matplotlib", "Create insightful data visualizations", "Draw actionable conclusions from data"]`)),
			Requirements: datatypes.JSON([]byte(`[ "Basic Python knowledge", "Familiarity with spreadsheets or data concepts", "Willingness to learn Python libraries", "No prior data analysis experience required"]`)),
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

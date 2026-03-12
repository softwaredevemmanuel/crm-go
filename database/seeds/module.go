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
			Description:   "This is the first module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "introduction-to-go",
			Status:        "published",
			Title:         "Introduction to Go",
		},
		{
			ID:            moduleID2,
			ModuleNumber:  2,
			CourseID:      courseId1,
			Description:   "This is the second module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "data-structures-in-go",
			Status:        "published",
			Title:         "Data Structures in Go",
		},
		{
			ID:            moduleID3,
			ModuleNumber:  3,
			CourseID:      courseId1,
			Description:   "This is the third module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "concurrent-programming-in-go",
			Status:        "published",
			Title:         "Concurrent Programming in Go",
		},
		{
			ID:            moduleID4,
			ModuleNumber:  4,
			CourseID:      courseId1,
			Description:   "This is the fourth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "building-a-web-server-with-go",
			Status:        "published",
			Title:         "Building a Web Server with Go",
		},
		{
			ID:            moduleID5,
			ModuleNumber:  5,
			CourseID:      courseId1,
			Description:   "This is the fifth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "testing-and-debugging-in-go",
			Status:        "published",
			Title:         "Testing and Debugging in Go",
		},
		// React Modules
		{
			ID:            moduleID6,
			ModuleNumber:  1,
			CourseID:      courseId2,
			Description:   "This is the first module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "introduction-to-react",
			Status:        "published",
			Title:         "Introduction to React",
		},
		{
			ID:            moduleID7,
			ModuleNumber:  2,
			CourseID:      courseId2,
			Description:   "This is the second module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "data-structures-in-react",
			Status:        "published",
			Title:         "Data Structures in React",
		},
		{
			ID:            moduleID8,
			ModuleNumber:  3,
			CourseID:      courseId2,
			Description:   "This is the third module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "concurrent-programming-in-react",
			Status:        "published",
			Title:         "Concurrent Programming in React",
		},
		{
			ID:            moduleID9,
			ModuleNumber:  4,
			CourseID:      courseId2,
			Description:   "This is the fourth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "building-a-web-server-with-react",
			Status:        "published",
			Title:         "Building a Web Server with React",
		},
		{
			ID:            moduleID10,
			ModuleNumber:  5,
			CourseID:      courseId2,
			Description:   "This is the fifth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "testing-and-debugging-in-react",
			Status:        "published",
			Title:         "Testing and Debugging in React",
		},
		// Python Modules
		{
			ID:            moduleID11,
			ModuleNumber:  1,
			CourseID:      courseId3,
			Description:   "This is the first module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "introduction-to-python",
			Status:        "published",
			Title:         "Introduction to Python",
		},
		{
			ID:            moduleID12,
			ModuleNumber:  2,
			CourseID:      courseId3,
			Description:   "This is the second module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "data-structures-in-python",
			Status:        "published",
			Title:         "Data Structures in Python",
		},
		{
			ID:            moduleID13,
			ModuleNumber:  3,
			CourseID:      courseId3,
			Description:   "This is the third module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "concurrent-programming-in-python",
			Status:        "published",
			Title:         "Concurrent Programming in Python",
		},
		{
			ID:            moduleID14,
			ModuleNumber:  4,
			CourseID:      courseId3,
			Description:   "This is the fourth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "building-a-web-server-with-python",
			Status:        "published",
			Title:         "Building a Web Server with Python",
		},
		{
			ID:            moduleID15,
			ModuleNumber:  5,
			CourseID:      courseId3,
			Description:   "This is the fifth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "testing-and-debugging-in-python",
			Status:        "published",
			Title:         "Testing and Debugging in Python",
		},
		// Django Web Development Modules
		{
			ID:            moduleID16,
			ModuleNumber:  1,
			CourseID:      courseId4,
			Description:   "This is the first module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "introduction-to-web-development",
			Status:        "published",
			Title:         "Introduction to Web Development",
		},
		{
			ID:            moduleID17,
			ModuleNumber:  2,
			CourseID:      courseId4,
			Description:   "This is the second module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "html-and-css-fundamentals",
			Status:        "published",
			Title:         "HTML and CSS Fundamentals",
		},
		{
			ID:            moduleID18,
			ModuleNumber:  3,
			CourseID:      courseId4,
			Description:   "This is the third module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "styling-in-css",
			Status:        "published",
			Title:         "Styling in CSS",
		},
		{
			ID:            moduleID19,
			ModuleNumber:  4,
			CourseID:      courseId4,
			Description:   "This is the fourth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "building-a-div-frame",
			Status:        "published",
			Title:         "Building a div frame",
		},
		{
			ID:            moduleID20,
			ModuleNumber:  5,
			CourseID:      courseId4,
			Description:   "This is the fifth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "testing-and-debugging-in-django",
			Status:        "published",
			Title:         "Testing and Debugging in Django",
		},
		// Machine Learning Modules
		{
			ID:            moduleID21,
			ModuleNumber:  1,
			CourseID:      courseId5,
			Description:   "This is the first module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "introduction-to-machine-learning",
			Status:        "published",
			Title:         "Introduction to Machine Learning",
		},
		{
			ID:            moduleID22,
			ModuleNumber:  2,
			CourseID:      courseId5,
			Description:   "This is the second module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "supervised-learning-algorithms",
			Status:        "published",
			Title:         "Supervised Learning Algorithms",
		},
		{
			ID:            moduleID23,
			ModuleNumber:  3,
			CourseID:      courseId5,
			Description:   "This is the third module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "unsupervised-learning-algorithms",
			Status:        "published",
			Title:         "Unsupervised Learning Algorithms",
		},
		{
			ID:            moduleID24,
			ModuleNumber:  4,
			CourseID:      courseId5,
			Description:   "This is the fourth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "reinforcement-learning-algorithms",
			Status:        "published",
			Title:         "Reinforcement Learning Algorithms",
		},
		{
			ID:            moduleID25,
			ModuleNumber:  5,
			CourseID:      courseId5,
			Description:   "This is the fifth module of the course.",
			EstimatedTime: 30,
			IsFree:        true,
			Slug:          "testing-and-debugging-in-machine-learning",
			Status:        "published",
			Title:         "Testing and Debugging in Machine Learning",
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

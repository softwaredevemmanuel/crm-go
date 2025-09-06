package seeds

import (
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedCourseCategories() {
	db := config.GetDB()
	categoryID1, err1 := uuid.Parse("283eaf1d-bf89-40fe-b363-b119e5815107")
	categoryID2, err2 := uuid.Parse("6109ae19-2a10-43b3-9daf-94c62f133724")
	categoryID3, err3 := uuid.Parse("84bd048a-f146-4191-99ea-5f45ebd551a3")
	categoryID4, err4 := uuid.Parse("d6b4bc52-c5c5-4af1-96f2-8daf6d977984")
	categoryID5, err5 := uuid.Parse("ec532f97-9de7-4fb4-8b93-dff24ca4acbd")

	courseID1, err6 := uuid.Parse("27d8ae14-4311-4380-8397-057ad5043fd6")
	courseID2, err7 := uuid.Parse("909b6026-30da-41f7-868f-42e6acba72c3")
	courseID3, err8 := uuid.Parse("b8ef3c14-d8ef-46fd-b63e-01b50cc9f227")
	courseID4, err9 := uuid.Parse("c40c9c00-0779-490a-931b-e8dbd91549bf")
	courseID5, err10 := uuid.Parse("cdbe9a63-c659-4912-abb0-58dcb9d2f341")
	
	courseCategoryID1, err11 := uuid.Parse("577d2ae4-9f46-411f-9e2c-21bf24bf1415")
	courseCategoryID2, err12 := uuid.Parse("6276970a-c708-4b3f-a79e-49305ebe5eae")
	courseCategoryID3, err13 := uuid.Parse("7053ace4-7061-4878-9884-9d895990c050")
	courseCategoryID4, err14 := uuid.Parse("b8236077-885f-4bc4-8433-eb017931a005")
	courseCategoryID5, err15 := uuid.Parse("bc68642d-6eb3-448b-9644-24499d7037a7")
	courseCategoryID6, err16 := uuid.Parse("bf7ac956-a842-402b-89f5-854e748599e6")
	courseCategoryID7, err17 := uuid.Parse("cfe2cefe-197a-4442-9618-545ab04acd4a")
	courseCategoryID8, err18 := uuid.Parse("def8f4d1-c8d1-45f3-93c2-7f6875974e0b")
	courseCategoryID9, err19 := uuid.Parse("e75e9747-89b0-4125-b473-0f5a532495f4")
	courseCategoryID10, err20 := uuid.Parse("eb915723-8fea-48c7-aee5-813be2ed0a2f")

if err1 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err3)
	}
	if err4 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err4)
	}
	if err5 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err5)
	}
	if err6 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err6)
	}
	if err7 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err7)
	}
	if err8 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err8)
	}
	if err9 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err9)
	}
	if err10 != nil {
		log.Fatalf("❌ Invalid course UUID: %v", err10)
	}
if err11 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err11)
	}
	if err12 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err12)
	}
	if err13 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err13)
	}
	if err14 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err14)
	}
	if err15 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err15)
	}
	if err16 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err16)
	}
	if err17 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err17)
	}
	if err18 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err18)
	}
	if err19 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err19)
	}
	if err20 != nil {
		log.Fatalf("❌ Invalid course category UUID: %v", err20)
	}



	course_categories := []models.CourseCategoryTable{
			{
		ID: 	  courseCategoryID1,
		CourseID: courseID1,
		CategoryID: categoryID3,
		},
		{
		ID: 	  courseCategoryID2,
		CourseID: courseID2,
		CategoryID: categoryID4,	
		},
		{
		ID: 	  courseCategoryID3,
		CourseID: courseID3,
		CategoryID: categoryID1,	
		},
		{
		ID: 	  courseCategoryID4,
		CourseID: courseID4,
		CategoryID: categoryID2,	
		},
		{
		ID: 	  courseCategoryID5,
		CourseID: courseID5,
		CategoryID: categoryID5,	
		},	
		{
		ID: 	  courseCategoryID6,
		CourseID: courseID1,
		CategoryID: categoryID1,
		},	
		{
		ID: 	  courseCategoryID7,
		CourseID: courseID2,
		CategoryID: categoryID2,
		},
		{
		ID: 	  courseCategoryID8,
		CourseID: courseID3,
		CategoryID: categoryID3,
		},	
		{
		ID: 	  courseCategoryID9,
		CourseID: courseID4,
		CategoryID: categoryID4,
		},
		{
		ID: 	  courseCategoryID10,
		CourseID: courseID5,
		CategoryID: categoryID5,
		},	

		
		
	}

	

	for _, course_category := range course_categories {
		if err := db.Create(&course_category).Error; err != nil {
			log.Printf("❌ Failed to seed Course Category: %v", err)
		} else {
			log.Printf("✅ Course Category Seeded Successfully:")
		}
	}
}


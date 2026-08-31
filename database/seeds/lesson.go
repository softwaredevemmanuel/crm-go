package seeds

import (
	"fmt"
	"log"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
)

func SeedLessons() error {
	db := config.GetDB()

	
	// SS1 Class
	classID := uuid.MustParse("9f7c2a41-6b83-4d95-a1e7-2c8f5b3d9046")

	// SS1 Arm
	armID := uuid.MustParse("73d9b526-1a84-4e60-8c35-617f2a948bcd")

	// Biology SS1 First Term Scheme of Work Items

	item1 := uuid.MustParse("8f3a1c72-5d94-4b68-a021-739e5f2c8146")

	item2 := uuid.MustParse("c4e7b912-3a56-4d80-b135-628f9a1e4732")

	item3 := uuid.MustParse("1d8b6f43-9275-4ca1-a850-316e7d5942bc")

	item4 := uuid.MustParse("7a2e9c51-4f83-46bd-b027-915c6a3e7481")

	item5 := uuid.MustParse("b6d4f820-1c59-4e73-9a26-583b7d14e9c0")

	item6 := uuid.MustParse("3e91a745-6b28-4c50-b839-172d5f6a904e")

	item7 := uuid.MustParse("f2c8d631-9a47-45be-a103-764e2b9158cd")

	item8 := uuid.MustParse("5b7e4a29-83d1-4f60-b925-318c6e7a5042")

	item9 := uuid.MustParse("a9c3f815-2e64-4b71-8d30-547a1f96c2be")

	item10 := uuid.MustParse("6d18b953-7c42-4ea0-a615-829f3b4c7061")

	item11 := uuid.MustParse("e5a7c240-1d83-4f96-b528-639e2a7154cd")

	item12 := uuid.MustParse("42f9b681-5c37-4ad2-8e10-753c6b9a214f")

	item13 := uuid.MustParse("9c6e2a57-4b91-46d3-a820-185f7c3946be")

	lessonDate := ParseDate("2026-09-01")

	lessons := []models.Lesson{
		{
			ID:                 uuid.MustParse("e2a83449-fd98-4b24-8fcf-e51b4e17456b"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item1,
			Status:             "scheduled",
			Title:              "Introduction to Micro-organisms",
			Week:               1,
		},
		{
			ID:                 uuid.MustParse("9e12aaf1-fdb9-430e-8f5e-e98c36ec3bf9"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item2,
			Status:             "scheduled",
			Title:              "Culturing and Identification of Micro-organisms",
			Week:               2,
		},
		{
			ID:                 uuid.MustParse("d837b61a-d747-4cb9-bca8-c165bc81e6fe"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item3,
			Status:             "scheduled",
			Title:              "Sexually Transmitted Infections",
			Week:               3,
		},
		{
			ID:                 uuid.MustParse("caafaae2-7efa-4cd8-89fb-a893c35a88e2"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item4,
			Status:             "scheduled",
			Title:              "Micro-organisms in Action",
			Week:               4,
		},
		{
			ID:                 uuid.MustParse("3f93af4b-9667-4772-8ada-b53875136304"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item5,
			Status:             "scheduled",
			Title:              "Causes, Effects and Prevention of STIs",
			Week:               5,
		},
		{
			ID:                 uuid.MustParse("710f0501-9c6b-49f7-b8b9-5b1ef22f519f"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item6,
			Status:             "scheduled",
			Title:              "Control of Harmful Micro-organisms",
			Week:               6,
		},
		{
			ID:                 uuid.MustParse("90169597-ef87-4451-8d79-e6f3451b344d"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item7,
			Status:             "scheduled",
			Title:              "Classification of Plants",
			Week:               7,
		},
		{
			ID:                 uuid.MustParse("e6492333-05af-4d29-8e38-8f83d050304d"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item8,
			Status:             "scheduled",
			Title:              "Plant Pests and Diseases",
			Week:               8,
		},
		{
			ID:                 uuid.MustParse("ae02e07e-e0f9-4d58-aa29-1f7f1d35e382"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item9,
			Status:             "scheduled",
			Title:              "Animal Pests and Diseases",
			Week:               9,
		},
		{
			ID:                 uuid.MustParse("ebe29753-3c96-4a0b-a7c0-5539d4d63a0b"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item10,
			Status:             "scheduled",
			Title:              "Control of Animal Diseases",
			Week:               10,
		},
		{
			ID:                 uuid.MustParse("e4698004-cad7-4c38-979c-4d0e78101bc6"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item11,
			Status:             "scheduled",
			Title:              "Food Production",
			Week:               11,
		},
		{
			ID:                 uuid.MustParse("3c8f106c-0203-4c69-bfb2-f507b1fc5819"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item12,
			Status:             "scheduled",
			Title:              "Food Preservation and Storage",
			Week:               12,
		},
		{
			ID:                 uuid.MustParse("9983316d-d951-4ba5-8959-90faa942dd49"),
			ArmID:              armID,
			ClassID:            classID,
			Duration:           60,
			LessonDate:         &lessonDate,
			Period:             1,
			SchemeOfWorkItemID: item13,
			Status:             "scheduled",
			Title:              "Revision of Term's Work",
			Week:               13,
		},
	}

	for _, lesson := range lessons {
		if err := db.Create(&lesson).Error; err != nil {
			return fmt.Errorf(
				"failed to seed lesson %s: %w",
				lesson.Title,
				err,
			)
		}

		log.Printf(
			"✅ Seeded lesson: Week %d - %s",
			lesson.Week,
			lesson.Title,
		)
	}

	return nil
}



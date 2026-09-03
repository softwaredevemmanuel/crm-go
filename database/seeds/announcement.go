package seeds

import (
	"log"
	"time"

	"crm-go/config"
	"crm-go/models"
	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func SeedAnnouncements() error {
	db := config.GetDB()
	announcementID1, err1 := uuid.Parse("1e1cbb15-07b6-4f75-9918-7d2e982c8133")
	announcementID2, err2 := uuid.Parse("22629645-9649-4bb3-b9f3-eed88ba36d60")
	announcementID3, err3 := uuid.Parse("2d71fe16-bf9f-4f05-9b5a-618197bd764e")
	announcementID4, err4 := uuid.Parse("8f36b5ea-72c8-479e-8b6a-cdb64ecddd4a")
	announcementID5, err5 := uuid.Parse("9d11215e-5c9d-43e5-834a-04a732c1225d")
	announcementID6, err6 := uuid.Parse("c9df9148-2617-43c8-98a3-24c2c1134613")

	if err1 != nil {
		log.Fatalf("❌ Invalid category UUID: %v", err1)
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
		log.Fatalf("❌ Invalid category UUID: %v", err6)
	}

	createdByID, err2 := uuid.Parse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")

	if err2 != nil {
		log.Fatalf("❌ Invalid CreatedBy UUID: %v", err2)
	}

	endDate, _ := time.Parse(time.RFC3339, "2026-01-22T04:00:00Z")
	startDate, _ := time.Parse(time.RFC3339, "2026-01-22T02:00:00Z")

	announcements := []models.Announcement{
		{
			ID:        announcementID1,
			Audience:  "all",
			CreatedBy: createdByID,
			EndDate:   &endDate,
			IsPinned:  true,
			Message:   "The platform will be unavailable from 2AM to 4AM.",
			StartDate: &startDate,
			Title:     "System Maintenance",
			Type:      "maintenance",
		},
		{
			ID:        announcementID2,
			Audience:  "all",
			CreatedBy: createdByID,
			EndDate:   &endDate,
			IsPinned:  true,
			Message:   "There would be a general meeting at 4:00 PM.",
			StartDate: &startDate,
			Title:     "System Maintenance",
			Type:      "maintenance",
		},
		{
			ID:        announcementID3,
			Audience:  "all",
			CreatedBy: createdByID,
			EndDate:   &endDate,
			IsPinned:  true,
			Message:   "All students are required to attend the meeting.",
			StartDate: &startDate,
			Title:     "System Maintenance",
			Type:      "maintenance",
		},
		{
			ID:        announcementID4,
			Audience:  "all",
			CreatedBy: createdByID,
			EndDate:   &endDate,
			IsPinned:  true,
			Message:   "Phones and Laptops must be returned by the end of the semester.",
			StartDate: &startDate,
			Title:     "System Maintenance",
			Type:      "maintenance",
		},
		{
			ID:        announcementID5,
			Audience:  "all",
			CreatedBy: createdByID,
			EndDate:   &endDate,
			IsPinned:  true,
			Message:   "The CEO will be attending the meeting.",
			StartDate: &startDate,
			Title:     "System Maintenance",
			Type:      "maintenance",
		},
		{
			ID:        announcementID6,
			Audience:  "all",
			CreatedBy: createdByID,
			EndDate:   &endDate,
			IsPinned:  true,
			Message:   "Hello, everyone! The platform will be undergoing scheduled maintenance from 2 AM to 4 AM on January 22, 2026. During this time, the platform will be unavailable. We apologize for any inconvenience this may cause and appreciate your understanding as we work to improve our services.",
			StartDate: &startDate,
			Title:     "System Maintenance",
			Type:      "maintenance",
		},
	}

	for _, announcement := range announcements {
		result := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "id"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"audience",
				"created_by",
				"end_date",
				"is_pinned",
				"message",
				"start_date",
				"title",
				"type",
			}),
		}).Create(&announcement)

		if result.Error != nil {
			log.Printf(
				"❌ Failed to seed Announcement: %v",
				result.Error,
			)
		} else {
			log.Printf(
				"✅ Seeded/updated Announcement: %s",
				announcement.Message,
			)
		}
	}
	return nil
}

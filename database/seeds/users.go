package seeds

import (
	"log"
	"time"

	"crm-go/config"
	"crm-go/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm/clause"

)
func hashPassword(password string) string {
    hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        panic("Failed to hash password: " + err.Error())
    }
    return string(hashed)
}

func SeedUsers() {
	db := config.GetDB()

	userID1, err1 := uuid.Parse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")
	userID2, err2:= uuid.Parse("5a853260-31fc-44ee-9d69-bb2a2957ba48")
	userID3, err3 := uuid.Parse("9c47dbea-5c34-4a35-9084-148c363eddaf")
	if err1 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err1)
	}
	if err2 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err2)
	}
	if err3 != nil {
		log.Fatalf("❌ Invalid tutor UUID: %v", err3)
	}

	
	users := []models.User{
		{
			ID:          userID1,
			FirstName:   "Emmanuel",
			LastName: 	"Okereke",
			Email:       "eokereke47@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "admin",
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID2,
			FirstName:   "Nathan",
			LastName: 	"Chigoziem",
			Email:       "nathan47@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "tutor",
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID3,
			FirstName:   "Caleb",
			LastName: 	"Kachimside",
			Email:       "caleb47@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "student",
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		
	}

	for _, user := range users {
    err := db.Clauses(clause.OnConflict{
        Columns:   []clause.Column{{Name: "email"}}, // avoid duplicates by email
        DoNothing: true, // don’t insert if exists
    }).Create(&user).Error

    if err != nil {
        log.Printf("❌ Failed to seed user: %v", err)
    
	} else {
			log.Printf("✅ Seeded course: %s", user.FirstName + " " + user.LastName)
		}
}
}
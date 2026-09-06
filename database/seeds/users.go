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


func parseDOB(date string) *time.Time {
	if date == "" {
		return nil
	}

	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Printf("Invalid DOB: %s", date)
		return nil
	}

	return &d
}

func SeedUsers() error {
	db := config.GetDB()
	userID1 := uuid.MustParse("fe4547a7-4c81-4bc2-bc81-5bbbce2fb5bd")
	userID2 := uuid.MustParse("5a853260-31fc-44ee-9d69-bb2a2957ba48")
	userID3 := uuid.MustParse("9c47dbea-5c34-4a35-9084-148c363eddaf")
	userID4 := uuid.MustParse("0affa19c-9419-4712-ba3e-9b17bc21c71c")
	userID5 := uuid.MustParse("6182a0b3-dcf1-4f0a-b73e-3c5cdafde0cf")
	userID6 := uuid.MustParse("8f965d41-2072-41e1-ad76-9bb579a9130d")
	userID7 := uuid.MustParse("618da138-da85-47c3-a96c-60a248864da1")
	userID8 := uuid.MustParse("a11962e0-870c-4568-a43c-5f1850b10a95")

	users := []models.User{
		{
			ID:          userID1,
			FirstName:   "Emmanuel",
			LastName: 	"Okereke",
			Email:       "eokereke47@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "admin",
			Position:        "admin",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID2,
			FirstName:   "Nathan",
			LastName: 	"Chigoziem",
			Email:       "upskill@ehizuahub.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "staff",
			Position:    "teacher",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID3,
			FirstName:   "Hannah",
			LastName: 	"Kachimside",
			Email:       "hanniebeke47@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "student",
			Position:     "student",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID4,
			FirstName:   "Mercy",
			LastName: 	"Abeke",
			Email:       "mercyabeke@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "student",
			Position:     "student",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID5,
			FirstName:   "Yemi",
			LastName: 	"Ogumbiyi",
			Email:       "yemi@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "student",
			Position:     "student",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID6,
			FirstName:   "Joseph",
			LastName: 	"Abeke",
			Email:       "Joseph@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "student",
			Position:     "student",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID7,
			FirstName:   "Testimony",
			LastName: 	"Florence",
			Email:       "testimony@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "student",
			Position:     "student",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          userID8,
			FirstName:   "Blessig",
			LastName: 	"Uche",
			Email:       "blessing@gmail.com",
        	Password:  	 hashPassword("mypassword"), 
			LoginID:  "QWERTY",
			Picture:     "https://lh3.googleusercontent.com/a/ACg8ocIucwnbi0gu-NdunUN5er6sqCwOouqNOuQ2dpU-1qR_yH0Kpw=s96-c",
			Role:        "student",
			Position:     "student",
			DOB:       parseDOB("2000-10-08"),
			Provider:   "local",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		
	}

for _, user := range users {
	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "email"},
		},
		DoUpdates: clause.AssignmentColumns([]string{
			"first_name",
			"last_name",
			"password",
			"login_id",
			"picture",
			"role",
			"position",
			"dob",
			"provider",
			"created_at",
			"updated_at",
		}),
	}).Create(&user)

	if result.Error != nil {
		log.Printf(
			"❌ Failed to seed user %s: %v",
			user.Email,
			result.Error,
		)
	} else {
		log.Printf(
			"✅ Seeded/updated user: %s",
			user.FirstName+" "+user.LastName,
		)
	}
}
return nil
}
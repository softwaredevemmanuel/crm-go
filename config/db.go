package config

import (
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
    dsn := os.Getenv("DATABASE_URL") // e.g. "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
    if dsn == "" {
        log.Fatal("❌ DATABASE_URL not set in .env")
    }

    database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("❌ Failed to connect database: %v", err)
    }
    DB= database
    log.Println("✅ Database connected")
    
}



func GetDB() *gorm.DB {
	if DB == nil {
		log.Println("⚠️ Database not initialized")
	}
	return DB
}


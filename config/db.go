package config

import (
    "log"
    "os"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
    dsn := os.Getenv("DATABASE_URL") // e.g. "host=localhost user=postgres password=postgres dbname=mydb port=5432 sslmode=disable"
    if dsn == "" {
        log.Fatal("❌ DATABASE_URL not set in .env")
    }

    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatalf("❌ Failed to connect database: %v", err)
    }

    log.Println("✅ Database connected")
    DB = db
}

func GetDB() *gorm.DB {
    return DB
}


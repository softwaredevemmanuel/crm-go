package config

import (
    "log"
    "os"
    "time"

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

    // Get the underlying *sql.DB to configure the connection pool
	sqlDB, err := database.DB()
	if err != nil {
		log.Fatalf("❌ Failed to get database instance: %v", err)
	}

	// Configure connection pool
	sqlDB.SetMaxOpenConns(20)             // Maximum number of open connections
	sqlDB.SetMaxIdleConns(5)              // Maximum number of idle connections
	sqlDB.SetConnMaxLifetime(time.Hour)   // Recreate connections after 1 hour
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	// Verify the connection
	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("❌ Database ping failed: %v", err)
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


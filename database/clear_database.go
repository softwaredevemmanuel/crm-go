package database

import (
	"crm-go/config"
	"fmt"
	"log"

)

func ClearDatabase() error {
	db := config.DB

	// Disable FK constraints temporarily (important for truncate)
	if err := db.Exec("SET session_replication_role = 'replica';").Error; err != nil {
		return err
	}

	// List all tables (you can filter if you want to keep some)
	var tables []string
	if err := db.Raw(`
        SELECT tablename 
        FROM pg_tables 
        WHERE schemaname = 'public';
    `).Scan(&tables).Error; err != nil {
		return err
	}

	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE;", table)
		if err := db.Exec(query).Error; err != nil {
			log.Printf("❌ Failed truncating %s: %v\n", table, err)
		} else {
			log.Printf("✅ Cleared table: %s\n", table)
		}
	}

	// Re-enable FK constraints
	if err := db.Exec("SET session_replication_role = 'origin';").Error; err != nil {
		return err
	}

	return nil
}

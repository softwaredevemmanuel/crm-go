package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// JWT
	JWTSecret string
	JWTExpire int

	// SMTP
	SMTPServer   string
	SMTPPort     int
	SMTPLogin    string
	SMTPPassword string
	SMTPFrom     string
}

func LoadConfig() *Config {

	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️ No .env file found, falling back to system environment")
	}

	// Parse DB port
	dbPortStr := os.Getenv("DB_PORT")
	if dbPortStr == "" {
		dbPortStr = "5432"
	}
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("❌ Invalid DB_PORT: %v", err)
	}

	// Parse SMTP port
	smtpPortStr := os.Getenv("SMTP_PORT")
	if smtpPortStr == "" {
		smtpPortStr = "587"
	}
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		log.Fatalf("❌ Invalid SMTP_PORT: %v", err)
	}

	// Parse JWT expiration hours
	jwtExpStr := os.Getenv("JWT_EXPIRATION_HOURS")
	if jwtExpStr == "" {
		jwtExpStr = "72"
	}
	jwtExp, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		log.Fatalf("❌ Invalid SMTP_PORT: %v", err)
	}

	return &Config{
		// DB
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "go_crm"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// JWT
		JWTSecret: getEnv("JWT_SECRET", "default_secret"),
		JWTExpire: jwtExp,

		// SMTP
		SMTPServer:   getEnv("SMTP_SERVER", "smtp-relay.brevo.com"),
		SMTPPort:     smtpPort,
		SMTPLogin:    getEnv("SMTP_LOGIN", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("FROM_EMAIL", ""),
	}
}

// Helper
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

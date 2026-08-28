// config/email.go
package config

import (
	"os"
	"fmt"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPLogin    string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

func LoadEmailConfig() EmailConfig {
	return EmailConfig{
		SMTPHost:     getEnv("SMTP_SERVER", "smtp.gmail.com"),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPLogin:    getEnv("SMTP_LOGIN", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@ehizuahub.com"),
		FromName:     getEnv("FROM_NAME", "Ehizua Hub"),
	}
}


func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		var intVal int
		fmt.Sscanf(value, "%d", &intVal)
		return intVal
	}
	return defaultValue
}
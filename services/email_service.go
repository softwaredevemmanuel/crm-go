// services/email_service.go
package services

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/gomail.v2"
)

type EmailService struct {
	dialer *gomail.Dialer
	from   string
}

type EmailData struct {
	To           string
	Subject      string
	Template     string
	TemplateData map[string]interface{}
}

// EmailConfig holds SMTP configuration
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPLogin    string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

// Initialize email config - replace with your values
var emailConfig = EmailConfig{
	SMTPHost:     os.Getenv("SMTP_SERVER"),
	SMTPPort:     465,
	SMTPLogin:    os.Getenv("SMTP_LOGIN"),
	SMTPPassword: os.Getenv("SMTP_PASSWORD"),
	FromEmail:    os.Getenv("FROM_EMAIL"),
	FromName:     os.Getenv("FROM_NAME"),
}

func NewEmailService() *EmailService {
	host := os.Getenv("SMTP_SERVER")
	port := 465
	username := os.Getenv("SMTP_LOGIN")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("FROM_EMAIL")

	if host == "" {
		host = "smtp.gmail.com"
	}
	if username == "" {
		username = os.Getenv("SMTP_LOGIN")
	}
	if password == "" {
		password = os.Getenv("SMTP_PASSWORD")
	}
	if from == "" {
		from = "noreply@edutech.com"
	}

	dialer := gomail.NewDialer(host, port, username, password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &EmailService{
		dialer: dialer,
		from:   from,
	}
}

// loadTemplate loads and parses an HTML template from the templates/emails directory
func loadTemplate(templateName string, data map[string]interface{}) (string, error) {
	// Get the template path
	templatePath := filepath.Join("templates", "emails", templateName)
	
	// Check if the template file exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", templatePath)
	}

	// Parse the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template %s: %w", templatePath, err)
	}

	// Execute the template with the provided data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

// SendEmail sends an email using SMTP with TLS
func SendEmail(emailData EmailData) error {
	if len(emailData.To) == 0 {
		return fmt.Errorf("no recipients provided")
	}

	// Set SMTP config from environment or use defaults
	smtpHost := emailConfig.SMTPHost
	if smtpHost == "" {
		smtpHost = os.Getenv("SMTP_SERVER")
		if smtpHost == "" {
			smtpHost = "smtp.gmail.com"
		}
	}

	smtpPort := emailConfig.SMTPPort
	if smtpPort == 0 {
		smtpPort = 465
	}

	smtpLogin := emailConfig.SMTPLogin
	if smtpLogin == "" {
		smtpLogin = os.Getenv("SMTP_LOGIN")
	}

	smtpPassword := emailConfig.SMTPPassword
	if smtpPassword == "" {
		smtpPassword = os.Getenv("SMTP_PASSWORD")
	}

	fromEmail := emailConfig.FromEmail
	if fromEmail == "" {
		fromEmail = os.Getenv("FROM_EMAIL")
		if fromEmail == "" {
			fromEmail = "noreply@edutech.com"
		}
	}

	fromName := emailConfig.FromName
	if fromName == "" {
		fromName = os.Getenv("FROM_NAME")
		if fromName == "" {
			fromName = "EduTech"
		}
	}

	// Determine body content
	var body string
	var err error

	// If template is specified, load it from the templates/emails directory
	if emailData.Template != "" {
		body, err = loadTemplate(emailData.Template, emailData.TemplateData)
		if err != nil {
			return fmt.Errorf("failed to load template: %w", err)
		}
	} else if emailData.TemplateData != nil {
		// Fallback to Body from TemplateData
		if b, ok := emailData.TemplateData["Body"].(string); ok {
			body = b
		}
	}

	// If no body, use a default message
	if body == "" {
		body = "Email content not available"
	}

	// Determine content type - default to HTML since we're using HTML templates
	contentType := "text/html; charset=UTF-8"
	if emailData.TemplateData != nil {
		if isPlain, ok := emailData.TemplateData["IsPlainText"].(bool); ok && isPlain {
			contentType = "text/plain; charset=UTF-8"
		}
	}

	// Authentication
	auth := smtp.PlainAuth(
		"",
		smtpLogin,
		smtpPassword,
		smtpHost,
	)

	// Build the email message
	message := fmt.Sprintf(
		"From: %s <%s>\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: %s\r\n"+
			"\r\n"+
			"%s",
		fromName,
		fromEmail,
		strings.Join([]string{emailData.To}, ", "),
		emailData.Subject,
		contentType,
		body,
	)

	// Connect to SMTP server with TLS
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)

	tlsConfig := &tls.Config{
		ServerName: smtpHost,
		MinVersion: tls.VersionTLS12,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate with SMTP: %w", err)
	}

	if err := client.Mail(fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %w", err)
	}

	// Add recipients
	recipients := strings.Split(emailData.To, ",")
	for _, recipient := range recipients {
		recipient = strings.TrimSpace(recipient)
		if recipient != "" {
			if err := client.Rcpt(recipient); err != nil {
				return fmt.Errorf("failed to set recipient %s: %w", recipient, err)
			}
		}
	}

	// Send the email body
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start email data: %w", err)
	}

	if _, err := writer.Write([]byte(message)); err != nil {
		writer.Close()
		return fmt.Errorf("failed to write email: %w", err)
	}

	if err := writer.Close(); err != nil {
		return fmt.Errorf("failed to close email data: %w", err)
	}

	if err := client.Quit(); err != nil {
		return fmt.Errorf("failed to quit SMTP connection: %w", err)
	}

	return nil
}

// SendVerificationEmail sends email verification link
func (s *EmailService) SendVerificationEmail(to, firstName, lastName, token, frontendURL string) error {
	verificationURL := fmt.Sprintf("%s/verify-email?token=%s", frontendURL, token)

	emailData := EmailData{
		To:       to,
		Subject:  "Verify Your Email Address - EduTech",
		Template: "verification_email.html",
		TemplateData: map[string]interface{}{
			"FirstName":       firstName,
			"LastName":        lastName,
			"VerificationURL": verificationURL,
			"Token":           token,
			"Email":           to,
			"Year":            2024,
		},
	}

	return SendEmail(emailData)
}

// SendPasswordResetEmail sends password reset link
func (s *EmailService) SendPasswordResetEmail(to, firstName, lastName, token, frontendURL  string) error {
	resetURL := fmt.Sprintf("%s/reset-password?token=%s", frontendURL, token)

	emailData := EmailData{
		To:       to,
		Subject:  "Reset Your Password - EduTech",
		Template: "password_reset_email.html",
		TemplateData: map[string]interface{}{
			"Email":     to,
			"FirstName": firstName,
			"LastName":  lastName,
			"Token":     token,
			"ResetURL":  resetURL,
			"Year":      2024,
		},
	}

	return SendEmail(emailData)
}

// SendEmailWithTemplate sends a custom email with the specified template
func (s *EmailService) SendEmailWithTemplate(to, subject, templateName string, data map[string]interface{}) error {
	emailData := EmailData{
		To:           to,
		Subject:      subject,
		Template:     templateName,
		TemplateData: data,
	}

	return SendEmail(emailData)
}
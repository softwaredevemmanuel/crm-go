// handlers/auth_handler.go
package controllers

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/smtp"
	"html/template"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"crm-go/config"
	"crm-go/models"
)

// SignUpInput represents the signup request body
type SignUpInput struct {
	Email      string `json:"email" binding:"required,email"`
	FirstName  string `json:"first_name" binding:"required"`
	LastName   string `json:"last_name" binding:"required"`
	MiddleName string `json:"middle_name"`
	Phone      string `json:"phone" binding:"required"`
	DOB        string `json:"dob"`
	Password   string `json:"password" binding:"required,min=8"`
	Position   string `json:"position"`
	Role       string `json:"role" binding:"required"`
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

// EmailData represents the data for the email template
type EmailData struct {
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	Password    string
	Role        string
	Position    string
	LoginURL    string
	CurrentYear int
	CompanyName string
}

// Initialize email config - replace with your values
var emailConfig = EmailConfig{
	SMTPHost:     "smtp-relay.brevo.com", // Replace with your SMTP host
	SMTPPort:     587,               // Replace with your SMTP port
	SMTPLogin:    "unitimarket@outlook.com", // Replace with your email
	SMTPPassword: "sOk1YNQdJq3SBHyx", // Replace with your app password
	FromEmail:    "unitimarket@outlook.com",
	FromName:     "Ehizua Hub",
}

// SendEmail sends an email using SMTP
func SendEmail(to []string, subject string, body string, isHTML bool) error {
	// SMTP server configuration
	smtpHost := emailConfig.SMTPHost
	smtpPort := emailConfig.SMTPPort
	smtpLogin := emailConfig.SMTPLogin
	smtpPassword := emailConfig.SMTPPassword
	fromEmail := emailConfig.FromEmail

	// Message headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", emailConfig.FromName, fromEmail)
	headers["To"] = to[0]
	headers["Subject"] = subject
	if isHTML {
		headers["MIME-Version"] = "1.0"
		headers["Content-Type"] = "text/html; charset=\"UTF-8\""
	} else {
		headers["Content-Type"] = "text/plain; charset=\"UTF-8\""
	}

	// Build message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	// Authentication
	auth := smtp.PlainAuth("", smtpLogin, smtpPassword, smtpHost)

	// TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         smtpHost,
	}

	// Connect to SMTP server
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to SMTP server: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return fmt.Errorf("failed to create SMTP client: %v", err)
	}
	defer client.Close()

	// Auth
	if err = client.Auth(auth); err != nil {
		return fmt.Errorf("failed to authenticate: %v", err)
	}

	// Set sender
	if err = client.Mail(fromEmail); err != nil {
		return fmt.Errorf("failed to set sender: %v", err)
	}

	// Set recipient
	if err = client.Rcpt(to[0]); err != nil {
		return fmt.Errorf("failed to set recipient: %v", err)
	}

	// Send message
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("failed to start data: %v", err)
	}
	defer w.Close()

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("failed to write message: %v", err)
	}

	return nil
}

// GenerateWelcomeEmail generates the HTML email template
func GenerateWelcomeEmail(data EmailData) (string, error) {
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Welcome to {{.CompanyName}}</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            line-height: 1.6;
            color: #333;
            max-width: 600px;
            margin: 0 auto;
            padding: 20px;
            background-color: #f8fafc;
        }
        .container {
            background-color: #ffffff;
            border-radius: 12px;
            padding: 40px;
            box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        }
        .header {
            text-align: center;
            padding-bottom: 20px;
            border-bottom: 2px solid #f0f0f0;
            margin-bottom: 30px;
        }
        .header h1 {
            color: #134574;
            font-size: 28px;
            margin: 0;
        }
        .header .subtitle {
            color: #F13178;
            font-size: 16px;
            margin-top: 5px;
        }
        .welcome {
            font-size: 18px;
            margin-bottom: 20px;
        }
        .user-details {
            background-color: #f8fafc;
            border-radius: 8px;
            padding: 20px;
            margin: 20px 0;
            border-left: 4px solid #F13178;
        }
        .user-details p {
            margin: 8px 0;
        }
        .user-details .label {
            font-weight: bold;
            color: #134574;
            display: inline-block;
            width: 120px;
        }
        .credentials-box {
            background-color: #fff3cd;
            border: 1px solid #ffc107;
            border-radius: 8px;
            padding: 15px;
            margin: 20px 0;
        }
        .credentials-box .password {
            font-weight: bold;
            color: #134574;
            font-family: monospace;
            font-size: 16px;
            letter-spacing: 1px;
        }
        .button {
            display: inline-block;
            background-color: #F13178;
            color: white;
            padding: 12px 30px;
            text-decoration: none;
            border-radius: 25px;
            margin: 20px 0;
            font-weight: bold;
        }
        .button:hover {
            background-color: #d42a6a;
        }
        .footer {
            margin-top: 30px;
            padding-top: 20px;
            border-top: 1px solid #f0f0f0;
            text-align: center;
            color: #666;
            font-size: 14px;
        }
        .footer .company {
            color: #134574;
            font-weight: bold;
        }
        .note {
            background-color: #e8f4fd;
            border-radius: 8px;
            padding: 15px;
            margin: 20px 0;
            font-size: 14px;
            color: #0056b3;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>{{.CompanyName}}</h1>
            <div class="subtitle">Welcome to the Family! 🎉</div>
        </div>

        <div class="welcome">
            Dear <strong>{{.FirstName}} {{.LastName}}</strong>,
        </div>

        <p>We are thrilled to welcome you to {{.CompanyName}}! Your account has been successfully created.</p>

        <div class="user-details">
            <h3 style="color: #134574; margin-top: 0;">Account Details</h3>
            <p><span class="label">Name:</span> {{.FirstName}} {{.LastName}}</p>
            <p><span class="label">Email:</span> {{.Email}}</p>
            <p><span class="label">Phone:</span> {{.Phone}}</p>
            <p><span class="label">Role:</span> {{.Role}}</p>
            {{if .Position}}<p><span class="label">Position:</span> {{.Position}}</p>{{end}}
        </div>

        <div class="credentials-box">
            <h4 style="margin-top: 0; color: #856404;">🔑 Your Login Credentials</h4>
            <p><strong>Email:</strong> {{.Email}}</p>
            <p><strong>Password:</strong> <span class="password">{{.Password}}</span></p>
        </div>

        <div class="note">
            <strong>⚠️ Important:</strong> Please change your password after your first login for security purposes.
        </div>

        <div style="text-align: center;">
            <a href="{{.LoginURL}}" class="button">Login to Your Account</a>
        </div>

        <p>If you have any questions or need assistance, please don't hesitate to contact our support team.</p>

        <div class="footer">
            <p>© {{.CurrentYear}} <span class="company">{{.CompanyName}}</span>. All rights reserved.</p>
            <p style="font-size: 12px; color: #999;">This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>
`

	t, err := template.New("welcome").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var body bytes.Buffer
	if err := t.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return body.String(), nil
}

// SignUp handles user registration
// @Summary Register a new user
// @Description Create a new user account with first name, last name, email, password, and role
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   input body SignUpInput true "User signup credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/signup [post]
func SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User with this email already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create user
	user := models.User{
		ID:         uuid.New(),
		FirstName:  input.FirstName,
		LastName:   input.LastName,
		MiddleName: input.MiddleName,
		Email:      input.Email,
		Phone:      input.Phone,
		Password:   string(hashedPassword),
		Role:       input.Role,
		Position:   input.Position,
		IsActive:   true,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Save user to database
	result := config.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	// Send welcome email
	go func() {
		// Prepare email data
		emailData := EmailData{
			FirstName:   user.FirstName,
			LastName:    user.LastName,
			Email:       user.Email,
			Phone:       user.Phone,
			Password:    input.Password, // Send the plain password in the email
			Role:        user.Role,
			Position:    user.Position,
			LoginURL:    "https://your-app-url.com/login", // Replace with your login URL
			CurrentYear: time.Now().Year(),
			CompanyName: "Ehizua Hub",
		}

		// Generate email HTML
		emailBody, err := GenerateWelcomeEmail(emailData)
		if err != nil {
			fmt.Printf("Failed to generate email template: %v\n", err)
			return
		}

		// Send email
		to := []string{user.Email}
		subject := fmt.Sprintf("Welcome to Ehizua Hub, %s!", user.FirstName)

		if err := SendEmail(to, subject, emailBody, true); err != nil {
			fmt.Printf("Failed to send welcome email to %s: %v\n", user.Email, err)
		} else {
			fmt.Printf("Welcome email sent successfully to %s\n", user.Email)
		}
	}()

	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully. A welcome email has been sent.",
		"user": gin.H{
			"id":         user.ID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
			"role":       user.Role,
			"is_active":  user.IsActive,
		},
	})
}

// SendTestEmail is a helper function to test email configuration
// @Summary Send test email
// @Description Send a test email to verify SMTP configuration
// @Tags Authentication
// @Accept  json
// @Produce  json
// @Param   email body object true "Email address to send test email to"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/test-email [post]
func SendTestEmail(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subject := "Test Email from Ehizua Hub"
	body := `
	<html>
	<body>
		<h2>✅ Email Configuration Test</h2>
		<p>This is a test email to verify that the SMTP configuration is working correctly.</p>
		<p>If you received this email, your email setup is working!</p>
		<p>Sent at: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
	</body>
	</html>
	`

	if err := SendEmail([]string{req.Email}, subject, body, true); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send test email: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Test email sent successfully to " + req.Email,
	})
}
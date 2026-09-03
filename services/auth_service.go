// services/auth_service.go
package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"crm-go/dto"
	"crm-go/models"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
)

type AuthService struct {
	db           *gorm.DB
	emailService *EmailService
	frontendURL  string
}

func NewAuthService(db *gorm.DB, emailService *EmailService) *AuthService {
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	return &AuthService{
		db:           db,
		emailService: emailService,
		frontendURL:  frontendURL,
	}
}

// Generate secure random token
func (s *AuthService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Send Verification Email
func (s *AuthService) SendVerificationEmail(req *dto.SendVerificationEmailRequest) error {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ? AND deleted_at IS NULL", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if user is already verified
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// Generate token
	token, err := s.generateToken()
	if err != nil {
		return err
	}

	// Delete any existing verification tokens for this user
	if err := s.db.Where("user_id = ?", user.ID).Delete(&models.EmailVerification{}).Error; err != nil {
		return err
	}

	// Create new verification record
	verification := &models.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour), // 24 hours
		Used:      false,
	}

	if err := s.db.Create(verification).Error; err != nil {
		return err
	}

	// Send email
	if err := s.emailService.SendVerificationEmail(user.Email, user.FirstName, user.LastName, token, s.frontendURL); err != nil {
		return err
	}

	return nil
}

// Verify Email
func (s *AuthService) VerifyEmail(req *dto.VerifyEmailRequest) (*dto.VerifyEmailResponse, error) {
	// Find verification record
	var verification models.EmailVerification
	if err := s.db.Where("token = ? AND used = ?", req.Token, false).
		Preload("User").
		First(&verification).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid or expired verification token")
		}
		return nil, err
	}

	// Check if token expired
	if time.Now().After(verification.ExpiresAt) {
		return nil, errors.New("verification token has expired")
	}

	// Mark token as used
	verification.Used = true
	if err := s.db.Save(&verification).Error; err != nil {
		return nil, err
	}

	// Update user
	now := time.Now()
	user := verification.User
	user.IsVerified = true
	user.EmailVerifiedAt = &now

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &dto.VerifyEmailResponse{
		UserID:     user.ID.String(),
		Email:      user.Email,
		IsVerified: true,
		VerifiedAt: now,
	}, nil
}

// Send Password Reset Email
func (s *AuthService) SendPasswordReset(req *dto.SendPasswordResetRequest) error {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ? AND deleted_at IS NULL", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Don't reveal if user exists or not for security
			return nil
		}
		return err
	}

	// Check if user is active
	if !user.IsActive {
		return errors.New("account is deactivated")
	}

	// Generate token
	token, err := s.generateToken()
	if err != nil {
		return err
	}

	// Delete any existing password reset tokens for this user
	if err := s.db.Where("user_id = ?", user.ID).Delete(&models.PasswordReset{}).Error; err != nil {
		return err
	}

	// Create new password reset record
	reset := &models.PasswordReset{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		ExpiresAt: time.Now().Add(1 * time.Hour), // 1 hour
		Used:      false,
	}

	if err := s.db.Create(reset).Error; err != nil {
		return err
	}

	// Send email
	if err := s.emailService.SendPasswordResetEmail(user.Email, user.FirstName, user.LastName, token, s.frontendURL); err != nil {
		return err
	}

	return nil
}

// Verify Reset Token
func (s *AuthService) VerifyResetToken(req *dto.VerifyResetTokenRequest) (*dto.VerifyResetTokenResponse, error) {
	// Find password reset record
	var reset models.PasswordReset
	if err := s.db.Where("token = ? AND used = ?", req.Token, false).
		First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &dto.VerifyResetTokenResponse{
				Valid: false,
			}, nil
		}
		return nil, err
	}

	// Check if token expired
	if time.Now().After(reset.ExpiresAt) {
		return &dto.VerifyResetTokenResponse{
			Valid: false,
		}, nil
	}

	return &dto.VerifyResetTokenResponse{
		Valid:  true,
		Email:  reset.Email,
		UserID: reset.UserID.String(),
	}, nil
}

// Reset Password
func (s *AuthService) ResetPassword(req *dto.ResetPasswordRequest) error {
	// Find password reset record
	var reset models.PasswordReset
	if err := s.db.Where("token = ? AND used = ?", req.Token, false).
		Preload("User").
		First(&reset).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired reset token")
		}
		return err
	}

	// Check if token expired
	if time.Now().After(reset.ExpiresAt) {
		return errors.New("reset token has expired")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update user password
	user := reset.User
	user.Password = string(hashedPassword)

	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	// Mark token as used
	reset.Used = true
	if err := s.db.Save(&reset).Error; err != nil {
		return err
	}

	return nil
}

// Resend Verification Email
func (s *AuthService) ResendVerificationEmail(email string) error {
	// Find user by email
	var user models.User
	if err := s.db.Where("email = ? AND deleted_at IS NULL", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}

	// Check if user is already verified
	if user.IsVerified {
		return errors.New("email already verified")
	}

	// Generate new token
	token, err := s.generateToken()
	if err != nil {
		return err
	}

	// Delete old verification tokens
	if err := s.db.Where("user_id = ?", user.ID).Delete(&models.EmailVerification{}).Error; err != nil {
		return err
	}

	// Create new verification record
	verification := &models.EmailVerification{
		ID:        uuid.New(),
		UserID:    user.ID,
		Token:     token,
		Email:     user.Email,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Used:      false,
	}

	if err := s.db.Create(verification).Error; err != nil {
		return err
	}

	// Send email
	if err := s.emailService.SendVerificationEmail(user.Email, user.FirstName, user.LastName, token, s.frontendURL); err != nil {
		return err
	}

	return nil
}

// routes/auth_routes.go
package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"crm-go/controllers/authentication"
	"crm-go/services"
)

func AuthRoutes(router *gin.RouterGroup, db *gorm.DB) {
	emailService := services.NewEmailService()
	authService := services.NewAuthService(db, emailService)
	authController := controllers.NewAuthController(authService)

	authGroup := router.Group("/api/auth")
	{
		// Authentication Routes
		authGroup.POST("/signup", controllers.SignUp)
		authGroup.POST("/logout", controllers.Logout)
		authGroup.POST("/login", controllers.Login)
		authGroup.POST("/login/id", controllers.LoginId)

		// Google OAuth Routes
		authGroup.GET("/google/login", controllers.GoogleLoginHandler)
		authGroup.GET("/google/callback", controllers.GoogleCallbackHandler)
	
		// Email Verification Routes
		authGroup.POST("/send-verification", authController.SendVerificationEmail)
		authGroup.POST("/verify-email", authController.VerifyEmail)
		authGroup.POST("/resend-verification", authController.ResendVerificationEmail)

		// Password Reset Routes
		authGroup.POST("/send-password-reset", authController.SendPasswordReset)
		authGroup.POST("/verify-reset-token", authController.VerifyResetToken)
		authGroup.POST("/reset-password", authController.ResetPassword)
	}
}
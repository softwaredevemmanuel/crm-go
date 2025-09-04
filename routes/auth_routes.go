package routes

import (
	"crm-go/controllers/authentication"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/signup", controllers.SignUp)
		auth.POST("/login", controllers.Login)
		auth.GET("/google/login", controllers.GoogleLoginHandler)
		auth.GET("/google/callback", controllers.GoogleCallbackHandler)
		auth.POST("/forgot-password", controllers.ForgotPassword)
		auth.POST("/reset-password", controllers.ResetPassword)

	}
}

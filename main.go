package main

import (
	"crm-go/config"
	"crm-go/controllers"
	"crm-go/database"
	"crm-go/middleware"

	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// Initialize DB
	config.InitDB()

	database.ConnectDatabase()
	// Init Google OAuth
	config.InitGoogleOauthConfig()
	// r := gin.Default()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.SetTrustedProxies(nil) // trust no proxies in dev

	// Root route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GO CRM 🚀",
		})
	})

	// Signup route
	r.POST("/signup", controllers.SignUp)

	// Login route
	r.POST("/login", controllers.Login)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// Only logged-in users
	protected.GET("/profile", func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")

		c.JSON(200, gin.H{
			"message": "This is your profile",
			"user_id": userID,
			"role":    role,
		})
	})

	// Role-based access
	protected.GET("/admin", middleware.RoleMiddleware("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Admin!"})
	})

	protected.GET("/tutor", middleware.RoleMiddleware("tutor"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Tutor!"})
	})

	protected.GET("/student", middleware.RoleMiddleware("student"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Student!"})
	})

	r.POST("/forgot-password", controllers.ForgotPassword)
	r.POST("/reset-password", controllers.ResetPassword)

	r.GET("/auth/google/login", func(c *gin.Context) {
		role := c.Query("role")

		if role == "" {
			role = "student" // default if not provided
		}

		// Use the role as the state parameter
		url := config.GoogleOauthConfig.AuthCodeURL(role)
		c.Redirect(http.StatusFound, url) // 302
	})

	r.GET("/auth/google/callback", controllers.GoogleCallbackHandler)

	r.Run(":8080")
}

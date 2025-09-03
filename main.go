package main

import (
	"crm-go/config"
	"crm-go/database"
	"crm-go/middleware"
	"crm-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	// Migrate to database	
	database.MigrateDatabase()

	// Initialize DB connection
	config.ConnectDB()

	// Init Google OAuth
	config.InitGoogleOauthConfig()
	
	// Initialize Gin router
	// r := gin.Default()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	r.SetTrustedProxies(nil) // trust no proxies in dev

	
	// Home route
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GO CRM 🚀",
		})
	})
	
	// Register all routes (split into files/packages)
	routes.RegisterAuthRoutes(r)
	routes.RegisterCourseRoutes(r)
	routes.AdminRoutes(r)


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
	


	protected.GET("/tutor", middleware.RoleMiddleware("tutor"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Tutor!"})
	})

	protected.GET("/student", middleware.RoleMiddleware("student"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Student!"})
	})


	r.Run(":8080")
}

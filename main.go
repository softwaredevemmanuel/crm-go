package main

import (
	"crm-go/config"
	"crm-go/database"
	"crm-go/middleware"
	"crm-go/routes"
	"flag"
	"crm-go/database/seeds"

	"github.com/gin-gonic/gin"
	_ "crm-go/docs" // ✅ Import generated docs package

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title GO CRM API
// @version 1.0
// @description This is a comprehensive CRM system for course management with authentication and role-based access control.

// @contact.name API Support
// @contact.email support@gocrm.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Migrate to database	
	database.MigrateDatabase()

	// Initialize DB connection
	config.ConnectDB()

	// Init Google OAuth
	config.InitGoogleOauthConfig()
	
	// Initialize Gin router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.SetTrustedProxies(nil) // trust no proxies in dev
	r.Use(middleware.SessionMiddleware())

	// ✅ Swagger documentation route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	
	// Home route
	// @Summary Welcome endpoint
	// @Description Returns a welcome message
	// @Tags General
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router / [get]
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GO CRM 🚀",
		})
	})
	
	// Register all routes
	routes.RegisterAuthRoutes(r)
	routes.RegisterCourseRoutes(r)
	routes.AdminRoutes(r)
	routes.CategoryRoutes(r)
	routes.CourseProductRoutes(r)
	routes.CourseCategoryRoutes(r)
	routes.ProductRoutes(r)

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())

	// Only logged-in users
	// @Summary Get user profile
	// @Description Get authenticated user's profile information
	// @Tags User
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Router /api/profile [get]
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
	// @Summary Tutor endpoint
	// @Description Endpoint accessible only to tutors
	// @Tags User
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 403 {object} map[string]interface{}
	// @Router /api/tutor [get]
	protected.GET("/tutor", middleware.RoleMiddleware("tutor"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Tutor!"})
	})

	// @Summary Student endpoint
	// @Description Endpoint accessible only to students
	// @Tags User
	// @Security BearerAuth
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Failure 401 {object} map[string]interface{}
	// @Failure 403 {object} map[string]interface{}
	// @Router /api/student [get]
	protected.GET("/student", middleware.RoleMiddleware("student"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome Student!"})
	})

	// Add a flag to run seeder
	seed := flag.Bool("seed", false, "Run database seeder")
	flag.Parse()

	if *seed {
		seeds.SeedCourses()
		seeds.SeedUsers()
		seeds.SeedCategories()
		seeds.SeedCourseCategories()
		seeds.SeedCourseProductsTable()
		seeds.SeedProducts()
		return
	}

	r.Run(":8080")
}
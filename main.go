package main

import (
	"crm-go/config"
	"crm-go/database"
	"crm-go/database/seeds"
	"crm-go/middleware"
	"crm-go/routes"
	"flag"
	"github.com/gin-contrib/cors"
	"time"
	"os"
	"log"

	_ "crm-go/docs"

	"github.com/gin-gonic/gin"

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
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080", "https://hr-app-ecru-nine.vercel.app"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "email"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.Use(middleware.SessionMiddleware())



	// ✅ Swagger documentation route
	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Home route
	// @Summary Welcome endpoint
	// @Description Returns a welcome message
	// @Tags General
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router / [get]
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GO CRM 🚀. For Swagger Documentation, visit: http://localhost:8080/swagger/index.html#",
		})
	})

	// Register all routes
	routes.RegisterAuthRoutes(r)
	routes.CourseRoutes(r)
	routes.AdminRoutes(r)
	routes.CategoryRoutes(r)
	routes.CourseProductRoutes(r)
	routes.CourseCategoryRoutes(r)
	routes.ProductRoutes(r)
	routes.EnrollmentRoutes(r)
	routes.AnnouncementRoutes(r)
	routes.AssignmentRoutes(r)
	routes.AssignmentSubmissionRoutes(r, config.DB)
	routes.LessonRoutes(r, config.DB)
	routes.GradeRoutes(r, config.DB)
	routes.ModuleRoutes(r)
	routes.TopicRoutes(r)
	routes.CourseMaterialRoutes(r)
	routes.LiveClassRoutes(r, config.DB)
	routes.ObjectiveQuestionRoutes(r, config.DB)

	// Example curl command to clear DB (replace with your server address):
	// curl -X DELETE "http://localhost:8080/admin/clear-db" \
	//      -H "Content-Type: application/json" \
	//      -d '{"password":"mypassword"}'
	routes.AdminDangerRoutes(r)

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
		seeds.SeedAnnouncements()
		seeds.SeedModules()
		seeds.SeedLessons()
		seeds.SeedTopics()
		seeds.SeedCourseMaterials()
		seeds.SeedObjectiveQuestions()
		return
	}

	SetupSwagger(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // local default
	}

	log.Printf("Server starting on port %s", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatal(err)
	}
}

func SetupSwagger(r *gin.Engine) {
	url := "http://localhost:8080/swagger/doc.json" // your swagger.json
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL(url),
		ginSwagger.DefaultModelsExpandDepth(-1),
		ginSwagger.PersistAuthorization(true), // 👈 keeps token after refresh
	))
}

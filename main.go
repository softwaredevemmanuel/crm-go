package main

import (
	"crm-go/controllers"
	"crm-go/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()

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

	r.Run(":8080")
}

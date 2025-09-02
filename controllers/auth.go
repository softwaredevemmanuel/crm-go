package controllers

import (
	"crm-go/database"
	"crm-go/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type SignUpInput struct {
	FirstName     string       `json:"first_name" binding:"required"`
	LastName     string       `json:"last_name" binding:"required"`
	Email    string       `json:"email" binding:"required,email"`
	Password string       `json:"password" binding:"required,min=6"`
	Role     models.Role  `json:"role" binding:"required"`
}

func SignUp(c *gin.Context) {
	var input SignUpInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	user := models.User{
		FirstName:     input.FirstName,
		LastName:     input.LastName,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
	}

	result := database.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": user})
}

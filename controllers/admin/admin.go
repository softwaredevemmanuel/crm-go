package admin

import (
	"crm-go/database"
	"github.com/gin-gonic/gin"
	"net/http"
)
// Example curl command to clear DB (replace with your server address):
// curl -X DELETE "http://localhost:8080/admin/clear-db" \
//      -H "Content-Type: application/json" \
//      -d '{"password":"mypassword"}'

const adminPassword = "mypassword" // ⚠️ better to keep this in env variable

type ClearRequest struct {
	Password string `json:"password"`
}

// ClearDatabaseHandler - Danger zone: Wipes all data!
func ClearDatabaseHandler(c *gin.Context) {
	var req ClearRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password is required in request body"})
		return
	}

	// Password check
	if req.Password != adminPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
		return
	}

	// Run clear
	if err := database.ClearDatabase(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to clear database", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "✅ All data cleared successfully"})
}

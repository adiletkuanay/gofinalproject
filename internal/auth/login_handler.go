package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// LoginHandler processes login requests and returns a JWT if successful
func LoginHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	sqlDB := db.(*sql.DB)

	var user User
	err := sqlDB.QueryRow("SELECT id, password, role FROM users WHERE username = $1", credentials.Username).Scan(&user.ID, &user.Password, &user.Role)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	} else if err != nil {
		log.Printf("Error querying user by username: %v", err) // Log error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query error"})
		return
	}

	if !CheckPassword(user.Password, credentials.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		return
	}

	token, err := GenerateToken(user.ID, user.Role) // Generate JWT
	if err != nil {
		log.Printf("Error generating JWT: %v", err) // Log error message
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token}) // Return the JWT
}

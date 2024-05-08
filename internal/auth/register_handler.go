package auth

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRegistrationHandler(c *gin.Context) {
	db, exists := c.Get("db")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not found"})
		return
	}

	sqlDB, ok := db.(*sql.DB)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid database connection"})
		return
	}

	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	if newUser.Username == "" || newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	hashedPassword, err := HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = sqlDB.Exec("INSERT INTO users (username, password, role) VALUES ($1, $2, $3)", newUser.Username, hashedPassword, "user")
	if err != nil {
		log.Printf("Error creating user: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create account"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created"})
}

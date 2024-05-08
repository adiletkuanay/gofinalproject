package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func AdminOnlyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "No user found"})
			return
		}

		userRole := user.(User).Role

		if userRole != "admin" {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Access denied. Admins only."})
			return
		}

		c.Next()
	}
}

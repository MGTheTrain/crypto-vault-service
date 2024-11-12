package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware is a simple authentication middleware
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Example authentication logic
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token required"})
			c.Abort()
			return
		}

		// Optionally, validate the token here

		c.Next()
	}
}

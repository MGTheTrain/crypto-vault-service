package v1

import "github.com/gin-gonic/gin"

// SendSuccess is a utility to send successful responses
func SendSuccess(c *gin.Context, data interface{}) {
	c.JSON(200, gin.H{"data": data})
}

// SendError is a utility to send error responses
func SendError(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, gin.H{"error": message})
}

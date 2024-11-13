package v1

import "github.com/gin-gonic/gin"

// SetupRoutes sets up all the API routes for version 1.
func SetupRoutes(r *gin.Engine) {
	v1 := r.Group("/api/v1") // Prefix for v1 routes

	// Define v1 API routes
	v1.GET("/users/:id", GetUser)
	v1.POST("/users", CreateUser)
}

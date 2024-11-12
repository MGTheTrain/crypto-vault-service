package main

import (
	v1 "crypto_vault_service/internal/api/v1"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Set up version 1 routes
	v1.SetupRoutes(r)

	// Optional: Apply a global middleware
	r.Use(v1.AuthMiddleware())

	// Start the server
	r.Run(":8080") // By default it will listen on :8080
}

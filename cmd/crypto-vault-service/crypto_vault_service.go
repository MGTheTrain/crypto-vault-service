package main

import (
	v1 "crypto_vault_service/internal/api/v1"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Create a new Gin router
	r := gin.Default()

	// Set up version 1 routes
	v1.SetupRoutes(r)

	// Optional: Apply a global middleware
	r.Use(v1.AuthMiddleware())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

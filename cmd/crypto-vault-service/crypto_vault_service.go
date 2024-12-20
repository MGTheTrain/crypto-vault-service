package main

import (
	v1 "crypto_vault_service/internal/api/v1"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	v1.SetupRoutes(r)

	// r.Use(v1.AuthMiddleware())

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}

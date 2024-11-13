package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUser handles the GET request to fetch a user by ID
func GetUser(c *gin.Context) {
	id := c.Param("id")
	// Here you would typically fetch user data from the database based on `id`
	user := User{
		ID:    1,
		Name:  "John Doe",
		Email: "john.doe@example.com",
	}

	fmt.Printf("ID is: %s\n", id)

	c.JSON(http.StatusOK, user)
}

// CreateUser handles POST requests to create a new user
func CreateUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid data"})
		return
	}

	// Logic to save the user to the database

	c.JSON(http.StatusCreated, user)
}

package middleware

import (
	"go_api_product_review/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// AuthMiddleware validates a Bearer token for incoming requests.
func AuthMiddleware(c *gin.Context) {
	// Load environment variables from .env file to access SECRET_KEY
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file") // Log error and terminate if loading fails
	}

	// Retrieve the token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	// Skip authentication for the Swagger documentation endpoint
	if c.FullPath() == "/swagger/*any" {
		c.Next() // Proceed with the next handler
		return
	}

	// If no token is provided, return an unauthorized error
	if tokenString == "" {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "No token provided",
			Details: err.Error(),
		})
		c.Abort()
		return
	}

	// Remove the "Bearer " prefix from the token if it's present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:] // Strip "Bearer " to get the actual token
	} else {
		// Return an error if the token format is invalid
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Invalid token format",
			Details: err.Error(),
		})
		c.Abort()
		return
	}

	// If the token is empty after stripping, return an error
	if len(tokenString) == 0 {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Token is empty",
			Details: err.Error(),
		})
		c.Abort()
		return
	}

	// Retrieve the secret token from the environment variables
	secretToken := os.Getenv("SECRET_KEY")
	if secretToken == "" {
		// Return an error if the secret token is missing
		c.JSON(http.StatusInternalServerError, models.ErrorResponse{
			Message: "SECRET_KEY environment variable is missing",
			Details: err.Error(),
		})
		c.Abort()
		return
	}

	// If the provided token doesn't match the secret, return an invalid token error
	if tokenString != secretToken {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse{
			Message: "Invalid token",
			Details: err.Error(),
		})
		c.Abort()
		return
	}

	// If the token is valid, proceed to the next handler
	c.Next()
}

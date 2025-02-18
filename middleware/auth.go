package middleware

import (
	"go_api_product_review/models"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// AuthMiddleware is a middleware that validates a simple Bearer token
func AuthMiddleware(c *gin.Context) {
	// Load environment variables from .env file
	// This is used to access the SECRET_KEY for token validation
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file") // If loading fails, log an error and terminate
	}

	// Retrieve the token from the Authorization header
	tokenString := c.GetHeader("Authorization")

	// Skip authentication for the Swagger documentation endpoint
	if c.FullPath() == "/swagger/*any" {
		c.Next() // Proceed with the next middleware/handler
		return
	}

	// Check if the Authorization header is missing or empty
	if tokenString == "" {
		// Respond with an error if no token is found
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "No token provided",
				Details: err.Error(),
			})
		}
		c.Abort() // Stop further request processing
		return
	}

	// Remove the "Bearer " prefix if it's present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:] // Strip "Bearer " to get the actual token
	} else {
		// If the token format is incorrect, respond with an error
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Invalid token format",
				Details: err.Error(),
			})
		}
		c.Abort() // Stop further request processing
		return
	}

	// If the token is empty after removing the prefix, respond with an error
	if len(tokenString) == 0 {
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Token is empty",
				Details: err.Error(),
			})
		}
		c.Abort() // Stop further request processing
		return
	}

	// Retrieve the secret token from the environment variables
	secretToken := os.Getenv("SECRET_KEY")
	if secretToken == "" {
		// If the secret token is missing, respond with a server error
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Message: "SECRET_KEY environment variable is missing",
				Details: err.Error(),
			})
		}
		c.Abort() // Stop further request processing
		return
	}

	// Compare the provided token with the secret token from the environment
	if tokenString != secretToken {
		// If the token does not match, respond with an invalid token error
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.ErrorResponse{
				Message: "Invalid token",
				Details: err.Error(),
			})
		}
		c.Abort() // Stop further request processing
		return
	}

	// If the token is valid, proceed to the next handler
	c.Next()
}

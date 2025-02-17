package main

import (
	"go_api_product_review/api"
	"go_api_product_review/cache"
	"go_api_product_review/db"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize database
	db.InitPostgres()

	// Initialize Redis cache
	cache.InitRedis()

	// Create Gin router
	router := gin.Default()

	// Set up routes
	api.RegisterProductRoutes(router)
	api.RegisterReviewRoutes(router)

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

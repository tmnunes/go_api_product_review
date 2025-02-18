package main

import (
	"go_api_product_review/api"
	"go_api_product_review/cache"
	"go_api_product_review/db"
	"go_api_product_review/middleware"
	"log"
	"os"

	_ "go_api_product_review/cmd/docs" // Import the docs from the cmd/docs folder

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Set up middleware (authentication)
	router.Use(middleware.AuthMiddleware)

	// Set up routes
	api.RegisterProductRoutes(router)
	api.RegisterReviewRoutes(router)

	// Swagger router
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

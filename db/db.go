package db

import (
	"go_api_product_review/models"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import the PostgreSQL dialect for GORM
)

// DB is the global database connection instance.
var DB *gorm.DB

// InitPostgres initializes the connection to the PostgreSQL database.
func InitPostgres() {
	var err error
	// Get the database connection string from environment variables.
	dsn := os.Getenv("POSTGRES_DSN")
	// Open a connection to the PostgreSQL database.
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		// Log and exit if the connection fails.
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Automatically migrate the Product and Review models.
	DB.AutoMigrate(&models.Product{}, &models.Review{})
}

// GetDB returns the current database instance.
func GetDB() *gorm.DB {
	// If the database isn't initialized, log an error and stop.
	if DB == nil {
		log.Fatal("Database not initialized")
	}
	return DB
}

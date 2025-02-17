package db

import (
	"go_api_product_review/models"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // Import the PostgreSQL dialect for GORM
)

// DB represents the database connection
var DB *gorm.DB

// InitPostgres initializes the connection to the PostgreSQL database
func InitPostgres() {
	var err error
	dsn := os.Getenv("POSTGRES_DSN")
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB.AutoMigrate(&models.Product{}, &models.Review{})
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}

package db

import (
	"go_api_product_review/models"
	"log"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitPostgres() {
	var err error
	dsn := os.Getenv("POSTGRES_DSN")
	DB, err = gorm.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	DB.AutoMigrate(&models.Product{}, &models.Review{})
}

func GetDB() *gorm.DB {
	return DB
}

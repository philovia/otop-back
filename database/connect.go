package database

import (
	// "fmt"
	"log"
	// "os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"

	"github.com/m/models"
)

var DB *gorm.DB

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

// SetupDatabase initializes the database connection using GORM
func SetupDatabase() {
	LoadEnv() // Load environment variables from .env file

	// Data source name (DSN) for connecting to PostgreSQL
	dsn := "host=localhost port=5432 user=postgres dbname=postgres password=postgres"

	// Connect to the PostgreSQL database using GORM
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL using GORM")

	// Auto-migrate the models
	DB.AutoMigrate(&models.Product{}, &models.User{}, &models.Supplier{})
}

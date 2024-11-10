package database

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/joho/godotenv"

	"github.com/m/models"
)

var DB *gorm.DB

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func SetupDatabase() {
	LoadEnv()

<<<<<<< HEAD
	dsn := "host=localhost port=5432 user=postgres dbname=postgres password=postgres"

=======
	dsn := "host=localhost port=5432 user=postgres dbname=Reggie_Macazar password=postgres"
>>>>>>> 36cf5b4b0c38771a532201f6a055694672691442
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL using GORM")

<<<<<<< HEAD
	DB.AutoMigrate(&models.Product{}, &models.User{}, &models.Supplier{}, &models.Order{})
=======
	DB.AutoMigrate(&models.Product{}, &models.User{}, &models.Supplier{})
>>>>>>> 36cf5b4b0c38771a532201f6a055694672691442
}

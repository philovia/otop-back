package main

import (
	// "fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/m/database"
	"github.com/m/models"
	"github.com/m/routes"
)

func main() {

	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	database.SetupDatabase()

	if err := database.DB.AutoMigrate(&models.Supplier{}); err != nil {
		log.Fatalf("Could not auto-migrate supplier table: %v", err)
	}

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	routes.UserRoutes(app)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Fatal(app.Listen(":" + port))
}

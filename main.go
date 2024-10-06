package main

import (
	// "fmt"
	"log"
	"os"

	// "gorm.io/driver/postgres"
	// "gorm.io/gorm"
	// "gorm.io/gorm/logger"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"

	// "github.com/gofiber/fiber/v2/middleware/cors"
	// "github.com/joho/godotenv"

	// "github.com/mreym/backend-otop-ph/controllers"
	"github.com/m/database"
	// middleware "github.com/mreym/backend-otop-ph/midleware"
	"github.com/m/routes"
)

func main() {

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll("uploads", os.ModePerm); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	// Setup the database connection
	database.SetupDatabase() // Ensure this is called before any routes

	// Create a new Fiber app
	app := fiber.New()

	// Enable CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Allow all origins (you may restrict this in production)
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// Setup the routes
	routes.UserRoutes(app)

	// Get the port from environment variables, default to 8081 if not set
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	// Start the server
	log.Fatal(app.Listen(":" + port))
}

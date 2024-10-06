package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/m/controllers"
	middleware "github.com/m/middleware"
)

func UserRoutes(app *fiber.App) {

	//  Public routes
	app.Post("/api/register", controllers.Register)
	app.Post("/api/login", controllers.Login)

	// Protected routes
	api := app.Group("/api", middleware.Authentication()) // Use JWT middleware

	// Example protected routes (e.g., add/update products)
	api.Post("/products", middleware.IsSupplier(), controllers.AddProduct)
	// api.Get("/products", controllers.GetProducts)
	// api.Put("/products/:id", controllers.UpdateProduct)

	// Logout route
	api.Post("/logout", controllers.Logout) // Add this line
}

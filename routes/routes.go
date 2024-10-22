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
	app.Post("/logout", controllers.Logout) // Add this line

	// Protected routes
	admin := app.Group("/admin", middleware.Authentication()) // Use JWT middleware
	admin.Post("/supplier", controllers.CreateSupplier)
	admin.Get("/suppliers", controllers.GetAllSuppliers)
	admin.Get("/suppliers/storeName", controllers.GetSupplierByStoreName)
	admin.Put("/suppliers/:id", controllers.UpdateSupplier)
	admin.Delete("/supplier/:id", controllers.DeleteSupplier)

}

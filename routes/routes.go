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

	app.Post("/logout", middleware.Authentication(), controllers.Logout)
	// Protected routes
	admin := app.Group("/admin", middleware.Authentication(), middleware.IsAdmin())
	admin.Post("/supplier", controllers.CreateSupplier)
	admin.Post("/login", controllers.SupplierLogin)
	admin.Get("/suppliers", controllers.GetAllSuppliers)
	admin.Get("/suppliers/email/:email", controllers.GetSupplierByEmail)
	admin.Get("/suppliers/store_name/:storeName", controllers.GetSupplierByStoreName)
	admin.Put("/suppliers/:id", controllers.UpdateSupplier)
	admin.Delete("/supplier/:id", controllers.DeleteSupplier)

	// routes for the supplier
	supplier := app.Group("/supplier", middleware.Authentication())
	supplier.Post("/add/product", controllers.AddProduct)
	supplier.Put("/update/product/:id", controllers.UpdateProduct)
	supplier.Delete("/delete/product/:id", controllers.DeleteProduct)
	supplier.Get("/products/name/:name", controllers.GetProductByName)

	// routes fpor the cashier
	cashier := app.Group("/cashier", middleware.IsCashier())
	cashier.Get("/products/:name", controllers.GetProductByName)

}

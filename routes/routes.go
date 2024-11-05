package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/m/controllers"
	middleware "github.com/m/middleware"
)

func UserRoutes(app *fiber.App) {

	app.Post("/register", controllers.Register)
	app.Post("/login", controllers.UnifiedLogin)
	app.Post("/logout", controllers.Logout)

	// Admin-only routes
	supplier := app.Group("/supplier")
	supplier.Use(middleware.JWTProtected)
	supplier.Use(middleware.IsAdmin)
	supplier.Post("/", controllers.CreateSupplier)
	supplier.Get("/", controllers.GetSuppliers)
	supplier.Get("/:storeName", controllers.GetSupplierByStoreName)
	supplier.Put("/:storeName", controllers.UpdateSupplier)
	supplier.Delete("/:storeName", controllers.DeleteSupplier)

	// Product management routes for suppliers
	supplierRoutes := app.Group("/products", middleware.IsSupplier, middleware.JWTProtected)
	supplierRoutes.Post("/", controllers.AddProduct)
	supplierRoutes.Get("/", controllers.GetProducts)
	supplierRoutes.Get("/:name", controllers.GetProductByName)
	supplierRoutes.Put("/:name", controllers.UpdateProduct)
	supplierRoutes.Delete("/:name", controllers.DeleteProduct)

	//for admin and cashier
	app.Get("/api/products/total_quantity", controllers.GetTotalQuantity)
	app.Get("/api/products/total_price", controllers.GetTotalPrice)
	app.Get("/products/:name", middleware.JWTProtected, controllers.GetProductByName)

}

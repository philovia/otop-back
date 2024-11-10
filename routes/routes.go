package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/m/controllers"
	middleware "github.com/m/middleware"
)

func UserRoutes(app *fiber.App) {

	// public routes (DONE)
	api := app.Group("/api")
	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.UnifiedLogin)
	api.Post("/logout", controllers.Logout)

	// Admin-only routes(DONE)
	supplier := app.Group("/supplier")
	supplier.Use(middleware.JWTProtected)
	supplier.Use(middleware.IsAdmin)
	supplier.Post("/", controllers.CreateSupplier)
	supplier.Get("/", controllers.GetSuppliers)
	supplier.Get("/:storeName", controllers.GetSupplierByStoreName)
	supplier.Put("/:storeName", controllers.UpdateSupplier)
	supplier.Delete("/:storeName", controllers.DeleteSupplier)

	// Product management routes for suppliers (DONE)
	supplierRoutes := app.Group("/products", middleware.IsSupplier, middleware.JWTProtected)
	supplierRoutes.Post("/", controllers.AddProduct)
	supplierRoutes.Get("/", controllers.GetProducts)
	supplierRoutes.Get("/:name", controllers.GetProductByName)
	supplierRoutes.Put("/:name", controllers.UpdateProduct)
	supplierRoutes.Delete("/:name", controllers.DeleteProduct)

	// the supplier will confirmed the order from admin
	supplierRoutes.Put("/orders/confirm/:id", controllers.ConfirmOrder)

	// Order Management for the admin with supplier
	admin := app.Group("/order", middleware.IsAdmin, middleware.JWTProtected)
	admin.Post("/", controllers.CreateOrder)
	admin.Get("/", controllers.GetOrders)
	admin.Get("/:id", controllers.GetOrder)
	admin.Put("/:id", controllers.UpdateOrder)
	admin.Delete("/:id", controllers.DeleteOrder)
 
	//for admin and cashier
	app.Get("/api/products/total_quantity", controllers.GetTotalQuantity)
	// app.Get("/api/products/total_price", controllers.GetTotalPrice)
	app.Get("/products/:name", middleware.JWTProtected, controllers.GetProductByName)

}

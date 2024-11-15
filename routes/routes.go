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

	// Add products by supplier (DONE)
	supplierRoutes := app.Group("/products", middleware.IsSupplier, middleware.JWTProtected)
	supplierRoutes.Post("/", controllers.AddProduct)
	supplierRoutes.Get("/", controllers.GetProducts)
	supplierRoutes.Get("/:id", controllers.GetProductByName)
	supplierRoutes.Put("/:id", controllers.UpdateProduct)
	supplierRoutes.Delete("/:id", controllers.DeleteProduct)

	// the supplier will confirmed the order from admin(NOT YET)
	app.Patch("/orders/confirm/:id", controllers.ConfirmOrder, middleware.IsSupplier)

	// Otop products (DONE)
	otop := app.Group("/otop", middleware.IsAdmin, middleware.JWTProtected)
	otop.Post("/", controllers.CreateOtopProduct)
	otop.Get("/", controllers.GetOtopProducts)
	otop.Get("/:id", controllers.GetOtopProductByID)
	otop.Put("/:id", controllers.UpdateOtopProduct)
	otop.Delete("/:id", controllers.DeleteOtopProduct)

	//Can Get Total Otop Products Stocks & Name(DONE)
	app.Get("/api/otop/total_quantity", controllers.GetOtopTotalQuantity)              // total of otop products quantity of all products
	app.Get("/api/otop/total_quantity_name", controllers.GetOtopTotalQuantityName)     // diffrerent store name and total  quantity of products
	app.Get("/api/otop/total_products", controllers.GetOtopTotalProducts)              // total number of products(USED)
	app.Get("/api/otop/total_categories", controllers.GetTotalProductsByCategory)      // total products on food and non-food(USED)
	app.Get("/api/otop/total_suppliers", controllers.GetTotalSuppliers)                // count all suppliers(USED)
	app.Get("/api/otop/total_suppliers_product", controllers.GetSupplierProductCounts) // supplier and number of products
	// Order Management for the admin with supplier (DONE)
	admin := app.Group("/order")
	admin.Post("/", controllers.CreateOrder)
	admin.Get("/", controllers.GetOrders)
	admin.Get("/:id", controllers.GetOrder)
	admin.Put("/:id", controllers.UpdateOrder)
	admin.Delete("/:id", controllers.DeleteOrder)

	//for admin and cashier (DONE)
	app.Get("/api/products/total_quantity", controllers.GetTotalQuantity)
	app.Get("/api/products", controllers.GetProducts)
	app.Get("/api/products/supplier/:supplier_id", controllers.GetProductsByStore)

}

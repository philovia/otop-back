package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/m/database"
	"github.com/m/models"
)

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order data"})
	}

	var product models.Product
	if err := database.DB.First(&product, order.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}
	order.Status = "pending"
	order.OrderDate = time.Now()

	// Check if there is enough stock available
	if int64(order.Quantity) > product.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient product stock"})
	}

	// Set the product name and price from the product record
	order.ProductName = product.Name
	order.Price = product.Price

	// Update the product stock in the database after the order is created
	product.Quantity -= int64(order.Quantity)
	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product stock"})
	}

	// Create the order in the database
	if err := database.DB.Create(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create order"})
	}

	return c.Status(fiber.StatusCreated).JSON(order)
}

func GetOrders(c *fiber.Ctx) error {
	var orders []models.Order

	database.DB.Find(&orders)

	return c.JSON(orders)
}

func GetOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	database.DB.First(&order, id)

	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	return c.JSON(order)
}

func UpdateOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	database.DB.First(&order, id)

	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	if err := c.BodyParser(&order); err != nil {
		return err
	}

	database.DB.Save(&order)

	return c.JSON(order)
}

func DeleteOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	database.DB.First(&order, id)

	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	database.DB.Delete(&order)

	return c.JSON(fiber.Map{
		"message": "Order deleted",
	})
}

func ConfirmOrder(c *fiber.Ctx) error {
	id := c.Params("id")
	var order models.Order

	// Retrieve the order from the database by its ID
	if err := database.DB.First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}

	// Check if the order is still "pending"
	if order.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order already confirmed or completed"})
	}

	// Get the JWT token from the request context (for the supplier)
	userToken := c.Locals("user")
	if userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// Parse the JWT claims
	claims, ok := userToken.(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Get the supplier ID from the claims
	supplierID, ok := claims["supplier_id"].(float64) // JWT stores numbers as float64
	if !ok {
		fmt.Println("Supplier ID not found in claims")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Supplier ID not found in token"})
	}
	fmt.Println("Supplier ID from Token: ", supplierID)
	fmt.Println("Supplier ID from Order: ", order.SupplierID)

	// Check if the supplier ID matches the order's supplier ID
	if uint(supplierID) != order.ProductID {
		fmt.Println("Supplier ID mismatch!")
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to confirm this order"})
	}

	// Update the order status to "confirmed_by_supplier"
	order.Status = "confirmed_by_supplier"

	// Save the updated order
	if err := database.DB.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to confirm order"})
	}

	// Return the updated order
	return c.JSON(order)
}

package controllers

import (
	// "fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	// "github.com/golang-jwt/jwt/v4"
	"github.com/m/database"
	"github.com/m/models"
)

func CreateOrder(c *fiber.Ctx) error {
	var order models.Order

	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid order data"})
	}
	order.Status = "pending"
	order.OrderDate = time.Now()

	var product models.Product
	if err := database.DB.First(&product, order.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	if product.Quantity < order.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient product stock"})
	}

	order.ProductName = product.Name
	order.Price = product.Price

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

	// Find the existing order
	database.DB.First(&order, id)

	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Order not found",
		})
	}

	// Parse the new order data
	if err := c.BodyParser(&order); err != nil {
		return err
	}

	// Fetch the product details from the database
	var product models.Product
	if err := database.DB.First(&product, order.ProductID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	// Check if the order quantity exceeds available product stock
	if order.Quantity > product.Quantity {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Insufficient product stock",
		})
	}

	// Update the product stock by deducting the ordered quantity
	product.Quantity -= order.Quantity
	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product stock"})
	}

	// Save the updated order
	database.DB.Save(&order)

	// Return the updated order
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

	if err := database.DB.First(&order, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Order not found"})
	}
	if order.Status != "pending" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Order already confirmed or completed"})
	}
	supplierID, ok := c.Locals("supplier_id").(float64)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	if uint(supplierID) != order.SupplierID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to confirm this order"})
	}

	// Update the order status to confirmed
	order.Status = "confirmed_by_supplier"
	order.UpdatedAt = time.Now()

	if err := database.DB.Save(&order).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to confirm order"})
	}

	return c.JSON(order)
}

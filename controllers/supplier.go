package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m/database"
	"github.com/m/models"
)

// CreateSupplier creates a new supplier (Admin only)
func CreateSupplier(c *fiber.Ctx) error {
	var supplier models.Supplier

	// Parse request body
	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Save the supplier in the database
	if err := database.DB.Create(&supplier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot create supplier",
		})
	}

	// Return the created supplier
	return c.Status(201).JSON(supplier)
}

// UpdateSupplier updates an existing supplier (Admin only)
func UpdateSupplier(c *fiber.Ctx) error {
	id := c.Params("id")

	var supplier models.Supplier
	if err := database.DB.First(&supplier, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	// Parse the updated data
	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Update the supplier in the database
	if err := database.DB.Save(&supplier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update supplier",
		})
	}

	return c.JSON(supplier)
}

// DeleteSupplier deletes a supplier (Admin only)
func DeleteSupplier(c *fiber.Ctx) error {
	id := c.Params("id")

	var supplier models.Supplier
	if err := database.DB.First(&supplier, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	// Delete the supplier
	if err := database.DB.Delete(&supplier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete supplier",
		})
	}

	return c.JSON(fiber.Map{"message": "Supplier deleted successfully"})
}

// GetAllSuppliers fetches all suppliers (Admin only)
func GetAllSuppliers(c *fiber.Ctx) error {
	var suppliers []models.Supplier

	// Fetch all suppliers
	if err := database.DB.Find(&suppliers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch suppliers",
		})
	}

	return c.JSON(suppliers)
}

func GetSupplierByStoreName(c *fiber.Ctx) error {
	storeName := c.Query("storeName")

	var supplier models.Supplier
	if err := database.DB.Where("store_name = ?", storeName).First(&supplier).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Supplier not found",
		})
	}

	return c.JSON(supplier)
}

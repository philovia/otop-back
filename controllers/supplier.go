package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m/database"
	"github.com/m/models"
)

func CreateSupplier(c *fiber.Ctx) error {
	var supplier models.Supplier
	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	database.DB.Create(&supplier)
	return c.JSON(supplier)
}

func GetSuppliers(c *fiber.Ctx) error {
	var suppliers []models.Supplier
	database.DB.Find(&suppliers)
	return c.JSON(suppliers)
}

func GetSupplierByStoreName(c *fiber.Ctx) error {
	storeName := c.Params("storeName")
	var supplier models.Supplier
	database.DB.Where("store_name = ?", storeName).First(&supplier)
	if supplier.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Supplier not found"})
	}
	return c.JSON(supplier)
}

func UpdateSupplier(c *fiber.Ctx) error {
	storeName := c.Params("storeName")
	var supplier models.Supplier

	if err := database.DB.First(&supplier, storeName).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Supplier not found"})
	}

	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	database.DB.Save(&supplier)
	return c.JSON(supplier)
}

func DeleteSupplier(c *fiber.Ctx) error {
	storeName := c.Params("storeName")
	var supplier models.Supplier
	if err := database.DB.First(&supplier, storeName).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Supplier not found"})
	}

	database.DB.Delete(&supplier)
	return c.JSON(fiber.Map{"message": "Supplier deleted successfully"})
}

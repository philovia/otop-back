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

	if err := database.DB.Create(&supplier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create supplier"})
	}
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
	// <<<<<<< HEAD/
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
	// <<<<<<< HEAD
	storeName := c.Params("storeName")
	var supplier models.Supplier
	if err := database.DB.First(&supplier, storeName).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Supplier not found"})
	}

	database.DB.Delete(&supplier)
	return c.JSON(fiber.Map{"message": "Supplier deleted successfully"})
}

func CountSuppliersByStoreName(c *fiber.Ctx) error {
	type StoreNameCount struct {
		StoreName string `json:"store_name"`
		Count     int64  `json:"count"`
	}

	var results []StoreNameCount
	err := database.DB.Model(&models.Supplier{}).
		Select("store_name, COUNT(*) as count").
		Group("store_name").
		Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count suppliers"})
	}

	return c.JSON(results)
}

func GetTotalSuppliers(c *fiber.Ctx) error {
	var count int64
	if err := database.DB.Model(&models.Supplier{}).Count(&count).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count suppliers"})
	}

	return c.JSON(fiber.Map{"total_suppliers": count})
}

func GetSupplierProductCounts(c *fiber.Ctx) error {
	type SupplierProductCount struct {
		StoreName    string `json:"store_name"`
		ProductCount int64  `json:"product_count"`
	}

	var results []SupplierProductCount
	err := database.DB.Model(&models.Product{}).
		Select("suppliers.store_name, COUNT(products.id) as product_count").
		Joins("JOIN suppliers ON suppliers.id = products.supplier_id").
		Group("suppliers.store_name").
		Scan(&results).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve supplier product counts"})
	}

	return c.JSON(results)
}

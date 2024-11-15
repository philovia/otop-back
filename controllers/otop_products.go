package controllers

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"

	// "github.com/golang-jwt/jwt/v4"
	"github.com/m/database"
	"github.com/m/models"
	// "gorm.io/gorm"
)

func CreateOtopProduct(c *fiber.Ctx) error {
	var otopProduct models.OtopProducts

	// Parse the request body into the product struct
	if err := c.BodyParser(&otopProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product data"})
	}

	// Set the creation time
	otopProduct.CreatedAt = time.Now()

	// Save the product in the database
	if err := database.DB.Create(&otopProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	// Return the created product as JSON
	return c.Status(fiber.StatusCreated).JSON(otopProduct)
}

func GetOtopProducts(c *fiber.Ctx) error {
	var otopProduct []models.OtopProducts

	// Fetch all products from the database
	database.DB.Find(&otopProduct)

	// Return the products as JSON
	return c.JSON(otopProduct)
}

func GetProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var otopProduct models.OtopProducts

	// Fetch the product by ID
	database.DB.First(&otopProduct, id)

	if otopProduct.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	// Return the product as JSON
	return c.JSON(otopProduct)
}

func UpdateOtopProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var otopProduct models.OtopProducts

	// Find the existing product
	database.DB.First(&otopProduct, id)

	if otopProduct.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	// Parse the new product data
	if err := c.BodyParser(&otopProduct); err != nil {
		return err
	}

	// Save the updated product
	if err := database.DB.Save(&otopProduct).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	// Return the updated product
	return c.JSON(otopProduct)
}

func DeleteOtopProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var otopProduct models.OtopProducts

	// Find the product by ID
	database.DB.First(&otopProduct, id)

	if otopProduct.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	// Delete the product
	database.DB.Delete(&otopProduct)

	// Return a success message
	return c.JSON(fiber.Map{
		"message": "Product deleted",
	})
}

func GetOtopTotalQuantity(c *fiber.Ctx) error {
	log.Println("Received request to calculate total quantity of products")

	var totalQuantity int64
	err := database.DB.Model(&models.OtopProducts{}).Select("SUM(quantity)").Scan(&totalQuantity).Error

	if err != nil {
		log.Println("Error calculating total quantity:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate total quantity"})
	}

	log.Println("Total quantity calculated:", totalQuantity)
	return c.JSON(fiber.Map{"total_quantity": totalQuantity})
}

func GetOtopProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.OtopProducts

	// Fetch the product by ID
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	// Return the product as JSON
	return c.JSON(product)
}

func GetOtopTotalQuantityName(c *fiber.Ctx) error {
	log.Println("Received request to calculate total quantity of products by name")

	var result []struct {
		ProductName   string `json:"product_name"`
		TotalQuantity int64  `json:"total_quantity"`
	}

	err := database.DB.Model(&models.OtopProducts{}).
		Select("name as product_name, SUM(quantity) as total_quantity").
		Group("name").
		Scan(&result).Error

	if err != nil {
		log.Println("Error calculating total quantity by product name:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate total quantity"})
	}

	log.Println("Total quantity calculated by product name:", result)
	return c.JSON(result)
}

func GetOtopTotalProducts(c *fiber.Ctx) error {
	log.Println("Received request to calculate total number of unique products")

	var result struct {
		TotalProducts int64 `json:"total_products"`
	}

	err := database.DB.Model(&models.OtopProducts{}).
		Select("COUNT(DISTINCT name) as total_products").
		Scan(&result).Error

	if err != nil {
		log.Println("Error calculating total number of products:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate total products"})
	}

	log.Println("Total number of products calculated:", result.TotalProducts)
	return c.JSON(result)
}

func GetTotalProductsByCategory(c *fiber.Ctx) error {
	var foodCount, nonFoodCount int64

	if err := database.DB.Model(&models.OtopProducts{}).Where("category = ?", "Food").Count(&foodCount).Error; err != nil {
		log.Println("Error counting food products:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count food products"})
	}

	if err := database.DB.Model(&models.OtopProducts{}).Where("category = ?", "Non-Food").Count(&nonFoodCount).Error; err != nil {
		log.Println("Error counting non-food products:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count non-food products"})
	}

	return c.JSON(fiber.Map{
		"Food":     foodCount,
		"Non-Food": nonFoodCount,
	})
}

package controllers

import (
	"fmt"
	"log"
	"net/http"

	// "path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/m/database"
	"github.com/m/models"
)

func AddProduct(c *fiber.Ctx) error {
	// Extract user role from JWT
	userRole, ok := c.Locals("role").(string)
	if !ok || userRole != "supplier" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User role not found"})
	}

	// Check if the user is a supplier
	if userRole != "supplier" {
		return c.Status(http.StatusForbidden).JSON(fiber.Map{"error": "Only suppliers can add products"})
	}

	// Get form data
	name := c.FormValue("name")
	log.Printf("Product Name: %s", name)
	description := c.FormValue("description")
	category := c.FormValue("category")
	log.Printf("Description: %s, Category: %s", description, category)
	price, err := strconv.ParseFloat(c.FormValue("price"), 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid price"})
	}
	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid stock"})
	}

	// // Get the file from the form field
	// file, err := c.FormFile("image")
	// if err != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Cannot process image"})
	// }

	// // Save the file temporarily in the uploads folder
	// filePath := filepath.Join("./uploads", file.Filename)
	// if err := c.SaveFile(file, filePath); err != nil {
	// 	return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot save image"})
	// }

	// Get user ID (Supplier) from JWT
	userIDInterface := c.Locals("userID")
	if userIDInterface == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found"})
	}
	userID, ok := userIDInterface.(uint)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "User ID is not of type uint"})
	}
	// Create product entry with file path stored
	product := models.Product{
		Name:        name,
		Description: description,
		Category:    category,
		Price:       price,
		Stock:       stock,
		// FilePath:    filePath, // Store the local file path in the database
		UserID: userID,
	}

	// Save product to the database
	if err := database.DB.Create(&product).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving product to database", "details": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(product)
}
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Product not found",
		})
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot update supplier",
		})
	}

	return c.JSON(product)

}

func GetProductByName(c *fiber.Ctx) error {
	name := c.Params("name")

	var product models.Product
	fmt.Println("Searching for product with name:", name)

	if err := database.DB.Where("name", name).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Product name not found in the given name",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve product name",
		})
	}

	return c.Status(fiber.StatusOK).JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")

	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": " Product bot found",
		})
	}
	if err := database.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot delete the product",
		})
	}
	return c.JSON(fiber.Map{"message": "Product delete successfully"})
}

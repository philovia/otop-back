package controllers

import (
	"log"
	"net/url"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/m/database"
	"github.com/m/models"
	"gorm.io/gorm"
)

func AddProduct(c *fiber.Ctx) error {
	var product models.Product

	userToken := c.Locals("supplier").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	supplierID := uint(claims["id"].(float64))

	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product data"})
	}

	if product.Name == "" || product.Price <= 0 || product.Quantity <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Name and Price are required"})
	}

	product.SupplierID = supplierID

	if err := database.DB.Create(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving product"})
	}

	return c.Status(fiber.StatusCreated).JSON(product)
}

func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	if err := database.DB.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}
	return c.JSON(products)
}

func GetMyProducts(c *fiber.Ctx) error {
	userToken := c.Locals("supplier").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	supplierID := uint(claims["id"].(float64))

	var products []models.Product
	if err := database.DB.Where("supplier_id = ?", supplierID).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch products"})
	}

	return c.JSON(products)
}

func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("name")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	userToken := c.Locals("supplier").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	supplierID := uint(claims["id"].(float64))

	if product.SupplierID != supplierID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized to update this product"})
	}

	var updatedProduct models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}
	if updatedProduct.Price <= 0 || updatedProduct.Quantity < 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Price must be greater than 0 and Quantity cannot be negative"})
	}

	product.Name = updatedProduct.Name
	product.Price = updatedProduct.Price
	product.Quantity = updatedProduct.Quantity

	if err := database.DB.Save(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	return c.JSON(product)
}

func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
	}

	userToken := c.Locals("supplier").(*jwt.Token)
	claims := userToken.Claims.(jwt.MapClaims)
	supplierID := uint(claims["id"].(float64))

	if product.SupplierID != supplierID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Not authorized to delete this product"})
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

func GetProductByName(c *fiber.Ctx) error {
	Name := c.Params("name")

	log.Println("Looking for product with name:", Name)

	decodedName, err := url.QueryUnescape(Name)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid product name"})
	}
	Name = strings.TrimSpace(decodedName)

	var product models.Product

	if err := database.DB.Where("LOWER(name) = LOWER(?)", Name).First(&product).Error; err != nil {
		log.Println("Error fetching product:", err)
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch product"})
	}

	return c.JSON(product)
}

func GetTotalQuantity(c *fiber.Ctx) error {
	log.Println("Received request to calculate total quantity of products")

	var totalQuantity int64
	err := database.DB.Model(&models.Product{}).Select("SUM(quantity)").Scan(&totalQuantity).Error

	if err != nil {
		log.Println("Error calculating total quantity:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate total quantity"})
	}

	log.Println("Total quantity calculated:", totalQuantity)
	return c.JSON(fiber.Map{"total_quantity": totalQuantity})
}

func GetTotalPrice(c *fiber.Ctx) error {
	log.Println("Received request to calculate total price of all products")

	var totalPrice float64
	err := database.DB.Model(&models.Product{}).Select("SUM(price * quantity)").Scan(&totalPrice).Error

	if err != nil {
		log.Println("Error calculating total price:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to calculate total price"})
	}

	log.Println("Total price calculated:", totalPrice)
	return c.JSON(fiber.Map{"total_price": totalPrice})
}

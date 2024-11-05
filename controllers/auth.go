package controllers

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m/database"
	"github.com/m/models"
	"github.com/m/utils"
)

func Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid user data"})
	}
	if user.Role == "admin" && !strings.HasSuffix(user.UserName, "_admin") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Admin username must end with '_admin'"})
	}
	if user.Role == "cashier" && !strings.HasSuffix(user.UserName, "_cashier") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cashier username must end with '_cashier'"})
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error saving user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

func UnifiedLogin(c *fiber.Ctx) error {
	var supplierCreds models.Supplier
	if err := c.BodyParser(&supplierCreds); err == nil && supplierCreds.StoreName != "" {
		return SupplierLogin(c)
	}

	var userCreds models.User
	if err := c.BodyParser(&userCreds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid login data"})
	}

	return Login(c)
}

func Login(c *fiber.Ctx) error {
	var creds models.User
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid login data"})
	}

	var storedUser models.User
	if err := database.DB.Where("email = ?", creds.Email).First(&storedUser).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	}

	token, err := utils.GenerateToken(storedUser.UserName, storedUser.Role, storedUser.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	return c.JSON(fiber.Map{"token": token, "role": storedUser.Role})
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Logout successful, please clear your token",
	})
}

func SupplierLogin(c *fiber.Ctx) error {
	var creds models.Supplier
	if err := c.BodyParser(&creds); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid login data"})
	}

	var storedSupplier models.Supplier
	if err := database.DB.Where("store_name = ?", creds.StoreName).First(&storedSupplier).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Supplier not found"})
	}

	if storedSupplier.Password != creds.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateToken(storedSupplier.StoreName, "supplier", storedSupplier.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Error generating token"})
	}

	return c.JSON(fiber.Map{"token": token, "role": "supplier"})
}

func SupplierLogout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Supplier logout successful, please clear your token",
	})
}

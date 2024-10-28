package controllers

import (
	// "fmt"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/m/database"
	"github.com/m/models"
	"github.com/m/utils"
	"gorm.io/gorm"
	// "golang.org/x/crypto/bcrypt"
)

func CreateSupplier(c *fiber.Ctx) error {
	var supplier models.Supplier

	// Parse the request body for supplier details (email, password, etc.)
	if err := c.BodyParser(&supplier); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Hash the password before saving (you should already have a function for this)
	hashedPassword, err := utils.HashPassword(supplier.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error hashing password",
		})
	}
	supplier.Password = hashedPassword
	supplier.Role = "supplier"

	if err := database.DB.Create(&supplier).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not create supplier account",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "Supplier account created successfully",
		"supplier": supplier,
	})
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

func GetAllSuppliers(c *fiber.Ctx) error {
	var suppliers []models.Supplier

	if err := database.DB.Find(&suppliers).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch suppliers",
		})
	}

	return c.JSON(suppliers)
}

func GetSupplierByEmail(c *fiber.Ctx) error {
	email := c.Params("email")

	var supplier models.Supplier

	if err := database.DB.Where("LOWER(email) = LOWER(?)", email).First(&supplier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Supplier not found for the given email",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve supplier",
		})
	}

	return c.Status(fiber.StatusOK).JSON(supplier)
}
func GetSupplierByStoreName(c *fiber.Ctx) error {
	storeName := c.Params("storeName")

	var supplier models.Supplier
	fmt.Println("Searching for supplier with store name:", storeName)

	if err := database.DB.Where("LOWER(store_name) = LOWER(?)", storeName).First(&supplier).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "Supplier not found for the given store name",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not retrieve supplier",
		})
	}

	return c.Status(fiber.StatusOK).JSON(supplier)
}

func SupplierLogin(c *fiber.Ctx) error {
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the login request body (email and password)
	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request format",
		})
	}

	var supplier models.User
	// Fetch the supplier from the database by email and ensure the role is 'supplier'
	if err := database.DB.Where("email = ? AND role = ?", loginRequest.Email, "supplier").First(&supplier).Error; err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Compare the hashed password with the provided password
	if !utils.CheckPasswordHash(loginRequest.Password, supplier.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	// Generate a JWT token for the supplier with userID, email, and role
	token, err := utils.GenerateJWT(supplier.Email, supplier.Role, uint(supplier.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not generate token",
		})
	}

	// Return the token and supplier information
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Supplier logged in successfully",
		"token":   token,
		"user": fiber.Map{
			"id":    supplier.ID,
			"email": supplier.Email,
			"role":  supplier.Role,
		},
	})
}

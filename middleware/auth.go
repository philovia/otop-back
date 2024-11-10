package middleware

import (
	// "os"÷\\\
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/m/utils"
)

func JWTProtected(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	token, err := utils.ParseToken(tokenString)
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized Pleased Input Token"})
	}
	c.Locals("user", token)
	return c.Next()
}
func IsAdmin(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	userToken, err := utils.ParseToken(tokenStr)
	if err != nil || !userToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized No token"})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized  Can't claim token"})
	}

	// Check if the user has the 'admin' role
	if role, exists := claims["role"].(string); !exists || role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	// Store claims in Locals for access in handlers
	c.Locals("user", claims)

	// Return nil if the check passes
	return nil
}

func IsSupplier(c *fiber.Ctx) error {
	// Retrieve the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	// Remove the "Bearer " prefix from the token
	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

	// Parse the token using your ParseToken function
	userToken, err := utils.ParseToken(tokenStr)
	if err != nil || !userToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized. Invalid token"})
	}

	// Parse the claims into your custom Claims struct
	claims, ok := userToken.Claims.(*utils.Claims)
	if !ok {
		// This will handle the case where the claims cannot be extracted correctly
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized. Unable to extract claims"})
	}

	// Ensure the role is 'supplier'
	if claims.Role != "supplier" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden. You are not authorized to perform this action"})
	}

	// Set the supplier_id in Locals for further use in handlers
	c.Locals("supplier_id", claims.SupplierID)

	// Proceed to the next handler
	return c.Next()
}

func IsCashier(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized access"})
	}

	tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
	userToken, err := utils.ParseToken(tokenStr)
	if err != nil || !userToken.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := userToken.Claims.(jwt.MapClaims)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	if role, exists := claims["role"].(string); !exists || role != "cashier" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	c.Locals("cashier", userToken)
	return c.Next()
}

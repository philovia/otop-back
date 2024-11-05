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
	userToken := c.Locals("user")
	if userToken == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	claims, ok := userToken.(*jwt.Token).Claims.(jwt.MapClaims)
	if !ok || claims["role"] != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	return c.Next()
}
func IsSupplier(c *fiber.Ctx) error {
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

	if role, exists := claims["role"].(string); !exists || role != "supplier" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Forbidden"})
	}

	c.Locals("supplier", userToken)
	return c.Next()
}

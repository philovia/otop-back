package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m/utils"
)

// var tokenBlacklist = make(map[string]bool)

func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authFormHeader := c.Get("Authorization")
		if authFormHeader == "" {
			log.Println("Unauthorized: No token provided")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized, no token provided",
			})
		}

		tokenString := strings.TrimPrefix(authFormHeader, "Bearer ")
		if tokenString == "" {
			log.Println("Unauthorized: No token found in header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized, no token found in header",
			})
		}

		// if _, exists := tokenBlacklist[tokenString]; exists {
		// 	log.Println("Unauthorized: Token is invalid (logged out)")
		// 	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		// 		"error": "Unauthorized, token is invalid",
		// 	})

		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			log.Println("Unauthorized: Invalid token", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized, invalid token",
			})
		}
		log.Printf("Authenticated user: %s with role: %s", claims.Username, claims.Role)

		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}

func Logout(c *fiber.Ctx) error {
	authFormHeader := c.Get("Authorization")
	if authFormHeader == "" {
		log.Println("Logout failed: No token provided")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token provided",
		})
	}

	tokenString := strings.TrimPrefix(authFormHeader, "Bearer ")
	if tokenString == "" {
		log.Println("Logout failed: No token found in header")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "No token found in header",
		})
	}

	// Invalidate the token by adding it to the blacklist
	// tokenBlacklist[tokenString] = true
	log.Println("Token invalidated for logout:", tokenString)

	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

func IsAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		if role != "admin" {
			log.Println("Forbidden: User is not an admin")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only admins can access this resource",
			})
		}

		return c.Next()
	}
}

func IsSupplier() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		if role != "supplier" {
			log.Println("Forbidden: User is not a supplier")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only suppliers can access this resource",
			})
		}

		return c.Next()
	}

}

func IsCashier() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role").(string)

		if role != "Cashier" {
			log.Println("Forbidden: User is not cashier")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only Cashiers can access this resource",
			})
		}
		return c.Next()
	}

}

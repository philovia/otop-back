package middleware

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/m/utils" // For JWT validation logic
)

// In-memory blacklist to store invalidated tokens
var tokenBlacklist = make(map[string]bool)

// Authentication middleware to verify JWT token
func Authentication() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authFormHeader := c.Get("Authorization")

		// Ensure the token is provided
		if authFormHeader == "" {
			log.Println("Unauthorized: No token provided")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized, no token provided",
			})
		}

		// Extract the token from the header (Bearer <token>)
		tokenString := strings.TrimPrefix(authFormHeader, "Bearer ")
		if tokenString == "" {
			log.Println("Unauthorized: No token found in header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized, no token found in header",
			})
		}

		// Check if the token is in the blacklist
		if _, exists := tokenBlacklist[tokenString]; exists {
			log.Println("Unauthorized: Token is invalid (logged out)")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized, token is invalid",
			})
		}

		// Validate the token
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			log.Println("Unauthorized: Invalid token", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized, invalid token",
			})
		}

		// Add user info (username and role) to the context
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)
		c.Locals("userID", claims.UserID)

		// Proceed to the next middleware or handler
		return c.Next()
	}
}

// IsSupplier middleware to verify if the user is a supplier
func IsSupplier() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the role from the request context (set in the Authentication middleware)
		role := c.Locals("role").(string)

		// Check if the user has the role "supplier"
		if role != "supplier" {
			log.Println("Forbidden: User is not a supplier")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Only suppliers can access this resource",
			})
		}

		// Proceed if the user is a supplier
		return c.Next()
	}
}

// Logout function to invalidate the token
func Logout(c *fiber.Ctx) error {
	// Extract the token from the header (Bearer <token>)
	authFormHeader := c.Get("Authorization")
	tokenString := strings.TrimPrefix(authFormHeader, "Bearer ")

	// Add the token to the blacklist
	tokenBlacklist[tokenString] = true

	log.Println("Token invalidated for logout:", tokenString)
	return c.JSON(fiber.Map{"message": "Logged out successfully"})
}

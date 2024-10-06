package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	UserID   uint   `json:"user_id"` // Add UserID field
	jwt.RegisteredClaims
}

// GenerateJWT generates a new JWT token for a user with their username and role
func GenerateJWT(username, role string, userID uint) (string, error) {
	// Set expiration time for the token (24 hours)
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, including the username and role
	claims := &Claims{
		Username: username,
		Role:     role,
		UserID:   userID, // Set UserID
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	// Create the token with the claims and signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token and returns the claims if valid
func ValidateToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	// Parse the token with claims and the signing key
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Handle token validation errors
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("invalid token signature")
		}
		return nil, errors.New("invalid token")
	}

	// Check if the token is valid
	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return claims, nil
}

// ExtractRoleAndUsername extracts the role and username from a valid token
func ExtractRoleAndUsername(tokenString string) (string, string, error) {
	claims, err := ValidateToken(tokenString)
	if err != nil {
		return "", "", err
	}
	return claims.Role, claims.Username, nil
}

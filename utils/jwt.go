package utils

import (
<<<<<<< HEAD
	"fmt"
=======
>>>>>>> 36cf5b4b0c38771a532201f6a055694672691442
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtKey = []byte("your_secret_key")

type Claims struct {
<<<<<<< HEAD
	UserName   string `json:"user_name"`
	Role       string `json:"role"`
	ID         uint   `json:"id"`
	SupplierID uint   `json:"supplier_id"`
	jwt.RegisteredClaims
}

func GenerateToken(username string, role string, ID uint, supplierID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserName:   username,
		Role:       role,
		ID:         ID,
		SupplierID: supplierID,
=======
	UserName string `json:"user_name"`
	Role     string `json:"role"`
	ID       uint   `json:"id"`
	jwt.RegisteredClaims
}

func GenerateToken(username string, role string, ID uint) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserName: username,
		Role:     role,
		ID:       ID,
>>>>>>> 36cf5b4b0c38771a532201f6a055694672691442
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseToken(tokenStr string) (*jwt.Token, error) {
<<<<<<< HEAD
	return jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
=======
	if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
		tokenStr = tokenStr[7:]
	}

	return jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
>>>>>>> 36cf5b4b0c38771a532201f6a055694672691442
		return jwtKey, nil
	})
}

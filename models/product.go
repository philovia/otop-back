package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name" form:"name"`
	Description string  `json:"description" form:"description"`
	Category    string  `json:"category" form:"category"`
	Price       float64 `json:"price" form:"price"`
	Stock       int     `json:"tock" form:"stock"`
	UserID      uint    `json:"user_id"`
}

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
	FilePath    string  `json:"ile_path" form:"file_path"` // This stores the local file path
	UserID      uint    `json:"user_id"`                   // Foreign key to User (Supplier)
}

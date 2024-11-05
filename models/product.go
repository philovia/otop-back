package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	SupplierID  uint    `json:"supplier_id"`
	Category    string  `json:"category"`
}

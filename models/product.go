package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
<<<<<<< HEAD
	ID          uint    `gorm:"primaryKey" json:"id"`
=======
>>>>>>> 36cf5b4b0c38771a532201f6a055694672691442
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	SupplierID  uint    `json:"supplier_id"`
	Category    string  `json:"category"`
}

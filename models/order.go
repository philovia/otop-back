package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	ID          uint      `gorm:"primaryKey" json:"id"`
	AdminID     uint      `json:"admin_id"`
	SupplierID  uint      `json:"supplier_id"`
	ProductID   uint      `json:"product_id"`
	ProductName string    `json:"product_name"`
	Quantity    int64     `json:"quantity"`
	Price       float64   `json:"price"`
	OrderDate   time.Time `json:"order_date"`
	Status      string    `json:"status"`
	Descriptiom string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

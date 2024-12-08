package models

import (
	"errors"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	SupplierID  uint    `json:"supplier_id"`
	Category    string  `json:"category"`
}

// BeforeCreate hook to validate the category field
func (p *Product) BeforeCreate(tx *gorm.DB) (err error) {
	if p.Category != "Food" && p.Category != "Non-Food" {
		return errors.New("category must be either 'Food' or 'Non-Food'")
	}
	return
}

// BeforeUpdate hook to validate the category field during an update
func (p *Product) BeforeUpdate(tx *gorm.DB) (err error) {
	if p.Category != "Food" && p.Category != "Non-Food" {
		return errors.New("category must be either 'Food' or 'Non-Food'")
	}
	return
}

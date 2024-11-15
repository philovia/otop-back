package models

import (
	"errors"

	"gorm.io/gorm"
)

type OtopProducts struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int64   `json:"quantity"`
	Category    string  `json:"category"`
}

// BeforeCreate hook to validate the category field
func (p *OtopProducts) BeforeCreate(tx *gorm.DB) (err error) {
	if p.Category != "Food" && p.Category != "Non-Food" {
		return errors.New("category must be either 'Food' or 'Non-Food'")
	}
	return
}

// BeforeUpdate hook to validate the category field during an update
func (p *OtopProducts) BeforeUpdate(tx *gorm.DB) (err error) {
	if p.Category != "Food" && p.Category != "Non-Food" {
		return errors.New("category must be either 'Food' or 'Non-Food'")
	}
	return
}

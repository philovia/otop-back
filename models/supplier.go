package models

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
	StoreName   string `json:"store_name" gorm:"unique"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	Role        string `gorm:"default:'supplier'"`
}

package models

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
	ID          int    `json:"id"`
	StoreName   string `json:"storename"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phonenumber"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	// Role     string `json:"role"`
}

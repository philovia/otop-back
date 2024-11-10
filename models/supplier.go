package models

import "gorm.io/gorm"

type Supplier struct {
	gorm.Model
<<<<<<< HEAD
	ID          uint   `json:"id" gorm:"primary_key"`
=======
>>>>>>> 36cf5b4b0c38771a532201f6a055694672691442
	StoreName   string `json:"store_name" gorm:"unique"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Address     string `json:"address"`
	Password    string `json:"password"`
	Role        string `gorm:"default:'supplier'"`
}

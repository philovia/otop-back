package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       uint   `json:"id" gorm:"primary_key"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

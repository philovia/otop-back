package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       int    `json:"id"`
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	UserID   uint   `json:"userid"`
}

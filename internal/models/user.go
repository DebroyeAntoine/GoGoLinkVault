package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email"`
	Password string `json:"password"`
	JWTToken string `json:"jwt_token"`
	Links    []Link `gorm:"foreignKey:UserID"`
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID           uint   `json:"id" gorm:"primary_key" example:"999"`
	EmailAddress string `json:"emailAddress" gorm:"size:255;not null;unique" example:"johndoe@gmail.com"`
	Password     string `json:"-" gorm:"size:255;not null;" binding:"required" example:"password"`
}

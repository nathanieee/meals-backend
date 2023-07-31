package models

import (
	"gorm.io/gorm"
)

type (
	User struct {
		gorm.Model
		ID           uint   `json:"id" gorm:"primary_key" example:"999"`
		EmailAddress string `json:"emailAddress" gorm:"size:255;not null;unique" example:"johndoe@gmail.com"`
		Password     string `json:"-" gorm:"size:255;not null;" binding:"required" example:"password"`
		FullName     string `json:"fullName" gorm:"not null" example:"user name"`
		Email        string `json:"email" gorm:"not null;unique" example:"email@email.com"`
		UserLevel    uint   `json:"userLevel" gorm:"not null" example:"1"`
		Country      string `json:"country" example:"country"`
		CountryCode  uint   `json:"countryCode" example:"62"`
	}
)

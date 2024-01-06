package models

import (
	"project-skbackend/internal/models/helper"
)

type (
	Address struct {
		helper.Model
		Name        string  `json:"name" gorm:"not null" binding:"required"`
		Address     string  `json:"address" gorm:"not null" binding:"required"`
		Description string  `json:"description" gorm:"not null" binding:"required"`
		Note        string  `json:"note" gorm:"not null" binding:"required"`
		Landmark    string  `json:"landmark" gorm:"not null" binding:"required"`
		Longitude   float64 `json:"langitude" gorm:"not null" binding:"required"`
		Latitude    float64 `json:"latitude" gorm:"not null" binding:"required"`
	}
)

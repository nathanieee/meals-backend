package models

import (
	"project-skbackend/internal/models/helper"

	"github.com/google/uuid"
)

type (
	Rating struct {
		helper.Model
		MealID      uuid.UUID `json:"meal_id" gorm:"not null" binding:"required"`
		Meal        Meal      `json:"meal"`
		UserID      uuid.UUID `json:"user_id" gorm:"not null" binding:"required"`
		User        User      `json:"user"`
		Value       float64   `json:"value" gorm:"not null" binding:"required"`
		Description string    `json:"description,omitempty"`
	}
)

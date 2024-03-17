package models

import (
	"project-skbackend/internal/models/base"

	"github.com/google/uuid"
)

type (
	Rating struct {
		base.Model

		MealID uuid.UUID `json:"meal_id" gorm:"required"`
		Meal   Meal      `json:"meal"`

		UserID uuid.UUID `json:"user_id" gorm:"required"`
		User   User      `json:"user"`

		Value       float64 `json:"value" gorm:"required"`
		Description string  `json:"description,omitempty"`
	}
)

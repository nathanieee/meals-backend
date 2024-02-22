package models

import (
	"project-skbackend/internal/models/helper"

	"github.com/google/uuid"
)

type (
	Cart struct {
		helper.Model

		MealID        uuid.UUID `json:"meal_id" gorm:"required"`
		Meal          Meal      `json:"meal"`
		ReferenceID   uuid.UUID `json:"reference_id" gorm:"required"`
		ReferenceType string    `json:"referenceType" gorm:"required; oneof='Member' 'Caregiver'"`
		Quantity      uint      `json:"quantity" gorm:"required"`
	}
)

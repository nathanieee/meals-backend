package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

type (
	Cart struct {
		helper.Model

		MealID        uuid.UUID           `json:"meal_id" gorm:"required"`
		Meal          Meal                `json:"meal"`
		ReferenceID   uuid.UUID           `json:"reference_id" gorm:"required"`
		ReferenceType consttypes.UserRole `json:"referenceType" gorm:"required; oneof='Member' 'Caregiver';type:user_role_enum;"`
		Quantity      uint                `json:"quantity" gorm:"required"`
	}
)

package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	Allergy struct {
		helper.Model
		Name        string               `json:"name" gorm:"not null" binding:"required" example:"Milk"`
		Description string               `json:"description" gorm:"not null" binding:"required" example:"A milk allergy, also known as a dairy allergy, is an adverse immune system response to one or more proteins found in cow's milk. It is different from lactose intolerance, which is a non-immune digestive disorder where the body has difficulty digesting lactose, a sugar found in milk. A milk allergy is an immune system disorder and can be more severe."`
		Allergens   consttypes.Allergens `json:"allergens" gorm:"not null" binding:"required" example:"Food"`
	}
)

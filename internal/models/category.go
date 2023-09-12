package models

import "project-skbackend/internal/models/helper"

type (
	FoodCategory struct {
		helper.Model
		Name string `json:"name" gorm:"not null" binding:"required" example:"Rice"`
	}
)

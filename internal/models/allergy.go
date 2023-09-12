package models

import "project-skbackend/internal/models/helper"

type (
	Allergy struct {
		helper.Model
		Name string `json:"name" gorm:"not null" binding:"required" example:"Lactose Intolerant"`
	}
)

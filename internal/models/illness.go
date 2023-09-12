package models

import "project-skbackend/internal/models/helper"

type (
	Illness struct {
		helper.Model
		Name string `json:"name" gorm:"not null" binding:"required" example:"Cold Sore"`
	}
)

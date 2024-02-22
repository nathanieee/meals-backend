package models

import (
	"project-skbackend/internal/models/helper"

	"github.com/google/uuid"
)

type (
	Address struct {
		helper.Model
		UserID    uuid.UUID `json:"user_id" gorm:"required"`
		Name      string    `json:"name" gorm:"required"`
		Address   string    `json:"address" gorm:"required"`
		Note      string    `json:"note" gorm:"required"`
		Longitude string    `json:"longitude" gorm:"required;longitude"`
		Latitude  string    `json:"latitude" gorm:"required;latitude"`
	}
)

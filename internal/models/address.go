package models

import (
	"project-skbackend/internal/models/base"

	"github.com/google/uuid"
)

type (
	Address struct {
		base.Model

		UserID uuid.UUID `json:"user_id" gorm:"required"`
		User   User      `json:"user"`

		Name      string `json:"name" gorm:"required"`
		Address   string `json:"address" gorm:"required"`
		Note      string `json:"note" gorm:"required;max:256"`
		Longitude string `json:"longitude" gorm:"required;longitude"`
		Latitude  string `json:"latitude" gorm:"required;latitude"`
	}
)

package models

import (
	"project-skbackend/internal/models/helper"

	"github.com/google/uuid"
)

type (
	Address struct {
		helper.Model
		UserID    uuid.UUID `json:"user_id" gorm:"not null" binding:"required"`
		Name      string    `json:"name" gorm:"not null" binding:"required"`
		Address   string    `json:"address" gorm:"not null" binding:"required"`
		Note      string    `json:"note" gorm:"not null" binding:"required"`
		Longitude string    `json:"longitude" gorm:"not null;longitude" binding:"required"`
		Latitude  string    `json:"latitude" gorm:"not null;latitude" binding:"required"`
	}
)

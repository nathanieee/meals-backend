package models

import "project-skbackend/internal/models/helper"

type (
	Partner struct {
		helper.Model
		UserID int    `json:"user_id" gorm:"not null" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		User   User   `json:"user"`
		Name   string `json:"name" gorm:"not null" binding:"required" example:"McDonald's"`
	}
)

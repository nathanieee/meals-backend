package requests

import "github.com/google/uuid"

type (
	CreateAddress struct {
		Name      string `json:"name" form:"name" binding:"required"`
		Address   string `json:"address" form:"address" binding:"required"`
		Note      string `json:"note" form:"note" binding:"required;max:32"`
		Longitude string `json:"longitude" form:"longitude" binding:"required"`
		Latitude  string `json:"latitude" form:"latitude" binding:"required"`
	}

	UpdateAddress struct {
		ID        uuid.UUID `json:"id" form:"id" binding:"required" example:"f7fbfa0d-5f95-42e0-839c-d43f0ca757a4"`
		Name      string    `json:"name" form:"name" binding:"required"`
		Address   string    `json:"address" form:"address" binding:"required"`
		Note      string    `json:"note" form:"note" binding:"required"`
		Longitude string    `json:"longitude" form:"longitude" binding:"required"`
		Latitude  string    `json:"latitude" form:"latitude" binding:"required"`
	}
)

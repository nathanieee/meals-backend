package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	Image struct {
		helper.Model
		Name string               `json:"name" gorm:"not null" binding:"required" example:"image.jpg"`
		Path string               `json:"path" gorm:"not null" binding:"required" example:"./files/images/profiles/image.jpg"`
		Type consttypes.ImageType `json:"image_type" gorm:"not null; type:image_type_enum" binding:"required" example:"Profile"`
	}
)
package models

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	Image struct {
		helper.Model
		Name string               `json:"name" gorm:"required" example:"image.jpg"`
		Path string               `json:"path" gorm:"required" example:"./files/images/profiles/image.jpg"`
		Type consttypes.ImageType `json:"image_type" gorm:"required; type:image_type_enum" example:"Profile"`
	}
)

package responses

import (
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
)

type (
	Image struct {
		helper.Model
		Name string               `json:"name"`
		Path string               `json:"path"`
		Type consttypes.ImageType `json:"image_type"`
	}
)

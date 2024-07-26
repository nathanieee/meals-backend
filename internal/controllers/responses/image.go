package responses

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Image struct {
		base.Model

		Name string               `json:"name"`
		Path string               `json:"path"`
		Type consttypes.ImageType `json:"image_type"`
	}
)

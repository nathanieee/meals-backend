package responses

import "project-skbackend/internal/models/helper"

type (
	Illness struct {
		helper.Model
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

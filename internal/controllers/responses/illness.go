package responses

import "project-skbackend/internal/models/base"

type (
	Illness struct {
		base.Model

		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

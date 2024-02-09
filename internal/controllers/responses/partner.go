package responses

import "project-skbackend/internal/models/helper"

type (
	Partner struct {
		helper.Model
		User User   `json:"user"`
		Name string `json:"name"`
	}
)

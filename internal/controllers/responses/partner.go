package responses

import "project-skbackend/internal/models/base"

type (
	Partner struct {
		base.Model

		User User `json:"user"`

		Name string `json:"name"`
	}
)

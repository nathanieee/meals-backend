package responses

import "project-skbackend/internal/models/helper"

type (
	IllnessResponse struct {
		helper.Model
		Name        string `json:"name"`
		Description string `json:"description"`
	}
)

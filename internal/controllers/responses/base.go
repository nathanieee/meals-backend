package responses

import "project-skbackend/packages/consttypes"

type (
	BaseRole struct {
		Data any                 `json:"data"`
		Role consttypes.UserRole `json:"role"`
	}
)

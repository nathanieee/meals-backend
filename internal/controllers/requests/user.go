package requests

import "project-skbackend/packages/consttypes"

type (
	CreateUserRequest struct {
		Email    string              `json:"email" binding:"required" example:"email@email.com"`
		Role     consttypes.UserRole `json:"role" default:"0" example:"0"`
		Password string              `json:"-" binding:"required" example:"password"`
	}
)

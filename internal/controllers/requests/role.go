package requests

import "project-skbackend/packages/consttypes"

type (
	CreateRoleRequest struct {
		Name    consttypes.Role  `json:"name,omitempty" binding:"required, unique"`
		LevelID consttypes.Level `json:"levelID,omitempty" binding:"required"`
	}
)

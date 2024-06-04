package utrole

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
)

func CartRoleCheck(role responses.BaseRole) (uuid.UUID, consttypes.UserRole, bool) {
	switch role.Role {
	case consttypes.UR_CAREGIVER:
		res, ok := role.Data.(*models.Caregiver)
		if !ok {
			return uuid.UUID{}, consttypes.UserRole(0), false
		}
		return res.ID, res.User.Role, true
	case consttypes.UR_MEMBER:
		res, ok := role.Data.(*models.Member)
		if !ok {
			return uuid.UUID{}, consttypes.UserRole(0), false
		}
		return res.ID, res.User.Role, true
	default:
		return uuid.UUID{}, consttypes.UserRole(0), false
	}
}

package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateUserRequest struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	UpdateUserRequest struct {
		Email   string               `json:"email" binding:"required,email" example:"email@email.com"`
		Address CreateAddressRequest `json:"address" binding:"required"`
	}
)

func (req *CreateUserRequest) ToModel(role consttypes.UserRole) *models.User {
	user := models.User{
		Role: role,
	}

	if err := copier.Copy(&user, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &user
}

func (req *UpdateUserRequest) ToModel(role consttypes.UserRole, uid uuid.UUID) *models.User {
	user := models.User{
		Role: role,
	}

	if err := copier.Copy(&user, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &user
}

package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateUser struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	UpdateUser struct {
		Email   string        `json:"email" binding:"required,email" example:"email@email.com"`
		Address CreateAddress `json:"address" binding:"required"`
	}
)

func (req *CreateUser) ToModel(role consttypes.UserRole) *models.User {
	user := models.User{
		Role: role,
	}

	if err := copier.Copy(&user, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &user
}

func (req *UpdateUser) ToModel(role consttypes.UserRole, uid uuid.UUID) *models.User {
	user := models.User{
		Role: role,
	}

	if err := copier.Copy(&user, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &user
}

package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/logger"

	"github.com/jinzhu/copier"
)

type (
	CreateUserRequest struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}
)

func (req *CreateUserRequest) ToModel(role consttypes.UserRole) *models.User {
	user := models.User{
		Role: role,
	}

	if err := copier.Copy(&user, &req); err != nil {
		logger.LogError(err)
		return nil
	}

	return &user
}

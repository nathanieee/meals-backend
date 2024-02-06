package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateUser struct {
		Image    *CreateImage     `json:"image"`
		Email    string           `json:"email" binding:"required,email" example:"email@email.com"`
		Password string           `json:"password" binding:"required" example:"password"`
		Address  *[]CreateAddress `json:"address"`
	}

	UpdateUser struct {
		Image    *UpdateImage     `json:"image"`
		Email    string           `json:"email" binding:"required,email" example:"email@email.com"`
		Password string           `json:"password" binding:"required" example:"password"`
		Address  *[]UpdateAddress `json:"address"`
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
		Model: helper.Model{ID: uid},
		Role:  role,
	}

	if err := copier.Copy(&user, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &user
}

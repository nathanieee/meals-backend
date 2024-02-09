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
		UserImage *CreateImage     `json:"user_image" form:"user_image" binding:"-"`
		Email     string           `json:"email" form:"email" binding:"required,email"`
		Password  string           `json:"password" form:"password" binding:"required"`
		Address   *[]CreateAddress `json:"address" form:"address" binding:"-"`
	}

	UpdateUser struct {
		UserImage *CreateImage     `json:"user_image" form:"user_image" binding:"-"`
		Email     string           `json:"email" form:"email" binding:"required,email"`
		Password  string           `json:"password" form:"password" binding:"required"`
		Address   *[]UpdateAddress `json:"address" form:"address" binding:"-"`
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

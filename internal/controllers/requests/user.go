package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreateUser struct {
		*CreateImage
		Email           string           `json:"email" form:"email" binding:"required,email"`
		Password        string           `json:"password" form:"password" binding:"required"`
		ConfirmPassword string           `json:"confirm_password" form:"confirm_password" binding:"eqfield=Password"`
		Address         *[]CreateAddress `json:"address" form:"address" binding:"-"`
	}

	UpdateUser struct {
		*UpdateImage
		Email           string           `json:"email,omitempty" form:"email" binding:"email,omitempty"`
		Password        string           `json:"password,omitempty" form:"password" binding:"-"`
		ConfirmPassword string           `json:"confirm_password" form:"confirm_password" binding:"eqfield=Password"`
		Address         *[]UpdateAddress `json:"address,omitempty" form:"address" binding:"-"`
	}
)

func (req *CreateUser) ToModel(role consttypes.UserRole) *models.User {
	var user models.User

	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil
	}

	if err := copier.CopyWithOption(&user, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil
	}

	user.Role = role
	user.Password = hash

	return &user
}

func (req *UpdateUser) ToModel(user models.User, role consttypes.UserRole) *models.User {
	hash, err := helper.HashPassword(req.Password)
	if err != nil {
		return nil
	}

	if err := copier.CopyWithOption(&user, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil
	}

	user.Role = role
	user.Password = hash

	return &user
}

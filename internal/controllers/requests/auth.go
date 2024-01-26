package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	Login struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	Register struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	VerifyToken struct {
		Token string `json:"token" binding:"required,number" example:""`
		Email string `json:"email" binding:"required,email" example:"email@email.com"`
	}

	ForgotPassword struct {
		Email string `json:"email" binding:"required,email" example:"email@email.com"`
	}

	ResetPassword struct {
		Email           string `json:"email" binding:"required,email" example:"email@email.com"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
		Token           string `json:"token" binding:"required"`
	}

	ResetPasswordRedirect struct {
		ResetToken string `uri:"token" binding:"required"`
	}
)

func (r *Register) ToUserModel() *models.User {
	var u models.User
	err := copier.Copy(&u, &r)
	if err != nil {
		utlogger.LogError(err)
	}

	return &u
}

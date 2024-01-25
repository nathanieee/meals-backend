package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	LoginRequest struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	RegisterRequest struct {
		Email    string `json:"email" binding:"required,email" example:"email@email.com"`
		Password string `json:"password" binding:"required" example:"password"`
	}

	VerifyTokenRequest struct {
		Token string `json:"token" binding:"required,number" example:""`
		Email string `json:"email" binding:"required,email" example:"email@email.com"`
	}

	ForgotPasswordRequest struct {
		Email string `json:"email" binding:"required,email" example:"email@email.com"`
	}

	ResetPasswordRequest struct {
		Email           string `json:"email" binding:"required,email" example:"email@email.com"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
		Token           string `json:"token" binding:"required"`
	}

	ResetPasswordRedirectRequest struct {
		ResetToken string `uri:"token" binding:"required"`
	}
)

func (r *RegisterRequest) ToUserModel() *models.User {
	var u models.User
	err := copier.Copy(&u, &r)
	if err != nil {
		utlogger.LogError(err)
	}

	return &u
}

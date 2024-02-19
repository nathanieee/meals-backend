package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	Login struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	Register struct {
		Email    string `json:"email" form:"email" binding:"required,email"`
		Password string `json:"password" form:"password" binding:"required"`
	}

	VerifyToken struct {
		Token string `json:"token" form:"token" binding:"required"`
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	ForgotPassword struct {
		Email string `json:"email" form:"email" binding:"required,email"`
	}

	ResetPassword struct {
		Email           string `json:"email" form:"email" binding:"required,email"`
		Password        string `json:"password" form:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`
		Token           string `json:"token" form:"token" binding:"required"`
	}

	ResetPasswordRedirect struct {
		ResetToken string `uri:"token" binding:"required"`
	}
)

func (r *Register) ToUserModel() *models.User {
	var u models.User
	err := copier.CopyWithOption(&u, &r, copier.Option{IgnoreEmpty: true, DeepCopy: true})
	if err != nil {
		utlogger.LogError(err)
	}

	return &u
}

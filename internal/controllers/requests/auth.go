package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utstring"

	"github.com/jinzhu/copier"
)

type (
	Signin struct {
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

func (r *ResetPassword) ToUserModel(
	u models.User,
) (*models.User, error) {
	hashpass, err := utstring.HashPassword(r.Password)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	if err := copier.CopyWithOption(&u, &r, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	u.Password = hashpass
	u.ResetPasswordToken = ""

	return &u, nil
}

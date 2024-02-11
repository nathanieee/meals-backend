package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreatePartner struct {
		User CreateUser `json:"user" form:"user" binding:"required"`
		Name string     `json:"name" form:"name" binding:"required"`
	}

	UpdatePartner struct {
		User UpdateUser `json:"user" form:"user" binding:"required"`
		Name string     `json:"name" form:"name" binding:"required"`
	}
)

func (req *CreatePartner) ToModel(
	user models.User,
) *models.Partner {
	partner := models.Partner{
		User: user,
	}

	if err := copier.Copy(&partner, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &partner
}

func (req *UpdatePartner) ToModel(
	partner models.Partner,
	user models.User,
) *models.Partner {
	if err := copier.CopyWithOption(&partner, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil
	}

	partner.User = user

	return &partner
}

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
		User UpdateUser `json:"user" form:"user" binding:"-"`
		Name string     `json:"name" form:"name" binding:"-"`
	}
)

func (req *CreatePartner) ToModel(
	user models.User,
) (*models.Partner, error) {
	var partner models.Partner

	if err := copier.CopyWithOption(&partner, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	partner.User = user

	return &partner, nil
}

func (req *UpdatePartner) ToModel(
	partner models.Partner,
	user models.User,
) (*models.Partner, error) {
	if err := copier.CopyWithOption(&partner, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	partner.User = user

	return &partner, nil
}

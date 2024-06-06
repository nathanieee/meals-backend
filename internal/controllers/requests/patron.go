package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreatePatron struct {
		User CreateUser            `json:"user" form:"user" binding:"required,dive"`
		Type consttypes.PatronType `json:"type" form:"type" binding:"required"`
		Name string                `json:"name" form:"name" binding:"required"`
	}

	UpdatePatron struct {
		User UpdateUser            `json:"user" form:"user" binding:"-"`
		Type consttypes.PatronType `json:"type" form:"type" binding:"-"`
		Name string                `json:"name" form:"name" binding:"-"`
	}
)

func (req *CreatePatron) ToModel(
	user models.User,
) (*models.Patron, error) {
	var (
		patron models.Patron
	)

	if err := copier.CopyWithOption(&patron, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	patron.User = user

	return &patron, nil
}

func (req *UpdatePatron) ToModel(
	patron models.Patron,
	user models.User,
) (*models.Patron, error) {
	if err := copier.CopyWithOption(&patron, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	patron.User = user

	return &patron, nil
}

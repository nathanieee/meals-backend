package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreateOrganization struct {
		User CreateUser `json:"user" form:"user" binding:"required,dive"`

		Type consttypes.OrganizationType `json:"type" form:"type" binding:"required"`
		Name string                      `json:"name" form:"name" binding:"required"`
	}

	UpdateOrganization struct {
		User UpdateUser `json:"user" form:"user" binding:"required,dive"`

		Type consttypes.OrganizationType `json:"type" form:"type" binding:"required"`
		Name string                      `json:"name" form:"name" binding:"required"`
	}
)

func (req *CreateOrganization) ToModel(
	user models.User,
) (*models.Organization, error) {
	var (
		organization models.Organization
	)

	if err := copier.CopyWithOption(&organization, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	organization.User = user

	return &organization, nil
}

func (req *UpdateOrganization) ToModel(
	organization models.Organization,
	user models.User,
) (*models.Organization, error) {
	if err := copier.CopyWithOption(&organization, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	organization.User = user

	return &organization, nil
}

func (req *CreateOrganization) ToSignin() *Signin {
	return &Signin{
		Email:    req.User.Email,
		Password: req.User.Password,
	}
}

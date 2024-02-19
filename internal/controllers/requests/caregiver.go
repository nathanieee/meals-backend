package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreateCaregiver struct {
		User        CreateUser        `json:"user" form:"user" binding:"required"`
		Gender      consttypes.Gender `json:"gender" form:"gender" binding:"required"`
		FirstName   string            `json:"first_name" form:"first_name" binding:"required"`
		LastName    string            `json:"last_name" form:"last_name" binding:"required"`
		DateOfBirth customs.CDT_DATE  `json:"date_of_birth" form:"date_of_birth" binding:"required"`
	}

	UpdateCaregiver struct {
		User        UpdateUser        `json:"user" form:"user" binding:"dive"`
		Gender      consttypes.Gender `json:"gender" form:"gender" binding:"-"`
		FirstName   string            `json:"first_name" form:"first_name" binding:"-"`
		LastName    string            `json:"last_name" form:"last_name" binding:"-"`
		DateOfBirth customs.CDT_DATE  `json:"date_of_birth" form:"date_of_birth" binding:"-"`
	}
)

func (req *CreateCaregiver) ToModel() (*models.Caregiver, error) {
	user, err := req.User.ToModel(consttypes.UR_CAREGIVER)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	caregiver := models.Caregiver{
		User: *user,
	}

	if err := copier.CopyWithOption(&caregiver, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &caregiver, nil
}

func (req *UpdateCaregiver) ToCreateCaregiver() (*CreateCaregiver, error) {
	create := CreateCaregiver{}

	if err := copier.CopyWithOption(&create, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &create, nil
}

func (req *UpdateCaregiver) ToModel(
	caregiver *models.Caregiver,
) (*models.Caregiver, error) {
	if caregiver == nil {
		ccg, err := req.ToCreateCaregiver()

		if err != nil {
			utlogger.LogError(err)
			return nil, err
		}

		caregiver, err = ccg.ToModel()
		if err != nil {
			utlogger.LogError(err)
			return nil, err
		}

		return caregiver, nil
	}

	if err := copier.CopyWithOption(&caregiver, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return caregiver, nil
}

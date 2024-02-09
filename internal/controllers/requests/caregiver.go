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
		User        UpdateUser        `json:"user" form:"user" binding:"required"`
		Gender      consttypes.Gender `json:"gender" form:"gender" binding:"required"`
		FirstName   string            `json:"first_name" form:"first_name" binding:"required"`
		LastName    string            `json:"last_name" form:"last_name" binding:"required"`
		DateOfBirth customs.CDT_DATE  `json:"date_of_birth" form:"date_of_birth" binding:"required"`
	}
)

func (req *CreateCaregiver) ToModel() *models.Caregiver {
	user := req.User.ToModel(consttypes.UR_CAREGIVER)
	caregiver := models.Caregiver{
		User: *user,
	}

	if err := copier.Copy(&caregiver, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &caregiver
}

func (req *UpdateCaregiver) ToModel(caregiver models.Caregiver) *models.Caregiver {
	if err := copier.Copy(&caregiver, &req); err != nil {
		utlogger.LogError(err)
		return nil
	}

	return &caregiver
}

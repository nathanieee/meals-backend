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
		User        CreateUser        `json:"user"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		FirstName   string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		DateOfBirth customs.CDT_DATE  `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20" time_format:"2006-01-02"`
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

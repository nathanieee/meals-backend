package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/logger"
	"time"

	"github.com/jinzhu/copier"
)

type (
	CreateCaregiverRequest struct {
		User        CreateUserRequest `json:"user"`
		Gender      consttypes.Gender `json:"gender" gorm:"not null" binding:"required" example:"Male"`
		FirstName   string            `json:"first_name" gorm:"not null" binding:"required" example:"Jonathan"`
		LastName    string            `json:"last_name" gorm:"not null" binding:"required" example:"Vince"`
		DateOfBirth time.Time         `json:"date_of_birth" gorm:"not null" binding:"required" example:"2000-10-20"`
	}
)

func (req *CreateCaregiverRequest) ToModel() *models.Caregiver {
	user := req.User.ToModel(consttypes.UR_CAREGIVER)
	caregiver := models.Caregiver{
		User: *user,
	}

	if err := copier.Copy(&caregiver, &req); err != nil {
		logger.LogError(err)
		return nil
	}

	return &caregiver
}

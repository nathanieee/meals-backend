package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreateIllness struct {
		Name        string `json:"name" example:"Cold Sore"`
		Description string `json:"description" example:"Infection with the herpes simplex virus around the border of the lips."`
	}

	UpdateIllness struct {
		Name        string `json:"name" example:"Cold Sore"`
		Description string `json:"description" example:"Infection with the herpes simplex virus around the border of the lips."`
	}
)

func (req *CreateIllness) ToModel() *models.Illness {
	return &models.Illness{
		Name:        req.Name,
		Description: req.Description,
	}
}

func (req *UpdateIllness) ToModel(ill models.Illness) (*models.Illness, error) {
	if req == nil {
		return &ill, nil
	}

	if err := copier.CopyWithOption(&ill, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &ill, nil
}

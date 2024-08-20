package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	CreateDonation struct {
		*CreateImage

		Value float64 `json:"value" form:"value" binding:"required"`
	}
)

func (req *CreateDonation) ToModel(
	pid uuid.UUID,
) (*models.Donation, error) {
	var (
		donation models.Donation
	)

	if err := copier.CopyWithOption(&donation, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * assign to patron's id
	donation.PatronID = pid

	// * default status is pending
	donation.Status = consttypes.DS_PENDING

	return &donation, nil
}

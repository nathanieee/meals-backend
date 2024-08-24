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

	UpdateDonation struct {
		// ! cannot update image that
		// ! was uploaded by patron
		Value float64 `json:"value" form:"value"`

		Status consttypes.DonationStatus `json:"status" form:"status"`
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

func (req *UpdateDonation) ToModel(don models.Donation) (*models.Donation, error) {
	if req == nil {
		return &don, nil
	}

	if err := copier.CopyWithOption(&don, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &don, nil
}

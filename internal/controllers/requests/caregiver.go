package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/customs/ctdatatype"
	"project-skbackend/packages/utils/utlogger"

	"github.com/jinzhu/copier"
)

type (
	CreateCaregiver struct {
		User CreateUser `json:"user" form:"user" binding:"required,dive"`

		Gender      consttypes.Gender   `json:"gender" form:"gender" binding:"required"`
		FirstName   string              `json:"first_name" form:"first_name" binding:"required"`
		LastName    string              `json:"last_name" form:"last_name" binding:"required"`
		DateOfBirth ctdatatype.CDT_DATE `json:"date_of_birth" form:"date_of_birth" binding:"required"`
	}

	UpdateCaregiver struct {
		User *UpdateUser `json:"user" form:"user" binding:"omitempty,dive"`

		Gender      consttypes.Gender   `json:"gender" form:"gender" binding:"omitempty"`
		FirstName   string              `json:"first_name" form:"first_name" binding:"omitempty"`
		LastName    string              `json:"last_name" form:"last_name" binding:"omitempty"`
		DateOfBirth ctdatatype.CDT_DATE `json:"date_of_birth" form:"date_of_birth" binding:"omitempty"`
	}
)

func (req *CreateCaregiver) ToModel() (*models.Caregiver, error) {
	var (
		caregiver models.Caregiver
	)

	user, err := req.User.ToModel(consttypes.UR_CAREGIVER)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	if err := copier.CopyWithOption(&caregiver, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	caregiver.User = *user

	return &caregiver, nil
}

func (req *CreateCaregiver) FromMemberAddition() (*models.Caregiver, error) {
	var (
		caregiver models.Caregiver
	)

	user, err := req.User.ToModel(consttypes.UR_CAREGIVER)
	user.ConfirmedAt = consttypes.TimeNow()
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	if err := copier.CopyWithOption(&caregiver, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	caregiver.User = *user

	return &caregiver, nil
}

func (req *UpdateCaregiver) ToCreateCaregiver() (*CreateCaregiver, error) {
	var (
		create CreateCaregiver
	)

	if err := copier.CopyWithOption(&create, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &create, nil
}

func (req *UpdateCaregiver) ToModel(
	caregiver *models.Caregiver,
) (*models.Caregiver, error) {
	// * if member does not have any caregiver
	// * then create a new caregiver and return it
	if caregiver == nil {
		createcaregiver, err := req.ToCreateCaregiver()

		if err != nil {
			utlogger.Error(err)
			return nil, err
		}

		caregiver, err = createcaregiver.ToModel()
		if err != nil {
			utlogger.Error(err)
			return nil, err
		}

		return caregiver, nil
	}

	// * if member does have a caregiver
	// * then copy the new data to the old data
	// * and return the new data
	if err := copier.CopyWithOption(&caregiver, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return caregiver, nil
}

package requests

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utgeolocation"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utstring"
	"strconv"
	"strings"

	"github.com/jinzhu/copier"
)

type (
	CreateUser struct {
		*CreateImage

		Email           string `json:"email" form:"email" binding:"required,email"`
		Password        string `json:"password" form:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`

		Addresses []CreateAddress `json:"addresses" form:"address" binding:"-"`
	}

	UpdateUser struct {
		*UpdateImage

		Email           string `json:"email" form:"email" binding:"omitempty"`
		Password        string `json:"password" form:"password" binding:"omitempty"`
		ConfirmPassword string `json:"confirm_password" form:"confirm_password" binding:"omitempty,eqfield=Password"`
	}
)

func (req *CreateUser) ToModel(
	role consttypes.UserRole,
) (*models.User, error) {
	var (
		user models.User
	)

	hash, err := utstring.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if err := copier.CopyWithOption(&user, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	user.Email = strings.ToLower(req.Email)
	user.Role = role
	user.Password = hash

	// * reset the user address
	user.Addresses = nil
	for _, address := range req.Addresses {
		adrdetail, err := utgeolocation.GetGeolocation(*address.ToDistanceMatrixGeolocation())
		if err != nil {
			return nil, consttypes.ErrGeolocationNotFound
		}

		// * format the float into string
		lat := strconv.FormatFloat(adrdetail.Lat, 'f', -1, 64)
		lng := strconv.FormatFloat(adrdetail.Lng, 'f', -1, 64)

		user.Addresses = append(user.Addresses, &models.Address{
			Name:    address.Name,
			Address: address.Address,
			Note:    address.Note,
			UserID:  user.ID,
			AddressDetail: &models.AddressDetail{
				Geolocation: models.Geolocation{
					Latitude:  lat,
					Longitude: lng,
				},
				FormattedAddress: adrdetail.FormattedAddress,
				PostCode:         adrdetail.PostCode,
				Country:          adrdetail.Country},
		})
	}

	return &user, nil
}

func (req *UpdateUser) ToModel(
	user models.User,
	role consttypes.UserRole,
) (*models.User, error) {
	if req == nil {
		return &user, nil
	}

	hash, err := utstring.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	if err := copier.CopyWithOption(&user, &req, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	user.Role = role
	user.Password = hash

	return &user, nil
}

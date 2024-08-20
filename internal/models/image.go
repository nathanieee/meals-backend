package models

import (
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/consttypes"
)

type (
	Image struct {
		base.Model

		Name string               `json:"name" gorm:"required" example:"image.jpg"`
		Path string               `json:"path" gorm:"required" example:"./files/images/profiles/image.jpg"`
		Type consttypes.ImageType `json:"image_type" gorm:"required; type:image_type_enum" example:"Profile"`
	}
)

func NewProfileImage(
	name string,
	path string,
) *Image {
	return &Image{
		Name: name,
		Path: path,
		Type: consttypes.IT_PROFILE,
	}
}

func NewDonationProof(
	name string,
	path string,
) *Image {
	return &Image{
		Name: name,
		Path: path,
		Type: consttypes.IT_DONATION_PROOF,
	}
}

func (i *Image) CreateDonationProof(
	donation Donation,
) *DonationProof {
	var (
		donationproof = DonationProof{}
	)

	donationproof.ImageID = i.ID
	donationproof.DonationID = donation.ID

	return &donationproof
}

func (i *Image) CreateUserImage(
	user User,
) *UserImage {
	var (
		userimage = UserImage{}
	)

	userimage.ImageID = i.ID
	userimage.UserID = user.ID

	return &userimage
}

func (i *Image) UpdateUserImage(
	userimage UserImage,
) *UserImage {
	userimage.ImageID = i.ID

	return &userimage
}

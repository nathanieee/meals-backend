package donationservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/repositories/donationrepo"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	DonationService struct {
		rdonation donationrepo.IDonationRepository
	}

	IDonationService interface {
		Read() ([]*responses.Donation, error)
		Update(req requests.UpdateDonation, donid uuid.UUID) (*responses.Donation, error)
		Delete(donid uuid.UUID) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(donid uuid.UUID) (*responses.Donation, error)
	}
)

func NewDonationService(
	rdonation donationrepo.IDonationRepository,
) *DonationService {
	return &DonationService{
		rdonation: rdonation,
	}
}

func (s *DonationService) Read() ([]*responses.Donation, error) {
	var (
		donreses []*responses.Donation
		err      error
	)

	donations, err := s.rdonation.Read()
	if err != nil {
		return nil, err
	}

	for _, donation := range donations {
		donres, err := donation.ToResponse()
		if err != nil {
			return nil, err
		}

		donreses = append(donreses, donres)
	}

	return donreses, nil
}

func (s *DonationService) Update(req requests.UpdateDonation, donid uuid.UUID) (*responses.Donation, error) {
	var (
		err error
	)

	donation, err := s.rdonation.GetByID(donid)
	if err != nil {
		return nil, err
	}

	donation, err = req.ToModel(*donation)
	if err != nil {
		return nil, err
	}

	donation, err = s.rdonation.Update(*donation)
	if err != nil {
		return nil, err
	}

	donationres, err := donation.ToResponse()
	if err != nil {
		return nil, err
	}

	return donationres, nil
}

func (s *DonationService) Delete(donid uuid.UUID) error {
	donation, err := s.rdonation.GetByID(donid)
	if err != nil {
		return err
	}

	return s.rdonation.Delete(*donation)
}

func (s *DonationService) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	donations, err := s.rdonation.FindAll(p)
	if err != nil {
		return nil, err
	}

	return donations, nil
}

func (s *DonationService) GetByID(donid uuid.UUID) (*responses.Donation, error) {
	donation, err := s.rdonation.GetByID(donid)
	if err != nil {
		return nil, err
	}

	donationres, err := donation.ToResponse()
	if err != nil {
		return nil, err
	}

	return donationres, nil
}

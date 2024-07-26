package partnerservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	PartnerService struct {
		rpart partnerrepo.IPartnerRepository
	}

	IPartnerService interface {
		Create(req requests.CreatePartner) (*responses.Partner, error)
		Read() ([]*models.Partner, error)
		Update(id uuid.UUID, req requests.UpdatePartner) (*responses.Partner, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Partner, error)
	}
)

func NewPartnerService(
	rpart partnerrepo.IPartnerRepository,
) *PartnerService {
	return &PartnerService{
		rpart: rpart,
	}
}

func (s *PartnerService) Create(req requests.CreatePartner) (*responses.Partner, error) {
	user, err := req.User.ToModel(consttypes.UR_PARTNER)
	if err != nil {
		return nil, err
	}

	partner, err := req.ToModel(*user)
	if err != nil {
		return nil, err
	}

	partner, err = s.rpart.Create(*partner)
	if err != nil {
		return nil, err
	}

	pres, err := partner.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

func (s *PartnerService) Read() ([]*models.Partner, error) {
	partners, err := s.rpart.Read()
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (s *PartnerService) Update(id uuid.UUID, req requests.UpdatePartner) (*responses.Partner, error) {
	partner, err := s.rpart.GetByID(id)
	if err != nil {
		return nil, err
	}

	user, err := req.User.ToModel(partner.User, consttypes.UR_PARTNER)
	if err != nil {
		return nil, err
	}

	partner, err = req.ToModel(*partner, *user)
	if err != nil {
		return nil, err
	}

	partner, err = s.rpart.Update(*partner)
	if err != nil {
		return nil, err
	}

	pres, err := partner.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

func (s *PartnerService) Delete(id uuid.UUID) error {
	partner, err := s.rpart.GetByID(id)
	if err != nil {
		return err
	}

	return s.rpart.Delete(*partner)
}

func (s *PartnerService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	partners, err := s.rpart.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (s *PartnerService) GetByID(id uuid.UUID) (*responses.Partner, error) {
	partner, err := s.rpart.GetByID(id)
	if err != nil {
		return nil, err
	}

	pres, err := partner.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

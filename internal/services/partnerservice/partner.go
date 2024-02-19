package partnerservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	PartnerService struct {
		prtrrepo partnerrepo.IPartnerRepository
	}

	IPartnerService interface {
		Create(req requests.CreatePartner) (*responses.Partner, error)
		Read() ([]*models.Partner, error)
		Update(id uuid.UUID, req requests.UpdatePartner) (*responses.Partner, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*responses.Partner, error)
	}
)

func NewPartnerService(
	prtrrepo partnerrepo.IPartnerRepository,
) *PartnerService {
	return &PartnerService{
		prtrrepo: prtrrepo,
	}
}

func (s *PartnerService) Create(req requests.CreatePartner) (*responses.Partner, error) {
	user := req.User.ToModel(consttypes.UR_PARTNER)

	partner := req.ToModel(*user)
	partner, err := s.prtrrepo.Create(*partner)
	if err != nil {
		return nil, err
	}

	pres := partner.ToResponse()

	return pres, nil
}

func (s *PartnerService) Read() ([]*models.Partner, error) {
	partners, err := s.prtrrepo.Read()
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (s *PartnerService) Update(id uuid.UUID, req requests.UpdatePartner) (*responses.Partner, error) {
	partner, err := s.prtrrepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	user := req.User.ToModel(partner.User, consttypes.UR_PARTNER)

	partner = req.ToModel(*partner, *user)
	partner, err = s.prtrrepo.Update(*partner)
	if err != nil {
		return nil, err
	}

	pres := partner.ToResponse()

	return pres, nil
}

func (s *PartnerService) Delete(id uuid.UUID) error {
	partner := models.Partner{
		Model: helper.Model{ID: id},
	}

	err := s.prtrrepo.Delete(partner)
	if err != nil {
		return err
	}

	return nil
}

func (s *PartnerService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	partners, err := s.prtrrepo.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return partners, nil
}

func (s *PartnerService) FindByID(id uuid.UUID) (*responses.Partner, error) {
	partner, err := s.prtrrepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	pres := partner.ToResponse()

	return pres, nil
}

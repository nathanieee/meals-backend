package organizationservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	OrganizationService struct {
		rorg organizationrepo.IOrganizationRepository
	}

	IOrganizationService interface {
		Create(req requests.CreateOrganization) (*responses.Organization, error)
		Read() ([]*models.Organization, error)
		Update(id uuid.UUID, req requests.UpdateOrganization) (*responses.Organization, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*responses.Organization, error)
	}
)

func NewOrganizationService(
	rorg organizationrepo.IOrganizationRepository,
) *OrganizationService {
	return &OrganizationService{
		rorg: rorg,
	}
}

func (s *OrganizationService) Create(req requests.CreateOrganization) (*responses.Organization, error) {
	user, err := req.User.ToModel(consttypes.UR_ORGANIZATION)
	if err != nil {
		return nil, err
	}

	org, err := req.ToModel(*user)
	if err != nil {
		return nil, err
	}

	org, err = s.rorg.Create(*org)
	if err != nil {
		return nil, err
	}

	orgres, err := org.ToResponse()
	if err != nil {
		return nil, err
	}

	return orgres, nil
}

func (s *OrganizationService) Read() ([]*models.Organization, error) {
	organizations, err := s.rorg.Read()
	if err != nil {
		return nil, err
	}

	return organizations, nil
}

func (s *OrganizationService) Update(id uuid.UUID, req requests.UpdateOrganization) (*responses.Organization, error) {
	org, err := s.rorg.FindByID(id)
	if err != nil {
		return nil, err
	}

	user, err := req.User.ToModel(org.User, consttypes.UR_ORGANIZATION)
	if err != nil {
		return nil, err
	}

	org, err = req.ToModel(*org, *user)
	if err != nil {
		return nil, err
	}

	org, err = s.rorg.Update(*org)
	if err != nil {
		return nil, err
	}

	orgres, err := org.ToResponse()
	if err != nil {
		return nil, err
	}

	return orgres, nil
}

func (s *OrganizationService) Delete(id uuid.UUID) error {
	org, err := s.rorg.FindByID(id)
	if err != nil {
		return err
	}

	return s.rorg.Delete(*org)
}

func (s *OrganizationService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	orgs, err := s.rorg.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return orgs, nil
}

func (s *OrganizationService) FindByID(id uuid.UUID) (*responses.Organization, error) {
	org, err := s.rorg.FindByID(id)
	if err != nil {
		return nil, err
	}

	orgres, err := org.ToResponse()
	if err != nil {
		return nil, err
	}

	return orgres, nil
}

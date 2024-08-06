package organizationservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
)

type (
	OrganizationService struct {
		rorg organizationrepo.IOrganizationRepository
	}

	IOrganizationService interface {
		Create(req requests.CreateOrganization) (*responses.Organization, error)
		Read() ([]*responses.Organization, error)
		Update(id uuid.UUID, req requests.UpdateOrganization) (*responses.Organization, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Organization, error)
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
	// TODO: upload the image to S3 bucket and get the image url

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

func (s *OrganizationService) Read() ([]*responses.Organization, error) {
	var (
		orgreses []*responses.Organization
	)

	org, err := s.rorg.Read()
	if err != nil {
		return nil, err
	}

	if err := copier.CopyWithOption(&orgreses, &org, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return orgreses, nil
}

func (s *OrganizationService) Update(id uuid.UUID, req requests.UpdateOrganization) (*responses.Organization, error) {
	org, err := s.rorg.GetByID(id)
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
	org, err := s.rorg.GetByID(id)
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

func (s *OrganizationService) GetByID(id uuid.UUID) (*responses.Organization, error) {
	org, err := s.rorg.GetByID(id)
	if err != nil {
		return nil, err
	}

	orgres, err := org.ToResponse()
	if err != nil {
		return nil, err
	}

	return orgres, nil
}

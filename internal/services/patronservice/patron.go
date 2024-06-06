package patronservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/patronrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	PatronService struct {
		rpatr patronrepo.IPatronRepository
	}

	IPatronService interface {
		Create(req requests.CreatePatron) (*responses.Patron, error)
		Read() ([]*models.Patron, error)
		Update(id uuid.UUID, req requests.UpdatePatron) (*responses.Patron, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*responses.Patron, error)
	}
)

func NewPatronService(
	rpatr patronrepo.IPatronRepository,
) *PatronService {
	return &PatronService{
		rpatr: rpatr,
	}
}

func (s *PatronService) Create(req requests.CreatePatron) (*responses.Patron, error) {
	user, err := req.User.ToModel(consttypes.UR_PATRON)
	if err != nil {
		return nil, err
	}

	patron, err := req.ToModel(*user)
	if err != nil {
		return nil, err
	}

	patron, err = s.rpatr.Create(*patron)
	if err != nil {
		return nil, err
	}

	pres, err := patron.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

func (s *PatronService) Read() ([]*models.Patron, error) {
	patrons, err := s.rpatr.Read()
	if err != nil {
		return nil, err
	}

	return patrons, nil
}

func (s *PatronService) Update(id uuid.UUID, req requests.UpdatePatron) (*responses.Patron, error) {
	patron, err := s.rpatr.FindByID(id)
	if err != nil {
		return nil, err
	}

	user, err := req.User.ToModel(patron.User, consttypes.UR_PATRON)
	if err != nil {
		return nil, err
	}

	patron, err = req.ToModel(*patron, *user)
	if err != nil {
		return nil, err
	}

	patron, err = s.rpatr.Update(*patron)
	if err != nil {
		return nil, err
	}

	pres, err := patron.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

func (s *PatronService) Delete(id uuid.UUID) error {
	patron, err := s.rpatr.FindByID(id)
	if err != nil {
		return err
	}

	err = s.rpatr.Delete(*patron)
	if err != nil {
		return err
	}

	return nil
}

func (s *PatronService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	patrons, err := s.rpatr.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return patrons, nil
}

func (s *PatronService) FindByID(id uuid.UUID) (*responses.Patron, error) {
	patron, err := s.rpatr.FindByID(id)
	if err != nil {
		return nil, err
	}

	pres, err := patron.ToResponse()
	if err != nil {
		return nil, err
	}

	return pres, nil
}

package allergyservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	AllergyService struct {
		ralle allergyrepo.IAllergyRepository
	}

	IAllergyService interface {
		Create(req requests.CreateAllergy) (*responses.Allergy, error)
		Read() ([]*responses.Allergy, error)
		Update(req requests.UpdateAllergy, alid uuid.UUID) (*responses.Allergy, error)
		Delete(alid uuid.UUID) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(alid uuid.UUID) (*responses.Allergy, error)
	}
)

func NewAllergyService(
	ralle allergyrepo.IAllergyRepository,
) *AllergyService {
	return &AllergyService{
		ralle: ralle,
	}
}

func (s *AllergyService) Create(req requests.CreateAllergy) (*responses.Allergy, error) {
	var (
		err     error
		allergy = req.ToModel()
	)

	allergy, err = s.ralle.Create(*allergy)
	if err != nil {
		return nil, err
	}

	allergyres, err := allergy.ToResponse()
	if err != nil {
		return nil, err
	}

	return allergyres, nil
}

func (s *AllergyService) Read() ([]*responses.Allergy, error) {
	var (
		allreses []*responses.Allergy
		err      error
	)

	allergies, err := s.ralle.Read()
	if err != nil {
		return nil, err
	}

	for _, allergy := range allergies {
		allergyres, err := allergy.ToResponse()
		if err != nil {
			return nil, err
		}

		allreses = append(allreses, allergyres)
	}

	return allreses, nil
}

func (s *AllergyService) Update(req requests.UpdateAllergy, alid uuid.UUID) (*responses.Allergy, error) {
	var (
		err error
	)

	allergy, err := s.ralle.GetByID(alid)
	if err != nil {
		return nil, err
	}

	allergy, err = req.ToModel(*allergy)
	if err != nil {
		return nil, err
	}

	allergy, err = s.ralle.Update(*allergy)
	if err != nil {
		return nil, err
	}

	allergyres, err := allergy.ToResponse()
	if err != nil {
		return nil, err
	}

	return allergyres, nil
}

func (s *AllergyService) Delete(alid uuid.UUID) error {
	allergy, err := s.ralle.GetByID(alid)
	if err != nil {
		return err
	}

	return s.ralle.Delete(*allergy)
}

func (s *AllergyService) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	allergies, err := s.ralle.FindAll(p)
	if err != nil {
		return nil, err
	}

	return allergies, nil
}

func (s *AllergyService) GetByID(alid uuid.UUID) (*responses.Allergy, error) {
	allergy, err := s.ralle.GetByID(alid)
	if err != nil {
		return nil, err
	}

	all, err := allergy.ToResponse()
	if err != nil {
		return nil, err
	}

	return all, nil
}

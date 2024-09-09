package caregiverservice

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/repositories/caregiverrepo"

	"github.com/google/uuid"
)

type (
	CaregiverService struct {
		// * repository
		rcare caregiverrepo.ICaregiverRepository
	}

	ICaregiverService interface {
		GetByID(id uuid.UUID) (*responses.Caregiver, error)
	}
)

func NewCaregiverService(
	// * repository
	rcare caregiverrepo.ICaregiverRepository,
) *CaregiverService {
	return &CaregiverService{
		// * repository
		rcare: rcare,
	}
}

func (s *CaregiverService) GetByID(id uuid.UUID) (*responses.Caregiver, error) {
	var (
		careres *responses.Caregiver
		err     error
	)

	caremod, err := s.rcare.GetByID(id)
	if err != nil {
		return nil, err
	}

	careres, err = caremod.ToResponse()
	if err != nil {
		return nil, err
	}

	return careres, nil
}

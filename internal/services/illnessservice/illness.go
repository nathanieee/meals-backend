package illnessservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	IllnessService struct {
		rill illnessrepo.IIllnessRepository
	}

	IIllnessService interface {
		Create(req requests.CreateIllness) (*responses.Illness, error)
		Read() ([]*responses.Illness, error)
		Update(id uuid.UUID, req requests.UpdateIllness) (*responses.Illness, error)
		Delete(id uuid.UUID) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Illness, error)
	}
)

func NewIllnessService(
	rill illnessrepo.IIllnessRepository,
) *IllnessService {
	return &IllnessService{
		rill: rill,
	}
}

func (s *IllnessService) Create(req requests.CreateIllness) (*responses.Illness, error) {
	var (
		err     error
		illness = req.ToModel()
	)

	illness, err = s.rill.Create(*illness)
	if err != nil {
		return nil, err
	}

	illress, err := illness.ToResponse()
	if err != nil {
		return nil, err
	}

	return illress, nil
}

func (s *IllnessService) Read() ([]*responses.Illness, error) {
	var (
		illreses []*responses.Illness
		err      error
	)

	illnesses, err := s.rill.Read()
	if err != nil {
		return nil, err
	}

	for _, illness := range illnesses {
		ill, err := illness.ToResponse()
		if err != nil {
			return nil, err
		}

		illreses = append(illreses, ill)
	}

	return illreses, nil
}

func (s *IllnessService) Update(id uuid.UUID, req requests.UpdateIllness) (*responses.Illness, error) {
	illness, err := s.rill.GetByID(id)
	if err != nil {
		return nil, err
	}

	illness, err = req.ToModel(*illness)
	if err != nil {
		return nil, err
	}

	illness, err = s.rill.Update(*illness)
	if err != nil {
		return nil, err
	}

	illress, err := illness.ToResponse()
	if err != nil {
		return nil, err
	}

	return illress, nil
}

func (s *IllnessService) Delete(id uuid.UUID) error {
	illness, err := s.rill.GetByID(id)
	if err != nil {
		return err
	}

	return s.rill.Delete(*illness)
}

func (s *IllnessService) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	illnesses, err := s.rill.FindAll(p)
	if err != nil {
		return nil, err
	}

	return illnesses, nil
}

func (s *IllnessService) GetByID(id uuid.UUID) (*responses.Illness, error) {
	illness, err := s.rill.GetByID(id)
	if err != nil {
		return nil, err
	}

	ill, err := illness.ToResponse()
	if err != nil {
		return nil, err
	}

	return ill, nil
}

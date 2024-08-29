package mealcategoryservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/repositories/mealcategoryrepo"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	MealCategoryService struct {
		rmcat mealcategoryrepo.IMealCategoryRepository
	}

	IMealCategoryService interface {
		Create(req requests.CreateMealCategory) (*responses.MealCategory, error)
		Read() ([]*responses.MealCategory, error)
		Update(id uuid.UUID, req requests.UpdateMealCategory) (*responses.MealCategory, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.MealCategory, error)
	}
)

func NewMealCategoryService(
	rmcat mealcategoryrepo.IMealCategoryRepository,
) *MealCategoryService {
	return &MealCategoryService{
		rmcat: rmcat,
	}
}

func (s *MealCategoryService) Create(req requests.CreateMealCategory) (*responses.MealCategory, error) {
	var (
		err error
	)

	mcat, err := req.ToModel()
	if err != nil {
		return nil, err
	}

	mcat, err = s.rmcat.Create(*mcat)
	if err != nil {
		return nil, err
	}

	mcatres, err := mcat.ToResponse()
	if err != nil {
		return nil, err
	}

	return mcatres, nil
}

func (s *MealCategoryService) Read() ([]*responses.MealCategory, error) {
	var (
		mcatreses []*responses.MealCategory
	)

	mcats, err := s.rmcat.Read()
	if err != nil {
		return nil, err
	}

	for _, mcat := range mcats {
		mcatres, err := mcat.ToResponse()
		if err != nil {
			return nil, err
		}

		mcatreses = append(mcatreses, mcatres)
	}

	return mcatreses, nil
}

func (s *MealCategoryService) Update(id uuid.UUID, req requests.UpdateMealCategory) (*responses.MealCategory, error) {
	mcat, err := s.rmcat.GetByID(id)
	if err != nil {
		return nil, err
	}

	mcat, err = req.ToModel(mcat)
	if err != nil {
		return nil, err
	}

	mcat, err = s.rmcat.Update(*mcat)
	if err != nil {
		return nil, err
	}

	mcatres, err := mcat.ToResponse()
	if err != nil {
		return nil, err
	}

	return mcatres, nil
}

func (s *MealCategoryService) Delete(id uuid.UUID) error {
	mcat, err := s.rmcat.GetByID(id)
	if err != nil {
		return err
	}

	return s.rmcat.Delete(*mcat)
}

func (s *MealCategoryService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	return s.rmcat.FindAll(preq)
}

func (s *MealCategoryService) GetByID(id uuid.UUID) (*responses.MealCategory, error) {
	mcat, err := s.rmcat.GetByID(id)
	if err != nil {
		return nil, err
	}

	mcatres, err := mcat.ToResponse()
	if err != nil {
		return nil, err
	}

	return mcatres, nil
}

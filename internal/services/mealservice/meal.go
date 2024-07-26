package mealservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	MealService struct {
		rmeal mealrepo.IMealRepository
		rill  illnessrepo.IIllnessRepository
		rall  allergyrepo.IAllergyRepository
		rpart partnerrepo.IPartnerRepository
	}

	IMealService interface {
		Create(req requests.CreateMeal) (*responses.Meal, error)
		Read() ([]*models.Meal, error)
		Update(id uuid.UUID, req requests.UpdateMeal) (*responses.Meal, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*responses.Meal, error)
	}
)

func NewMealService(
	rmeal mealrepo.IMealRepository,
	rill illnessrepo.IIllnessRepository,
	rall allergyrepo.IAllergyRepository,
	rpart partnerrepo.IPartnerRepository,

) *MealService {
	return &MealService{
		rmeal: rmeal,
		rill:  rill,
		rall:  rall,
		rpart: rpart,
	}
}

func (s *MealService) Create(req requests.CreateMeal) (*responses.Meal, error) {
	var (
		illnesses []*models.MealIllness
		allergies []*models.MealAllergy
		images    []*models.MealImage
		partner   *models.Partner
	)

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := s.rill.GetByID(*ill)
		if err != nil {
			return nil, consttypes.ErrIllnessNotFound
		}

		millness := illness.ToMealIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := s.rall.GetByID(*all)
		if err != nil {
			return nil, consttypes.ErrAllergiesNotFound
		}

		mallergy := allergy.ToMealAllergy()

		allergies = append(allergies, mallergy)
	}

	partner, err := s.rpart.GetByID(req.PartnerID)
	if err != nil {
		return nil, consttypes.ErrPartnerNotFound
	}

	meal, err := req.ToModel(images, illnesses, allergies, *partner)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	meal, err = s.rmeal.Create(*meal)
	if err != nil {
		return nil, consttypes.ErrFailedToCreateMeal
	}

	meres, err := meal.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return meres, nil
}

func (s *MealService) Read() ([]*models.Meal, error) {
	meals, err := s.rmeal.Read()
	if err != nil {
		return nil, consttypes.ErrFailedToReadMeals
	}

	return meals, nil
}

func (s *MealService) Update(id uuid.UUID, req requests.UpdateMeal) (*responses.Meal, error) {
	var (
		images    []*models.MealImage
		illnesses []*models.MealIllness
		allergies []*models.MealAllergy
		partner   *models.Partner
	)

	meal, err := s.rmeal.GetByID(id)
	if err != nil {
		return nil, consttypes.ErrMealsNotFound
	}

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := s.rill.GetByID(*ill)
		if err != nil {
			return nil, consttypes.ErrIllnessNotFound
		}

		millness := illness.ToMealIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := s.rall.GetByID(*all)
		if err != nil {
			return nil, consttypes.ErrAllergiesNotFound
		}

		mallergy := allergy.ToMealAllergy()

		allergies = append(allergies, mallergy)
	}

	partner, err = s.rpart.GetByID(req.PartnerID)
	if err != nil {
		return nil, consttypes.ErrPartnerNotFound
	}

	// TODO: handle the uploading of meal image
	meal, err = req.ToModel(*meal, images, illnesses, allergies, *partner)
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	meal, err = s.rmeal.Update(*meal)
	if err != nil {
		return nil, consttypes.ErrFailedToUpdateMeal
	}

	mres, err := meal.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return mres, nil
}

func (s *MealService) Delete(id uuid.UUID) error {
	meal := models.Meal{
		Model: base.Model{ID: id},
	}

	err := s.rmeal.Delete(meal)
	if err != nil {
		return consttypes.ErrFailedToDeleteMeal
	}

	return nil
}

func (s *MealService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	meals, err := s.rmeal.FindAll(preq)
	if err != nil {
		return nil, consttypes.ErrFailedToFindAllMeals
	}

	return meals, nil
}

func (s *MealService) GetByID(id uuid.UUID) (*responses.Meal, error) {
	meal, err := s.rmeal.GetByID(id)
	if err != nil {
		return nil, err
	}

	mres, err := meal.ToResponse()
	if err != nil {
		return nil, consttypes.ErrConvertFailed
	}

	return mres, nil
}

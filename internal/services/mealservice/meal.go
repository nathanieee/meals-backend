package mealservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	MealService struct {
		mealrepo mealrepo.IMealRepository
		illrepo  illnessrepo.IIllnessRepository
		allgrepo allergyrepo.IAllergyRepository
	}

	IMealService interface {
		Create(req requests.CreateMeal) (*responses.Meal, error)
		Read() ([]*models.Meal, error)
		Update(id uuid.UUID, req requests.UpdateMeal) (*responses.Meal, error)
		Delete(id uuid.UUID) error
		FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*responses.Meal, error)
	}
)

func NewMealService(
	mealrepo mealrepo.IMealRepository,
	illrepo illnessrepo.IIllnessRepository,
	allgrepo allergyrepo.IAllergyRepository,

) *MealService {
	return &MealService{
		mealrepo: mealrepo,
		illrepo:  illrepo,
		allgrepo: allgrepo,
	}
}

func (s *MealService) Create(req requests.CreateMeal) (*responses.Meal, error) {
	var illnesses []*models.MealIllness
	var allergies []*models.MealAllergy
	var images []*models.MealImage
	var partner models.Partner

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := s.illrepo.FindByID(*ill)
		if err != nil {
			return nil, err
		}

		millness := illness.ToMealIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := s.allgrepo.FindByID(*all)
		if err != nil {
			return nil, err
		}

		mallergy := allergy.ToMealAllergy()

		allergies = append(allergies, mallergy)
	}

	// TODO - create a partner repository to CRUD new partner
	// partner, err :=

	meal := req.ToModel(images, illnesses, allergies, partner)
	meal, err := s.mealrepo.Create(*meal)
	if err != nil {
		return nil, err
	}

	meres := meal.ToResponse()

	return meres, nil
}

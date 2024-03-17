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
		FindByID(id uuid.UUID) (*responses.Meal, error)
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
	var illnesses []*models.MealIllness
	var allergies []*models.MealAllergy
	var images []*models.MealImage
	var partner *models.Partner

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := s.rill.FindByID(*ill)
		if err != nil {
			return nil, err
		}

		millness := illness.ToMealIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := s.rall.FindByID(*all)
		if err != nil {
			return nil, err
		}

		mallergy := allergy.ToMealAllergy()

		allergies = append(allergies, mallergy)
	}

	partner, err := s.rpart.FindByID(req.PartnerID)
	if err != nil {
		return nil, err
	}

	meal, err := req.ToModel(images, illnesses, allergies, *partner)
	if err != nil {
		return nil, err
	}

	meal, err = s.rmeal.Create(*meal)
	if err != nil {
		return nil, err
	}

	meres := meal.ToResponse()

	return meres, nil
}

func (s *MealService) Read() ([]*models.Meal, error) {
	meals, err := s.rmeal.Read()
	if err != nil {
		return nil, err
	}

	return meals, nil
}

func (s *MealService) Update(id uuid.UUID, req requests.UpdateMeal) (*responses.Meal, error) {
	var images []*models.MealImage
	var illnesses []*models.MealIllness
	var allergies []*models.MealAllergy
	var partner *models.Partner

	meal, err := s.rmeal.FindByID(id)
	if err != nil {
		return nil, err
	}

	// * find illness object and append to the array.
	for _, ill := range req.IllnessID {
		illness, err := s.rill.FindByID(*ill)
		if err != nil {
			return nil, err
		}

		millness := illness.ToMealIllness()

		illnesses = append(illnesses, millness)
	}

	// * find allergy object and append to the array.
	for _, all := range req.AllergyID {
		allergy, err := s.rall.FindByID(*all)
		if err != nil {
			return nil, err
		}

		mallergy := allergy.ToMealAllergy()

		allergies = append(allergies, mallergy)
	}

	partner, err = s.rpart.FindByID(req.PartnerID)
	if err != nil {
		return nil, err
	}

	meal, err = req.ToModel(*meal, images, illnesses, allergies, *partner)
	if err != nil {
		return nil, err
	}

	meal, err = s.rmeal.Update(*meal)
	if err != nil {
		return nil, err
	}

	mres := meal.ToResponse()

	return mres, nil
}

func (s *MealService) Delete(id uuid.UUID) error {
	meal := models.Meal{
		Model: base.Model{ID: id},
	}

	err := s.rmeal.Delete(meal)
	if err != nil {
		return err
	}

	return nil
}

func (s *MealService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	meals, err := s.rmeal.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return meals, nil
}

func (s *MealService) FindByID(id uuid.UUID) (*responses.Meal, error) {
	meal, err := s.rmeal.FindByID(id)
	if err != nil {
		return nil, err
	}

	meres := meal.ToResponse()

	return meres, nil
}

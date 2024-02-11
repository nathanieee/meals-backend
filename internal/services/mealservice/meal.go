package mealservice

import (
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	MealService struct {
		mealrepo mealrepo.IMealRepository
		illrepo  illnessrepo.IIllnessRepository
		allgrepo allergyrepo.IAllergyRepository
		prtrrepo partnerrepo.IPartnerRepository
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
	prtrrepo partnerrepo.IPartnerRepository,

) *MealService {
	return &MealService{
		mealrepo: mealrepo,
		illrepo:  illrepo,
		allgrepo: allgrepo,
		prtrrepo: prtrrepo,
	}
}

func (s *MealService) Create(req requests.CreateMeal) (*responses.Meal, error) {
	var illnesses []*models.MealIllness
	var allergies []*models.MealAllergy
	var images []*models.MealImage
	var partner *models.Partner

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

	partner, err := s.prtrrepo.FindByID(req.PartnerID)
	if err != nil {
		return nil, err
	}

	meal := req.ToModel(images, illnesses, allergies, *partner)
	meal, err = s.mealrepo.Create(*meal)
	if err != nil {
		return nil, err
	}

	meres := meal.ToResponse()

	return meres, nil
}

func (s *MealService) Read() ([]*models.Meal, error) {
	meals, err := s.mealrepo.Read()
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

	meal, err := s.mealrepo.FindByID(id)
	if err != nil {
		return nil, err
	}

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

	partner, err = s.prtrrepo.FindByID(req.PartnerID)
	if err != nil {
		return nil, err
	}

	meal = req.ToModel(*meal, images, illnesses, allergies, *partner)
	meal, err = s.mealrepo.Update(*meal)
	if err != nil {
		return nil, err
	}

	mres := meal.ToResponse()

	return mres, nil
}

func (s *MealService) Delete(id uuid.UUID) error {
	meal := models.Meal{
		Model: helper.Model{ID: id},
	}

	err := s.mealrepo.Delete(meal)
	if err != nil {
		return err
	}

	return nil
}

func (s *MealService) FindAll(preq utpagination.Pagination) (*utpagination.Pagination, error) {
	meals, err := s.mealrepo.FindAll(preq)
	if err != nil {
		return nil, err
	}

	return meals, nil
}

func (s *MealService) FindByID(id uuid.UUID) (*responses.Meal, error) {
	meal, err := s.mealrepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	meres := meal.ToResponse()

	return meres, nil
}

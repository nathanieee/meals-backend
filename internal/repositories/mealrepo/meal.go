package mealrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/internal/repositories/paginationrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		partner_id,
		name,
		status,
		description,
		created_at,
		updated_at
	`
)

type (
	MealRepository struct {
		db *gorm.DB
	}

	IMealRepository interface {
		Create(m models.Meal) (*models.Meal, error)
		Read() ([]*models.Meal, error)
		Update(m models.Meal) (*models.Meal, error)
		Delete(m models.Meal) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(id uuid.UUID) (*models.Meal, error)
	}
)

func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{db: db}
}

func (r *MealRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("MealImage.Image").
		Preload("MealIllnesses.Illness").
		Preload("MealAllergies.Allergy").
		Preload("Partner.User.Address").
		Preload("Partner.User.Image.Image")
}

func (r *MealRepository) Create(m models.Meal) (*models.Meal, error) {
	err := r.db.
		Create(m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &m, err
}

func (r *MealRepository) Read() ([]*models.Meal, error) {
	var m []*models.Meal

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return m, nil
}

func (r *MealRepository) Update(m models.Meal) (*models.Meal, error) {
	err := r.db.
		Save(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &m, nil
}

func (r *MealRepository) Delete(m models.Meal) error {
	err := r.db.
		Delete(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *MealRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var m []models.Meal
	var mlres []responses.Meal

	result := r.
		preload().
		Model(&m).
		Select(SELECTED_FIELDS)

	p.Search = fmt.Sprintf("%%%s%%", p.Search)
	if p.Search != "" {
		result = result.
			Where(r.db.
				Where(&models.Meal{Name: p.Search}).
				Or(&models.Meal{Description: p.Search}),
			)
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) between ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	if p.Filter.Partner.ID != nil {
		result = result.
			Where(&models.Meal{PartnerID: *p.Filter.Partner.ID})
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&m, &p, result)).
		Find(&m)

	if err := result.Error; err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&mlres, &m, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = mlres
	return &p, result.Error
}

func (r *MealRepository) FindByID(id uuid.UUID) (*models.Meal, error) {
	var m *models.Meal

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Meal{Model: helper.Model{ID: id}}).
		First(&m).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return m, nil
}

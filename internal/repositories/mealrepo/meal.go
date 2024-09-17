package mealrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
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
		GetByID(id uuid.UUID) (*models.Meal, error)
		ReadByPartnerID(pid uuid.UUID) ([]models.Meal, error)
	}
)

func NewMealRepository(db *gorm.DB) *MealRepository {
	return &MealRepository{db: db}
}

func (r *MealRepository) omit() *gorm.DB {
	return r.db.Omit(
		"Illnesses.Illness",
		"Allergies.Allergy",
		"Partner",
	)
}

func (r *MealRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Images.Image").
		Preload("Illnesses.Illness").
		Preload("Allergies.Allergy").
		Preload("Partner.User.Addresses.AddressDetail").
		Preload("Partner.User.Image.Image")
}

func (r *MealRepository) Create(m models.Meal) (*models.Meal, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	mnew, err := r.GetByID(m.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mnew, nil
}

func (r *MealRepository) Read() ([]*models.Meal, error) {
	var (
		m []*models.Meal
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}

func (r *MealRepository) Update(m models.Meal) (*models.Meal, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	mnew, err := r.GetByID(m.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mnew, nil
}

func (r *MealRepository) Delete(m models.Meal) error {
	err := r.db.
		Delete(&m).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *MealRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		m     []models.Meal
		mlres []responses.Meal
	)

	result := r.
		preload().
		Debug().
		Model(&m).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(
				r.db.Where(`
					name ILIKE ?
						OR 
					description ILIKE ? 
			`, p.Search, p.Search),
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
			Where("partner_id = ?", p.Filter.Partner.ID)
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&m, &p, result)).
		Find(&m)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&mlres, &m, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = mlres
	return &p, nil
}

func (r *MealRepository) GetByID(id uuid.UUID) (*models.Meal, error) {
	var (
		m *models.Meal
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Meal{Model: base.Model{ID: id}}).
		First(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}

func (r *MealRepository) ReadByPartnerID(pid uuid.UUID) ([]models.Meal, error) {
	var (
		m []models.Meal
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Meal{PartnerID: pid}).
		Find(&m).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return m, nil
}

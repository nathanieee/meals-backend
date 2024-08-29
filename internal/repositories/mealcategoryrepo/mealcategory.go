package mealcategoryrepo

import (
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
		name,
		image_id,
		created_at,
		updated_at
	`
)

type (
	MealCategoryRepository struct {
		db *gorm.DB
	}

	IMealCategoryRepository interface {
		Create(m models.MealCategory) (*models.MealCategory, error)
		Read() ([]*models.MealCategory, error)
		Update(m models.MealCategory) (*models.MealCategory, error)
		Delete(m models.MealCategory) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*models.MealCategory, error)
	}
)

func NewMealCategoryRepository(db *gorm.DB) *MealCategoryRepository {
	return &MealCategoryRepository{db: db}
}

func (r *MealCategoryRepository) omit() *gorm.DB {
	return r.db.Omit(
		"",
	)
}

func (r *MealCategoryRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Image")
}

func (r *MealCategoryRepository) Create(mc models.MealCategory) (*models.MealCategory, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&mc).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	mcnew, err := r.GetByID(mc.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mcnew, nil
}

func (r *MealCategoryRepository) Read() ([]*models.MealCategory, error) {
	var (
		mcs []*models.MealCategory
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&mcs).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mcs, nil
}

func (r *MealCategoryRepository) Update(mc models.MealCategory) (*models.MealCategory, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&mc).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	mcnew, err := r.GetByID(mc.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mcnew, nil
}

func (r *MealCategoryRepository) Delete(mc models.MealCategory) error {
	err := r.db.
		Delete(&mc).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *MealCategoryRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		mcs    []models.MealCategory
		mcsres []responses.MealCategory
	)

	result := r.
		preload().
		Model(&mcs).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		result = result.Where("name LIKE ?", "%"+p.Search+"%")
		// TODO: add a like query here
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) between ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&mcs, &p, result)).
		Find(&mcs)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&mcsres, &mcs, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = mcsres
	return &p, nil
}

func (r *MealCategoryRepository) GetByID(id uuid.UUID) (*models.MealCategory, error) {
	var (
		mc *models.MealCategory
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.MealCategory{Model: base.Model{ID: id}}).
		First(&mc).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return mc, nil
}

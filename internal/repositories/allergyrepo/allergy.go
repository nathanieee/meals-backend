package allergyrepo

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
		name,
		description,
		allergens,
		created_at,
		updated_at
	`
)

type (
	AllergyRepository struct {
		db *gorm.DB
	}

	IAllergyRepository interface {
		Create(al models.Allergy) (*models.Allergy, error)
		Read() ([]*models.Allergy, error)
		Update(al models.Allergy) (*models.Allergy, error)
		Delete(al models.Allergy) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		GetByID(id uuid.UUID) (*models.Allergy, error)
	}
)

func NewAllergyRepository(db *gorm.DB) *AllergyRepository {
	return &AllergyRepository{db: db}
}

func (r *AllergyRepository) omit() *gorm.DB {
	return r.db.Omit(
		"",
	)
}

func (r *AllergyRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *AllergyRepository) Create(al models.Allergy) (*models.Allergy, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&al).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	alnew, err := r.GetByID(al.ID)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return alnew, nil
}

func (r *AllergyRepository) Read() ([]*models.Allergy, error) {
	var (
		al []*models.Allergy
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&al).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return al, nil
}

func (r *AllergyRepository) Update(al models.Allergy) (*models.Allergy, error) {
	err := r.db.
		Save(&al).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	alnew, err := r.GetByID(al.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return alnew, nil
}

func (r *AllergyRepository) Delete(al models.Allergy) error {
	err := r.db.
		Delete(&al).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *AllergyRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		al    []models.Allergy
		alres []responses.Admin
	)

	result := r.
		preload().
		Model(&al).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where("name LIKE ?", p.Search).
				Or("description LIKE ?", p.Search),
			)
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) BETWEEN ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT),
			)
	}

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&al, &p, result)).
		Find(&al)

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&alres, &al, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = alres
	return &p, nil
}

func (r *AllergyRepository) GetByID(id uuid.UUID) (*models.Allergy, error) {
	var (
		al *models.Allergy
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Allergy{Model: base.Model{ID: id}}).
		First(&al).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return al, nil
}

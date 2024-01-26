package illnessrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
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
		created_at,
		updated_at
	`
)

type (
	IllnessRepository struct {
		db *gorm.DB
	}

	IIllnessRepository interface {
		Create(ill models.Illness) (*models.Illness, error)
		Read() ([]*models.Illness, error)
		Update(ill models.Illness) (*models.Illness, error)
		Delete(ill models.Illness) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(illid uuid.UUID) (*models.Illness, error)
	}
)

func NewIllnessRepository(db *gorm.DB) *IllnessRepository {
	return &IllnessRepository{db: db}
}

func (r *IllnessRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *IllnessRepository) Create(ill models.Illness) (*models.Illness, error) {
	err := r.db.
		Create(ill).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &ill, nil
}

func (r *IllnessRepository) Read() ([]*models.Illness, error) {
	var ill []*models.Illness

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&ill).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return ill, nil
}

func (r *IllnessRepository) Update(ill models.Illness) (*models.Illness, error) {
	err := r.db.
		Save(&ill).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &ill, nil
}

func (r *IllnessRepository) Delete(ill models.Illness) error {
	err := r.db.
		Delete(&ill).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *IllnessRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var ill []models.Illness
	var illres []responses.Illness

	result := r.
		preload().
		Model(&ill).
		Select(SELECTED_FIELDS)

	p.Search = fmt.Sprintf("%%%s%%", p.Search)
	if p.Search != "" {
		result = result.
			Where(r.db.
				Where(&models.Illness{Name: p.Search}).
				Or(&models.Illness{Description: p.Search}),
			)
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
		Scopes(paginationrepo.Paginate(&ill, &p, result)).
		Find(&ill)

	if err := result.Error; err != nil {
		utlogger.LogError(result.Error)
		return nil, result.Error
	}

	// * copy the data from model to response
	copier.Copy(&illres, &ill)

	p.Data = illres
	return &p, nil
}

func (r *IllnessRepository) FindByID(illid uuid.UUID) (*models.Illness, error) {
	var illness *models.Illness
	err := r.
		preload().
		Model(&models.Illness{}).
		Group("id").
		First(&illness, illid).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return illness, nil
}

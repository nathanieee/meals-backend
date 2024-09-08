package illnessrepo

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
		GetByID(id uuid.UUID) (*models.Illness, error)
	}
)

func NewIllnessRepository(db *gorm.DB) *IllnessRepository {
	return &IllnessRepository{db: db}
}

func (r *IllnessRepository) omit() *gorm.DB {
	return r.db.
		Omit(
			"",
		)
}

func (r *IllnessRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *IllnessRepository) Create(ill models.Illness) (*models.Illness, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&ill).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	illnew, err := r.GetByID(ill.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return illnew, nil
}

func (r *IllnessRepository) Read() ([]*models.Illness, error) {
	var (
		ill []*models.Illness
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&ill).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return ill, nil
}

func (r *IllnessRepository) Update(ill models.Illness) (*models.Illness, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&ill).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	illnew, err := r.GetByID(ill.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return illnew, nil
}

func (r *IllnessRepository) Delete(ill models.Illness) error {
	err := r.db.
		Delete(&ill).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *IllnessRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		ill    []models.Illness
		illres []responses.Illness
	)

	result := r.
		preload().
		Model(&ill).
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

	result = result.
		Group("id").
		Scopes(paginationrepo.Paginate(&ill, &p, result)).
		Find(&ill)

	if err := result.Error; err != nil {
		utlogger.Error(result.Error)
		return nil, result.Error
	}

	// * copy the data from model to response
	copier.CopyWithOption(&illres, &ill, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = illres
	return &p, nil
}

func (r *IllnessRepository) GetByID(id uuid.UUID) (*models.Illness, error) {
	var (
		ill *models.Illness
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Illness{Model: base.Model{ID: id}}).
		First(&ill).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return ill, nil
}

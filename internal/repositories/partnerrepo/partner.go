package partnerrepo

import (
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
		user_id,
		name, 
		created_at,
		updated_at
	`
)

type (
	PartnerRepository struct {
		db *gorm.DB
	}

	IPartnerRepository interface {
		Create(p models.Partner) (*models.Partner, error)
		Read() ([]*models.Partner, error)
		Update(p models.Partner) (*models.Partner, error)
		Delete(p models.Partner) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(pid uuid.UUID) (*models.Partner, error)
	}
)

func NewPartnerRepository(db *gorm.DB) *PartnerRepository {
	return &PartnerRepository{db: db}
}

func (r *PartnerRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("User.Image.Image").
		Preload("User.Address")
}

func (r *PartnerRepository) Create(p models.Partner) (*models.Partner, error) {
	err := r.db.
		Create(&p).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &p, nil
}

func (r *PartnerRepository) Read() ([]*models.Partner, error) {
	var p []*models.Partner

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&p).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return p, nil
}

func (r *PartnerRepository) Update(p models.Partner) (*models.Partner, error) {
	err := r.db.
		Save(&p).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &p, nil
}

func (r *PartnerRepository) Delete(p models.Partner) error {
	err := r.db.
		Delete(&p).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *PartnerRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var pa []*models.Partner
	var pares []responses.Partner

	result := r.
		preload().
		Model(&pa).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		result = result.
			Where(r.db.
				Where(&models.Partner{Name: p.Search}),
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
		Scopes(paginationrepo.Paginate(&pa, &p, result)).
		Find(&pa)

	if result.Error != nil {
		utlogger.LogError(result.Error)
		return nil, result.Error
	}

	copier.Copy(&pares, &pa)

	p.Data = pares
	return &p, nil
}

func (r *PartnerRepository) FindByID(pid uuid.UUID) (*models.Partner, error) {
	var p *models.Partner

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Partner{Model: helper.Model{ID: pid}}).
		First(&p).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return p, nil
}

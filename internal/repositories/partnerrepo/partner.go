package partnerrepo

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
		FindByEmail(email string) (*models.Partner, error)
		FindByUserID(uid uuid.UUID) (*models.Partner, error)
	}
)

func NewPartnerRepository(db *gorm.DB) *PartnerRepository {
	return &PartnerRepository{db: db}
}

func (r *PartnerRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("User.Address").
		Preload("User.Image.Image")
}

func (r *PartnerRepository) Create(p models.Partner) (*models.Partner, error) {
	err := r.db.
		Create(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	pnew, err := r.FindByID(p.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return pnew, nil
}

func (r *PartnerRepository) Read() ([]*models.Partner, error) {
	var (
		p []*models.Partner
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return p, nil
}

func (r *PartnerRepository) Update(p models.Partner) (*models.Partner, error) {
	err := r.db.
		Save(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	pnew, err := r.FindByID(p.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return pnew, nil
}

func (r *PartnerRepository) Delete(p models.Partner) error {
	err := r.db.
		Delete(&p).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *PartnerRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		pa    []*models.Partner
		pares []responses.Partner
	)

	result := r.
		preload().
		Model(&pa).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where("name LIKE ?", p.Search),
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

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&pares, &pa, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = pares
	return &p, nil
}

func (r *PartnerRepository) FindByID(pid uuid.UUID) (*models.Partner, error) {
	var (
		p *models.Partner
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Partner{Model: base.Model{ID: pid}}).
		First(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return p, nil
}

func (r *PartnerRepository) FindByEmail(email string) (*models.Partner, error) {
	var (
		p *models.Partner
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Partner{User: models.User{Email: email}}).
		First(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return p, nil
}

func (r *PartnerRepository) FindByUserID(uid uuid.UUID) (*models.Partner, error) {
	var (
		p *models.Partner
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Partner{User: models.User{Model: base.Model{ID: uid}}}).
		First(&p).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return p, nil
}

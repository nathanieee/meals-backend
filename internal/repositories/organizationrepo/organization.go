package organizationrepo

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
		user_id,
		type,
		name,
		created_at,
		updated_at
	`
)

type (
	OrganizationRepository struct {
		db *gorm.DB
	}

	IOrganizationRepository interface {
		Create(o models.Organization) (*models.Organization, error)
		Read() ([]*models.Organization, error)
		Update(o models.Organization) (*models.Organization, error)
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(oid uuid.UUID) (*models.Organization, error)
		FindByEmail(email string) (*models.Organization, error)
	}
)

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("User.Address").
		Preload("User.UserImage.Image")
}

func (r *OrganizationRepository) Create(o models.Organization) (*models.Organization, error) {
	err := r.db.
		Create(o).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &o, nil
}

func (r *OrganizationRepository) Read() ([]*models.Organization, error) {
	var o []*models.Organization

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&o).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return o, nil
}

func (r *OrganizationRepository) Update(o models.Organization) (*models.Organization, error) {
	err := r.db.
		Save(&o).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &o, nil
}

func (r *OrganizationRepository) Delete(o models.Organization) error {
	err := r.db.
		Delete(&o).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *OrganizationRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var o []models.Organization
	var ores []responses.Organization

	result := r.
		preload().
		Model(&o).
		Select(SELECTED_FIELDS)

	p.Search = fmt.Sprintf("%%%s%%", p.Search)
	if p.Search != "" {
		result = result.
			Where(r.db.
				Where(&models.Organization{Name: p.Search}),
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
		Scopes(paginationrepo.Paginate(&o, &p, result)).
		Find(&o)

	if err := result.Error; err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.Copy(&ores, &o)

	p.Data = ores
	return &p, result.Error
}

func (r *OrganizationRepository) FindByID(oid uuid.UUID) (*models.Organization, error) {
	var o *models.Organization

	err := r.db.
		Model(&models.Organization{}).
		Select(SELECTED_FIELDS).
		First(&o, oid).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return o, nil
}

func (r *OrganizationRepository) FindByEmail(email string) (*models.Organization, error) {
	var o *models.Organization

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Organization{User: models.User{Email: email}}).
		First(&o).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return o, nil
}

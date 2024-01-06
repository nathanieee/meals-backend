package orgrepository

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/logger"

	"github.com/google/uuid"
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
		FindByID(oid uuid.UUID) (*models.Organization, error)
	}
)

func NewOrganizationRepository(db *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{db: db}
}

func (r *OrganizationRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations)
}

func (r *OrganizationRepository) Create(o models.Organization) (*models.Organization, error) {
	err := r.db.Create(o).Error
	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return &o, nil
}

func (r *OrganizationRepository) FindByID(oid uuid.UUID) (*models.Organization, error) {
	var o *models.Organization

	err := r.db.
		Model(&models.Organization{}).
		Select(SELECTED_FIELDS).
		First(&o, oid).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return o, nil
}

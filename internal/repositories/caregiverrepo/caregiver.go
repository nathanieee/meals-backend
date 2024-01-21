package caregiverrepo

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
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		user_id,
		gender,
		first_name,
		last_name,
		date_of_birth,
		created_at,
		updated_at
	`
)

type (
	CaregiverRepository struct {
		db *gorm.DB
	}

	ICaregiverRepository interface {
		Create(cg models.Caregiver) (*models.Caregiver, error)
		Read() ([]*models.Caregiver, error)
		Update(cg models.Caregiver) (*models.Caregiver, error)
		Delete(cg models.Caregiver) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(cgid uuid.UUID) (*models.Caregiver, error)
		FindByEmail(email string) (*models.Caregiver, error)
	}
)

func NewCaregiverRepository(db *gorm.DB) *CaregiverRepository {
	return &CaregiverRepository{db: db}
}

func (r *CaregiverRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("User").
		Preload("User.UserImages.Images").
		Preload("User.Addresses")
}

func (r *CaregiverRepository) Create(cg models.Caregiver) (*models.Caregiver, error) {
	err := r.db.
		Create(cg).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &cg, err
}

func (r *CaregiverRepository) Read() ([]*models.Caregiver, error) {
	var cg []*models.Caregiver

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&cg).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return cg, nil
}

func (r *CaregiverRepository) Update(cg models.Caregiver) (*models.Caregiver, error) {
	err := r.db.
		Save(&cg).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &cg, nil
}

func (r *CaregiverRepository) Delete(cg models.Caregiver) error {
	err := r.db.
		Delete(&cg).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *CaregiverRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var cg []models.Caregiver
	var cgres []responses.CaregiverResponse

	result := r.
		preload().
		Model(&cg).
		Select(SELECTED_FIELDS)

	p.Search = fmt.Sprintf("%%%s%%", p.Search)
	if p.Search != "" {
		result = result.
			Where(r.db.
				Where(&models.Caregiver{FirstName: p.Search}).
				Or(&models.Caregiver{LastName: p.Search}),
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
		Scopes(paginationrepo.Paginate(&cg, &p, result)).
		Find(&cgres)

	if result.Error != nil {
		utlogger.LogError(result.Error)
		return nil, result.Error
	}

	p.Data = cgres
	return &p, nil
}

func (r *CaregiverRepository) FindByID(cgid uuid.UUID) (*models.Caregiver, error) {
	var cg *models.Caregiver

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Caregiver{Model: helper.Model{ID: cgid}}).
		First(&cg).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return cg, nil
}

func (r *CaregiverRepository) FindByEmail(email string) (*models.Caregiver, error) {
	var cg *models.Caregiver

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Caregiver{User: models.User{Email: email}}).
		First(&cg).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return cg, nil
}

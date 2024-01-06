package crgvrrepository

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils"
	"project-skbackend/packages/utils/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	CaregiverRepository struct {
		db *gorm.DB
	}

	ICaregiverRepository interface {
		Create(cg models.Caregiver) (*models.Caregiver, error)
		Update(cg models.Caregiver, cgid uuid.UUID) (*models.Caregiver, error)
		FindAll(p utils.Pagination) (*utils.Pagination, error)
		FindByID(cgid uuid.UUID) (*responses.CaregiverResponse, error)
		FindByEmail(email string) (*responses.CaregiverResponse, error)
		Delete(cg models.Caregiver) error
	}
)

func NewCaregiverRepository(db *gorm.DB) *CaregiverRepository {
	return &CaregiverRepository{db: db}
}

func (r *CaregiverRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations).
		Preload("User").
		Preload("User.UserImages.Images").
		Preload("User.Addresses")
}

func (r *CaregiverRepository) Create(cg models.Caregiver) (*models.Caregiver, error) {
	err := r.db.
		Create(cg).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return &cg, err
}

func (r *CaregiverRepository) Update(cg models.Caregiver, cgid uuid.UUID) (*models.Caregiver, error) {
	err := r.db.
		Model(&cg).
		Where("id = ?", cgid).
		Updates(cg).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return &cg, nil
}

func (r *CaregiverRepository) FindAll(p utils.Pagination) (*utils.Pagination, error) {
	var cg []models.Caregiver
	var cgres []responses.CaregiverResponse

	result := r.preload(r.db).
		Model(&cg).
		Select("id, user_id, gender, first_name, last_name, date_of_birth, created_at, updated_at")

	if p.Search != "" {
		result = result.
			Where("first_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).
			Or("last_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.
			Where("date(created_at) between ? and ?",
				p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT),
				p.Filter.CreatedTo.Format(consttypes.DATEFORMAT))
	}

	result = result.
		Group("id").
		Scopes(pagination.Paginate(&cg, &p, result)).
		Find(&cgres)

	if result.Error != nil {
		logger.LogError(result.Error)
		return nil, result.Error
	}

	p.Data = cgres
	return &p, nil
}

func (r *CaregiverRepository) FindByID(cgid uuid.UUID) (*responses.CaregiverResponse, error) {
	var cgres *responses.CaregiverResponse
	err := r.preload(r.db).
		Model(&models.Caregiver{}).
		Select("id, user_id, gender, first_name, last_name, date_of_birth, created_at, updated_at").
		First(&cgres, cgid).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return cgres, nil
}

func (r *CaregiverRepository) FindByEmail(email string) (*responses.CaregiverResponse, error) {
	var cgres *responses.CaregiverResponse
	err := r.preload(r.db).
		Model(&models.Caregiver{}).
		Select("id, user_id, gender, first_name, last_name, date_of_birth, created_at, updated_at").
		Where("users.email = ?", email).
		Group("id").
		Take(&cgres).Error

	if err != nil {
		logger.LogError(err)
		return nil, err
	}

	return cgres, nil
}

func (r *CaregiverRepository) Delete(cg models.Caregiver) error {
	err := r.db.
		Delete(&cg).Error

	if err != nil {
		logger.LogError(err)
		return err
	}

	return nil
}

package allgrepository

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	AllergyRepository struct {
		db *gorm.DB
	}

	IAllergyRepository interface {
		Create(al models.Allergy) (*models.Allergy, error)
		FindByID(alid uuid.UUID) (*models.Allergy, error)
	}
)

func NewAllergyRepository(db *gorm.DB) *AllergyRepository {
	return &AllergyRepository{db: db}
}

func (r *AllergyRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations)
}

func (r *AllergyRepository) Create(al models.Allergy) (*models.Allergy, error) {
	err := r.db.Create(al).Error
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &al, nil
}

func (r *AllergyRepository) FindByID(alid uuid.UUID) (*models.Allergy, error) {
	var ally *models.Allergy
	err := r.
		preload(r.db).
		Model(&models.Allergy{}).
		Group("id").
		First(&ally, alid).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return ally, nil
}

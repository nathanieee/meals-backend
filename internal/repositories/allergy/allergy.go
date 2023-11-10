package allgrepository

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	AllergyRepository struct {
		db *gorm.DB
	}

	IAllergyRepository interface {
		Create(al *models.Allergy) (*models.Allergy, error)
		FindByID(alid uuid.UUID) (*responses.AllergyResponse, error)
	}
)

func NewAllergyRepository(db *gorm.DB) *AllergyRepository {
	return &AllergyRepository{db: db}
}

func (r *AllergyRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations)
}

func (r *AllergyRepository) Create(al *models.Allergy) (*models.Allergy, error) {
	err := r.db.Create(al).Error
	if err != nil {
		return nil, err
	}

	return al, nil
}

func (r *AllergyRepository) FindByID(alid uuid.UUID) (*responses.AllergyResponse, error) {
	var alres *responses.AllergyResponse
	err := r.db.
		Model(&models.Allergy{}).
		Select("id, name, created_at, updated_at").
		Group("id").
		First(&alres, alid).Error

	if err != nil {
		return nil, err
	}

	return alres, nil
}

package allergy

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AllergyRepo struct {
	db *gorm.DB
}

func NewAllergyRepo(db *gorm.DB) *AllergyRepo {
	return &AllergyRepo{db: db}
}

func (alr *AllergyRepo) Create(al *models.Allergy) (*models.Allergy, error) {
	err := alr.db.Create(al).Error
	if err != nil {
		return nil, err
	}

	return al, nil
}

func (alr *AllergyRepo) FindByID(alid uuid.UUID) (*responses.AllergyResponse, error) {
	var alres *responses.AllergyResponse
	err := alr.db.
		Model(&models.Allergy{}).
		Select("id, name, created_at, updated_at").
		Group("id").
		First(&alres, alid).Error

	if err != nil {
		return nil, err
	}

	return alres, nil
}

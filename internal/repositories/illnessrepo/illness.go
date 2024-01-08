package illnessrepo

import (
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	IllnessRepository struct {
		db *gorm.DB
	}

	IIllnessRepository interface {
		Create(ill models.Illness) (*models.Illness, error)
		FindByID(illid uuid.UUID) (*models.Illness, error)
	}
)

func NewIllnessRepository(db *gorm.DB) *IllnessRepository {
	return &IllnessRepository{db: db}
}

func (r *IllnessRepository) preload(db *gorm.DB) *gorm.DB {
	return db.
		Preload(clause.Associations)
}

func (r *IllnessRepository) Create(ill models.Illness) (*models.Illness, error) {
	err := r.db.Create(ill).Error
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &ill, nil
}

func (r *IllnessRepository) FindByID(illid uuid.UUID) (*models.Illness, error) {
	var illness *models.Illness
	err := r.
		preload(r.db).
		Model(&models.Illness{}).
		Group("id").
		First(&illness, illid).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return illness, nil
}

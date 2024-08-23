package imagerepo

import (
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
	"project-skbackend/packages/utils/utlogger"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		name, 
		path,
		type,
		created_at,
		updated_at
	`
)

type (
	ImageRepository struct {
		db *gorm.DB
	}

	IImageRepo interface {
		Create(i models.Image) (*models.Image, error)
		Read() ([]*models.Image, error)
		Update(i models.Image) (*models.Image, error)
		Delete(i models.Image) error
		GetByID(id uuid.UUID) (*models.Image, error)
	}
)

func NewImageRepository(db *gorm.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) omit() *gorm.DB {
	return r.db.Omit(
		"",
	)
}

func (r *ImageRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *ImageRepository) Create(i models.Image) (*models.Image, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&i).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	inew, err := r.GetByID(i.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return inew, nil
}

func (r *ImageRepository) Read() ([]*models.Image, error) {
	var (
		i []*models.Image
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&i).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return i, nil
}

func (r *ImageRepository) Update(i models.Image) (*models.Image, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Save(&i).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	inew, err := r.GetByID(i.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return inew, nil
}

func (r *ImageRepository) Delete(i models.Image) error {
	err := r.db.
		Delete(&i).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *ImageRepository) GetByID(id uuid.UUID) (*models.Image, error) {
	var (
		i *models.Image
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.Image{Model: base.Model{ID: id}}).
		First(&i).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return i, nil
}

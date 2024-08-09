package userimagerepo

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
		user_id,
		image_id,
		created_at,
		updated_at
	`
)

type (
	UserImageRepository struct {
		db *gorm.DB
	}

	IUserImageRepo interface {
		Create(ui models.UserImage) (*models.UserImage, error)
		Read() ([]*models.UserImage, error)
		Update(ui models.UserImage) (*models.UserImage, error)
		Delete(ui models.UserImage) error
		GetByID(id uuid.UUID) (*models.UserImage, error)
		GetByUserID(uid uuid.UUID) (*models.UserImage, error)
	}
)

func NewUserImageRepository(db *gorm.DB) *UserImageRepository {
	return &UserImageRepository{db: db}
}

func (r *UserImageRepository) omit() *gorm.DB {
	return r.db.Omit(
		"",
	)
}

func (r *UserImageRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations)
}

func (r *UserImageRepository) Create(ui models.UserImage) (*models.UserImage, error) {
	err := r.
		omit().
		Session(&gorm.Session{FullSaveAssociations: true}).
		Create(&ui).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	uinew, err := r.GetByID(ui.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return uinew, nil
}

func (r *UserImageRepository) Read() ([]*models.UserImage, error) {
	var (
		ui []*models.UserImage
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&ui).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return ui, nil
}

func (r *UserImageRepository) Update(ui models.UserImage) (*models.UserImage, error) {
	err := r.db.
		Save(&ui).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	uinew, err := r.GetByID(ui.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return uinew, nil
}

func (r *UserImageRepository) Delete(ui models.UserImage) error {
	err := r.db.
		Delete(&ui).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *UserImageRepository) GetByID(id uuid.UUID) (*models.UserImage, error) {
	var (
		ui *models.UserImage
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.UserImage{Model: base.Model{ID: id}}).
		First(&ui).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return ui, nil
}

func (r *UserImageRepository) GetByUserID(uid uuid.UUID) (*models.UserImage, error) {
	var (
		ui *models.UserImage
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.UserImage{UserID: uid}).
		First(&ui).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return ui, nil
}

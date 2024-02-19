package userrepo

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
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		email,
		password,
		role,
		reset_password_token,
		reset_password_sent_at,
		created_at,
		updated_at
	`
)

type (
	UserRepository struct {
		db *gorm.DB
	}

	IUserRepository interface {
		Create(u models.User) (*models.User, error)
		Read() ([]*models.User, error)
		Update(u models.User) (*models.User, error)
		Delete(u models.User) error
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		FindByID(uid uuid.UUID) (*models.User, error)
		FindByEmail(email string) (*models.User, error)
		FirstOrCreate(u models.User) (*models.User, error)
	}
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.
		Preload(clause.Associations).
		Preload("Image.Image").
		Preload("Address")

	return &UserRepository{db: db}
}

func (r *UserRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Address").
		Preload("Image.Image")
}

func (r *UserRepository) Create(u models.User) (*models.User, error) {
	err := r.db.Create(&u).Error
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Read() ([]*models.User, error) {
	var u []*models.User

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&u).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Update(u models.User) (*models.User, error) {
	err := r.db.
		Model(&u).
		Updates(u).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) Delete(u models.User) error {
	err := r.db.
		Delete(&u).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (r *UserRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var u []models.User
	var ures []responses.User

	result := r.
		preload().
		Model(&u).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(r.db.
				Where(&models.User{Email: p.Search}),
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
		Scopes(paginationrepo.Paginate(&u, &p, result)).
		Find(&u)

	if result.Error != nil {
		utlogger.LogError(result.Error)
		return nil, result.Error
	}

	// * copy the data from model to response
	copier.CopyWithOption(&ures, &u, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = ures
	return &p, nil
}

func (r *UserRepository) FindByID(uid uuid.UUID) (*models.User, error) {
	var u models.User
	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.User{Model: helper.Model{ID: uid}}).
		First(&u).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.User{Email: email}).
		First(&u).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &u, nil
}

func (r *UserRepository) FirstOrCreate(u models.User) (*models.User, error) {
	err := r.db.
		FirstOrCreate(&u, u).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &u, nil
}

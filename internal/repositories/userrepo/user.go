package userrepo

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/base"
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
		confirmed_at,
		confirmation_token,
		confirmation_sent_at,
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
		GetByID(id uuid.UUID) (*models.User, error)
		GetByEmail(email string) (*models.User, error)
		FirstOrCreate(u models.User) (*models.User, error)
	}
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) preload() *gorm.DB {
	return r.db.
		Preload(clause.Associations).
		Preload("Addresses").
		Preload("Image.Image")
}

func (r *UserRepository) Create(u models.User) (*models.User, error) {
	err := r.db.
		Create(&u).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	unew, err := r.GetByID(u.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return unew, nil
}

func (r *UserRepository) Read() ([]*models.User, error) {
	var (
		u []*models.User
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Find(&u).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Update(u models.User) (*models.User, error) {
	err := r.db.
		Save(&u).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	unew, err := r.GetByID(u.ID)

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return unew, nil
}

func (r *UserRepository) Delete(u models.User) error {
	err := r.db.
		Delete(&u).Error

	if err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (r *UserRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var (
		u    []models.User
		ures []responses.User
	)

	result := r.
		preload().
		Model(&u).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		p.Search = fmt.Sprintf("%%%s%%", p.Search)
		result = result.
			Where(
				r.db.Where(`
					email ILIKE ?
		`, p.Search),
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

	if err := result.Error; err != nil {
		utlogger.Error(err)
		return nil, err
	}

	// * copy the data from model to response
	copier.CopyWithOption(&ures, &u, copier.Option{IgnoreEmpty: true, DeepCopy: true})

	p.Data = ures
	return &p, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	var (
		u *models.User
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.User{Model: base.Model{ID: id}}).
		First(&u).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	var (
		u *models.User
	)

	err := r.
		preload().
		Select(SELECTED_FIELDS).
		Where(&models.User{Email: email}).
		First(&u).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return u, nil
}

func (r *UserRepository) FirstOrCreate(u models.User) (*models.User, error) {
	err := r.db.
		FirstOrCreate(&u, u).Error

	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return &u, nil
}

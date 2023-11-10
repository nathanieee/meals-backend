package userrepository

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	SELECTED_FIELDS = `
		id,
		email,
		role,
		password,
		reset_password_token,
		reset_password_sent_at,
		confirmation_token,
		confirmed_at,
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
		FindAll(p models.Pagination) (*models.Pagination, error)
		Create(user *models.User) (*models.User, error)
		Update(user models.User, uid uuid.UUID) (*models.User, error)
		FindByID(uid uuid.UUID) (*responses.UserResponse, error)
		FindByEmail(email string) (*responses.UserResponse, error)
		Delete(user models.User) error
	}
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.
		Preload(clause.Associations).
		Preload("UserImages.Images").
		Preload("Addresses")

	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(u *models.User) (*models.User, error) {
	err := ur.db.Create(u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Update(u models.User, uid uuid.UUID) (*models.User, error) {
	err := ur.db.
		Model(&u).
		Where("id = ?", uid).
		Updates(u).Error

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (ur *UserRepository) FindAll(p models.Pagination) (*models.Pagination, error) {
	var u []models.User
	var ures []responses.UserResponse

	result := ur.db.Model(&u).
		Select(SELECTED_FIELDS)

	if p.Search != "" {
		result = result.
			Where("full_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).
			Or("email LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
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
		Scopes(pagination.Paginate(&u, &p, result)).
		Find(&ures)

	if result.Error != nil {
		return nil, result.Error
	}

	p.Data = ures
	return &p, nil
}

func (ur *UserRepository) FindByID(uid uuid.UUID) (*responses.UserResponse, error) {
	var ures *responses.UserResponse
	err := ur.db.
		Model(&models.User{}).
		Select(SELECTED_FIELDS).
		Group("id").
		First(&ures, uid).Error

	if err != nil {
		return nil, err
	}

	return ures, nil
}

func (ur *UserRepository) FindByEmail(email string) (*responses.UserResponse, error) {
	var ures *responses.UserResponse
	err := ur.db.
		Model(&models.User{}).
		Select(SELECTED_FIELDS).
		Where("email = ?", email).
		Group("id").
		Take(&ures).Error

	if err != nil {
		return nil, err
	}

	return ures, nil
}

func (ur *UserRepository) Delete(u models.User) error {
	err := ur.db.
		Delete(&u).Error

	if err != nil {
		return err
	}

	return nil
}

package userrepository

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

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
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		Create(user models.User) (*models.User, error)
		Update(user models.User, uid uuid.UUID) (*models.User, error)
		FindByID(uid uuid.UUID) (*models.User, error)
		FindByEmail(email string) (*models.User, error)
		Delete(user models.User) error
		FirstOrCreate(user models.User) (*models.User, error)
	}
)

func NewUserRepository(db *gorm.DB) *UserRepository {
	db.
		Preload(clause.Associations).
		Preload("UserImages.Images").
		Preload("Addresses")

	return &UserRepository{db: db}
}

func (ur *UserRepository) Create(user models.User) (*models.User, error) {
	err := ur.db.Create(&user).Error
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Update(user models.User, uid uuid.UUID) (*models.User, error) {
	err := ur.db.
		Model(&user).
		Where("id = ?", uid).
		Updates(user).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	var user []models.User
	var ures []responses.UserResponse

	result := ur.db.Model(&user).
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
		Scopes(pagination.Paginate(&user, &p, result)).
		Find(&ures)

	if result.Error != nil {
		utlogger.LogError(result.Error)
		return nil, result.Error
	}

	p.Data = ures
	return &p, nil
}

func (ur *UserRepository) FindByID(uid uuid.UUID) (*models.User, error) {
	var user models.User
	err := ur.db.
		Model(&models.User{}).
		Select(SELECTED_FIELDS).
		Group("id").
		First(&user, uid).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := ur.db.
		Model(&models.User{}).
		Select(SELECTED_FIELDS).
		Where("email = ?", email).
		Group("id").
		Take(&user).Error

	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &user, nil
}

func (ur *UserRepository) Delete(user models.User) error {
	err := ur.db.
		Delete(&user).Error

	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (ur *UserRepository) FirstOrCreate(user models.User) (*models.User, error) {
	err := ur.db.FirstOrCreate(&user, user).Error
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return &user, nil
}

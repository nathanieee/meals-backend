package user

import (
	"fmt"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/pagination"
	"project-skbackend/packages/consttypes"

	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (ur *UserRepo) Store(u *models.User) (*models.User, error) {
	err := ur.db.Create(u).Error
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (ur *UserRepo) Update(u models.User, uid uint) (*models.User, error) {
	err := ur.db.Model(&u).Where("id = ?", uid).Updates(u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (ur *UserRepo) FindAll(p models.Pagination) (*models.Pagination, error) {
	var users []models.User
	var usersResponse []responses.UserResponse

	result := ur.db.Model(&users).Select("users.id as id, full_name, email, password, reset_password_token, reset_password_sent_at, confirmation_token, confirmed_at, confirmation_sent_at, users.created_at as created_at, users.updated_at as updated_at")

	if p.Search != "" {
		result = result.Where("full_name LIKE ?", fmt.Sprintf("%%%s%%", p.Search)).Or("email LIKE ?", fmt.Sprintf("%%%s%%", p.Search))
	}

	if !p.Filter.CreatedFrom.IsZero() && !p.Filter.CreatedTo.IsZero() {
		result = result.Where("date(users.created_at) between ? and ?", p.Filter.CreatedFrom.Format(consttypes.DATEFORMAT), p.Filter.CreatedTo.Format(consttypes.DATEFORMAT))
	}

	result = result.Group("users.id").Scopes(pagination.Paginate(&users, &p, result)).Find(&usersResponse)

	if result.Error != nil {
		return &p, result.Error
	}

	p.Data = usersResponse
	return &p, nil
}

func (ur *UserRepo) FindByID(uid uint) (*responses.UserResponse, error) {
	var u *responses.UserResponse
	err := ur.db.Model(&models.User{}).Select("users.id as id, full_name, email, password, reset_password_token, reset_password_sent_at, confirmation_token, confirmed_at, confirmation_sent_at, refresh_token, refresh_token_expiration, users.created_at as created_at, users.updated_at as updated_at").Group("users.id").First(&u, uid).Error
	if err != nil {
		return nil, err
	}

	return u, err
}

func (ur *UserRepo) FindByEmail(email string) (*responses.UserResponse, error) {
	var u *responses.UserResponse
	err := ur.db.Model(&models.User{}).Select("users.id as id, full_name, email, password, reset_password_token, reset_password_sent_at, confirmation_token, confirmed_at, confirmation_sent_at, users.created_at as created_at, users.updated_at as updated_at").Where("email = ?", email).Group("users.id").Take(&u).Error
	if err != nil {
		return nil, err
	}

	return u, err
}

func (ur *UserRepo) DeleteUser(u models.User) error {
	err := ur.db.Unscoped().Delete(&u).Error
	if err != nil {
		return err
	}

	return nil
}

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

func (u *UserRepo) Store(user *models.User) (*models.User, error) {
	err := u.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserRepo) Update(user models.User, userID uint) (*models.User, error) {
	err := u.db.Model(&user).Where("id = ?", userID).Updates(user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (u *UserRepo) FindAll(p models.Pagination) (*models.Pagination, error) {
	var users []models.User
	var usersResponse []responses.UserResponse

	result := u.db.Model(&users).Select("users.id as id, full_name, email, password, user_level, country, country_code, reset_password_token, reset_password_sent_at, confirmation_token, confirmed_at, confirmation_sent_at, users.created_at as created_at, users.updated_at as updated_at")

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

func (u *UserRepo) FindById(id uint) (*responses.UserResponse, error) {
	var user *responses.UserResponse
	err := u.db.Model(&models.User{}).Select("users.id as id, full_name, email, password, user_level, country, country_code, count(distinct(scenarios.id)) as scenario_count, reset_password_token, reset_password_sent_at, confirmation_token, confirmed_at, confirmation_sent_at, refresh_token, refresh_token_expiration, users.created_at as created_at, users.updated_at as updated_at").Group("users.id").Where("scenarios.deleted_at is null").First(&user, id).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *UserRepo) FindByEmail(email string) (*responses.UserResponse, error) {
	var user *responses.UserResponse
	err := u.db.Model(&models.User{}).Select("users.id as id, full_name, email, password, user_level, country, country_code, reset_password_token, reset_password_sent_at, confirmation_token, confirmed_at, confirmation_sent_at, users.created_at as created_at, users.updated_at as updated_at").Where("email = ?", email).Group("users.id").Take(&user).Error
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *UserRepo) DeleteUser(user models.User) error {
	err := u.db.Unscoped().Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}

package user

import (
	"encoding/json"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories"
	"project-skbackend/packages/utils"

	"gorm.io/gorm"
)

type UserService struct {
	ur repositories.IUserRepo
}

func NewUserService(ur repositories.IUserRepo) *UserService {
	return &UserService{ur: ur}
}

func (us *UserService) CreateUser(req requests.CreateUserRequest) (*responses.UserResponse, error) {
	var ures *responses.UserResponse

	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName: req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	user, err = us.ur.Store(user)
	if err != nil {
		return nil, err
	}

	marshaledUser, _ := json.Marshal(user)
	err = json.Unmarshal(marshaledUser, &ures)
	if err != nil {
		return nil, err
	}

	return ures, err
}

func (us *UserService) GetUser(uid uint) (*responses.UserResponse, error) {
	user, err := us.ur.FindByID(uid)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (us *UserService) GetUsers(paginationReq models.Pagination) (*models.Pagination, error) {
	users, err := us.ur.FindAll(paginationReq)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) DeleteUser(uid uint) error {
	userModel := models.User{
		Model: gorm.Model{ID: uid},
	}

	err := us.ur.DeleteUser(userModel)
	if err != nil {
		return err
	}

	return nil
}

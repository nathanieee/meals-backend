package user

import (
	"encoding/json"
	"fmt"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories"
	"project-skbackend/packages/utils"
)

type UserService struct {
	userRepo repositories.IUserRepo
}

func NewUserService(userRepo repositories.IUserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (u *UserService) CreateUser(req requests.CreateUserRequest) (*responses.UserResponse, error) {
	var userResponse *responses.UserResponse

	hashedPassword, err := utils.EncryptPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FullName: req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	user, err = u.userRepo.Store(user)
	if err != nil {
		return nil, err
	}

	marshaledUser, _ := json.Marshal(user)
	err = json.Unmarshal(marshaledUser, &userResponse)
	if err != nil {
		fmt.Println("err", err)
	}

	return userResponse, err
}

func (u *UserService) GetUser(id uint) (*responses.UserResponse, error) {
	user, err := u.userRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (u *UserService) GetUsers(paginationReq models.Pagination) (*models.Pagination, error) {
	users, err := u.userRepo.FindAll(paginationReq)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (u *UserService) DeleteUser(id uint) error {
	userModel := models.User{
		ID: id,
	}
	err := u.userRepo.DeleteUser(userModel)
	if err != nil {
		return err
	}

	return err
}

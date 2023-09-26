package user

import (
	"encoding/json"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/internal/repositories"

	"github.com/google/uuid"
)

type UserService struct {
	ur repositories.IUserRepo
}

func NewUserService(ur repositories.IUserRepo) *UserService {
	return &UserService{ur: ur}
}

func (us *UserService) Create(req requests.CreateUserRequest) (*responses.UserResponse, error) {
	var ures *responses.UserResponse

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}
	user, err := us.ur.Create(user)
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

func (us *UserService) FindByID(uid uuid.UUID) (*responses.UserResponse, error) {
	user, err := us.ur.FindByID(uid)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (us *UserService) FindAll(paginationReq models.Pagination) (*models.Pagination, error) {
	users, err := us.ur.FindAll(paginationReq)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) Delete(uid uuid.UUID) error {
	userModel := models.User{
		Model: helper.Model{ID: uid},
	}

	err := us.ur.Delete(userModel)
	if err != nil {
		return err
	}

	return nil
}

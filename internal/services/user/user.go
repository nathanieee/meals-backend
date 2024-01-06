package userservice

import (
	"encoding/json"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	userrepository "project-skbackend/internal/repositories/user"
	"project-skbackend/packages/utils"

	"github.com/google/uuid"
)

type (
	UserService struct {
		userrepo userrepository.IUserRepository
	}

	IUserService interface {
		Create(req requests.CreateUserRequest) (*responses.UserResponse, error)
		FindByID(id uuid.UUID) (*responses.UserResponse, error)
		FindAll(paginationReq utils.Pagination) (*utils.Pagination, error)
		Delete(id uuid.UUID) error
	}
)

func NewUserService(
	userrepo userrepository.IUserRepository,
) *UserService {
	return &UserService{
		userrepo: userrepo,
	}
}

func (us *UserService) Create(req requests.CreateUserRequest) (*responses.UserResponse, error) {
	var ures *responses.UserResponse

	user := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	user, err := us.userrepo.Create(*user)
	if err != nil {
		return nil, err
	}

	maru, _ := json.Marshal(user)
	err = json.Unmarshal(maru, &ures)
	if err != nil {
		return nil, err
	}

	return ures, err
}

func (us *UserService) FindByID(uid uuid.UUID) (*responses.UserResponse, error) {
	user, err := us.userrepo.FindByID(uid)
	if err != nil {
		return nil, err
	}

	return user.ToResponse(), err
}

func (us *UserService) FindAll(paginationReq utils.Pagination) (*utils.Pagination, error) {
	users, err := us.userrepo.FindAll(paginationReq)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (us *UserService) Delete(uid uuid.UUID) error {
	userModel := models.User{
		Model: helper.Model{ID: uid},
	}

	err := us.userrepo.Delete(userModel)
	if err != nil {
		return err
	}

	return nil
}

package userservice

import (
	"encoding/json"
	"project-skbackend/internal/controllers/requests"
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
	"project-skbackend/internal/models/helper"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils/utlogger"
	"project-skbackend/packages/utils/utpagination"

	"github.com/google/uuid"
)

type (
	UserService struct {
		userrepo userrepo.IUserRepository
	}

	IUserService interface {
		Create(req requests.CreateUserRequest) (*responses.UserResponse, error)
		FindByID(uid uuid.UUID) (*responses.UserResponse, error)
		FindAll(p utpagination.Pagination) (*utpagination.Pagination, error)
		Delete(uid uuid.UUID) error
		Update(req requests.UpdateUserRequest, uid uuid.UUID) (*responses.UserResponse, error)
	}
)

func NewUserService(
	userrepo userrepo.IUserRepository,
) *UserService {
	return &UserService{
		userrepo: userrepo,
	}
}

func (s *UserService) Create(req requests.CreateUserRequest) (*responses.UserResponse, error) {
	var ures *responses.UserResponse

	u := &models.User{
		Email:    req.Email,
		Password: req.Password,
	}

	u, err := s.userrepo.Create(*u)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	umar, _ := json.Marshal(u)
	err = json.Unmarshal(umar, &ures)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return ures, err
}

func (s *UserService) FindByID(uid uuid.UUID) (*responses.UserResponse, error) {
	u, err := s.userrepo.FindByID(uid)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return u.ToResponse(), err
}

func (s *UserService) FindAll(p utpagination.Pagination) (*utpagination.Pagination, error) {
	users, err := s.userrepo.FindAll(p)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return users, nil
}

func (s *UserService) Delete(uid uuid.UUID) error {
	u := models.User{
		Model: helper.Model{ID: uid},
	}

	err := s.userrepo.Delete(u)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (s *UserService) Update(req requests.UpdateUserRequest, uid uuid.UUID) (*responses.UserResponse, error) {
	u := req.ToModel(consttypes.UR_USER, uid)

	u, err := s.userrepo.Update(*u)
	if err != nil {
		utlogger.LogError(err)
		return nil, err
	}

	return u.ToResponse(), err
}

package user

import "project-skbackend/internal/repositories"

type UserService struct {
	userRepo repositories.IUserRepository
}

func NewUserService(userRepo repositories.IUserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

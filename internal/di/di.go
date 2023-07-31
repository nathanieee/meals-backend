package di

import (
	"project-skbackend/configs"
	userRepo "project-skbackend/internal/repositories/user"
	userService "project-skbackend/internal/services/user"

	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService *userService.UserService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config) *DependencyInjection {

	/* --------------------------- setup user service --------------------------- */

	userRepo := userRepo.NewUserRepo(db)
	userService := userService.NewUserService(userRepo)

	return &DependencyInjection{
		UserService: userService,
	}
}

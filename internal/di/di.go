package di

import (
	"project-skbackend/configs"
	urepo "project-skbackend/internal/repositories/user"
	ausvc "project-skbackend/internal/services/auth"
	usvc "project-skbackend/internal/services/user"

	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService *usvc.UserService
	AuthService *ausvc.AuthService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config) *DependencyInjection {

	/* ------------------------------ USER SERVICE ------------------------------ */

	urepo := urepo.NewUserRepo(db)
	usvc := usvc.NewUserService(urepo)

	/* ------------------------------ AUTH SERVICE ------------------------------ */

	ausvc := ausvc.NewAuthService(urepo, cfg)

	return &DependencyInjection{
		UserService: usvc,
		AuthService: ausvc,
	}
}

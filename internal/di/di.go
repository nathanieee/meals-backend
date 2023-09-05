package di

import (
	"project-skbackend/configs"
	urepo "project-skbackend/internal/repositories/user"
	ausvc "project-skbackend/internal/services/auth"
	msvc "project-skbackend/internal/services/mail"
	usvc "project-skbackend/internal/services/user"

	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService *usvc.UserService
	AuthService *ausvc.AuthService
	MailService *msvc.MailService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config) *DependencyInjection {

	/* ------------------------------ user service ------------------------------ */

	urepo := urepo.NewUserRepo(db)
	usvc := usvc.NewUserService(urepo)

	/* ------------------------------ mail service ------------------------------ */

	msvc := msvc.NewMailService(cfg)

	/* ------------------------------ auth service ------------------------------ */

	ausvc := ausvc.NewAuthService(urepo, cfg, msvc)

	return &DependencyInjection{
		UserService: usvc,
		AuthService: ausvc,
		MailService: msvc,
	}
}

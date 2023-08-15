package di

import (
	"project-skbackend/configs"
	lrepo "project-skbackend/internal/repositories/level"
	rrepo "project-skbackend/internal/repositories/role"
	urepo "project-skbackend/internal/repositories/user"
	ausvc "project-skbackend/internal/services/auth"
	lsvc "project-skbackend/internal/services/level"
	msvc "project-skbackend/internal/services/mail"
	rsvc "project-skbackend/internal/services/role"
	usvc "project-skbackend/internal/services/user"

	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService  *usvc.UserService
	AuthService  *ausvc.AuthService
	RoleService  *rsvc.RoleService
	LevelService *lsvc.LevelService
	MailService  *msvc.MailService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config) *DependencyInjection {

	/* ------------------------------ user service ------------------------------ */

	urepo := urepo.NewUserRepo(db)
	usvc := usvc.NewUserService(urepo)

	/* ------------------------------ mail service ------------------------------ */

	msvc := msvc.NewMailService(cfg)

	/* ------------------------------ auth service ------------------------------ */

	ausvc := ausvc.NewAuthService(urepo, cfg, msvc)

	/* ------------------------------ role service ------------------------------ */

	rrepo := rrepo.NewRoleRepo(db)
	rsvc := rsvc.NewRoleService(rrepo)

	/* ------------------------------ level service ----------------------------- */

	lrepo := lrepo.NewLevelRepo(db)
	lsvc := lsvc.NewLevelService(lrepo)

	return &DependencyInjection{
		UserService:  usvc,
		AuthService:  ausvc,
		RoleService:  rsvc,
		LevelService: lsvc,
		MailService:  msvc,
	}
}

package di

import (
	"project-skbackend/configs"
	lrepo "project-skbackend/internal/repositories/level"
	rrepo "project-skbackend/internal/repositories/role"
	urepo "project-skbackend/internal/repositories/user"
	ausvc "project-skbackend/internal/services/auth"
	lsvc "project-skbackend/internal/services/level"
	rsvc "project-skbackend/internal/services/role"
	usvc "project-skbackend/internal/services/user"

	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService  *usvc.UserService
	AuthService  *ausvc.AuthService
	RoleService  *rsvc.RoleService
	LevelService *lsvc.LevelService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config) *DependencyInjection {

	/* ------------------------------ USER SERVICE ------------------------------ */

	urepo := urepo.NewUserRepo(db)
	usvc := usvc.NewUserService(urepo)

	/* ------------------------------ AUTH SERVICE ------------------------------ */

	ausvc := ausvc.NewAuthService(urepo, cfg)

	/* ------------------------------ ROLE SERVICE ------------------------------ */

	rrepo := rrepo.NewRoleRepo(db)
	rsvc := rsvc.NewRoleService(rrepo)

	/* ------------------------------ LEVEL SERVICE ----------------------------- */

	lrepo := lrepo.NewLevelRepo(db)
	lsvc := lsvc.NewLevelService(lrepo)

	return &DependencyInjection{
		UserService:  usvc,
		AuthService:  ausvc,
		RoleService:  rsvc,
		LevelService: lsvc,
	}
}

package di

import (
	"project-skbackend/configs"
	merepo "project-skbackend/internal/repositories/member"
	urepo "project-skbackend/internal/repositories/user"
	ausvc "project-skbackend/internal/services/auth"
	masvc "project-skbackend/internal/services/mail"
	mesvc "project-skbackend/internal/services/member"
	usvc "project-skbackend/internal/services/user"

	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService   *usvc.UserService
	AuthService   *ausvc.AuthService
	MailService   *masvc.MailService
	MemberService *mesvc.MemberService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config) *DependencyInjection {

	/* ------------------------------ user service ------------------------------ */

	urepo := urepo.NewUserRepo(db)
	usvc := usvc.NewUserService(urepo)

	/* ------------------------------ mail service ------------------------------ */

	masvc := masvc.NewMailService(cfg)

	/* ------------------------------ auth service ------------------------------ */

	ausvc := ausvc.NewAuthService(urepo, cfg, masvc)

	/* ----------------------------- member service ----------------------------- */
	merepo := merepo.NewMemberRepo(db)
	mesvc := mesvc.NewMemberService(merepo, urepo)

	return &DependencyInjection{
		UserService:   usvc,
		AuthService:   ausvc,
		MailService:   masvc,
		MemberService: mesvc,
	}
}

package di

import (
	"project-skbackend/configs"
	allgrepository "project-skbackend/internal/repositories/allergy"
	crgvrrepository "project-skbackend/internal/repositories/caregiver"
	mmbrrepository "project-skbackend/internal/repositories/member"
	userrepository "project-skbackend/internal/repositories/user"
	authservice "project-skbackend/internal/services/auth"
	mailservice "project-skbackend/internal/services/mail"
	mmbrservice "project-skbackend/internal/services/member"
	userservice "project-skbackend/internal/services/user"

	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService   *userservice.UserService
	AuthService   *authservice.AuthService
	MailService   *mailservice.MailService
	MemberService *mmbrservice.MemberService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config) *DependencyInjection {

	/* ---------------------------------- user ---------------------------------- */
	user := userrepository.NewUserRepository(db)
	us := userservice.NewUserService(user)

	/* ---------------------------------- mail ---------------------------------- */
	mais := mailservice.NewMailService(cfg)

	/* ---------------------------------- auth ---------------------------------- */
	auts := authservice.NewAuthService(cfg, user, mais)

	/* ---------------------------- caregiver service --------------------------- */
	carr := crgvrrepository.NewCaregiverRepository(db)

	/* --------------------------------- allergy -------------------------------- */
	allr := allgrepository.NewAllergyRepository(db)

	/* ----------------------------- member service ----------------------------- */
	memr := mmbrrepository.NewMemberRepository(db)
	mems := mmbrservice.NewMemberService(memr, user, carr, allr)

	return &DependencyInjection{
		UserService:   us,
		AuthService:   auts,
		MailService:   mais,
		MemberService: mems,
	}
}

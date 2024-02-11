package di

import (
	"project-skbackend/configs"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/internal/services/mealservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/partnerservice"
	"project-skbackend/internal/services/userservice"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService    *userservice.UserService
	AuthService    *authservice.AuthService
	MailService    *mailservice.MailService
	MemberService  *memberservice.MemberService
	PartnerService *partnerservice.PartnerService
	MealService    *mealservice.MealService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config, rdb *redis.Client) *DependencyInjection {
	/* -------------------------------- database -------------------------------- */
	if cfg.DB.LogMode {
		db = db.Debug()
	}

	/* ---------------------------------- user ---------------------------------- */
	ruser := userrepo.NewUserRepository(db)
	suser := userservice.NewUserService(ruser)

	/* --------------------------------- partner -------------------------------- */
	rpartner := partnerrepo.NewPartnerRepository(db)
	spartner := partnerservice.NewPartnerService(rpartner)

	/* ---------------------------------- mail ---------------------------------- */
	smail := mailservice.NewMailService(cfg)

	/* ---------------------------------- auth ---------------------------------- */
	sauth := authservice.NewAuthService(cfg, ruser, smail, rdb)

	/* -------------------------------- caregiver ------------------------------- */
	rcaregiver := caregiverrepo.NewCaregiverRepository(db)

	/* --------------------------------- allergy -------------------------------- */
	rallergy := allergyrepo.NewAllergyRepository(db)

	/* --------------------------------- illness -------------------------------- */
	rillness := illnessrepo.NewIllnessRepository(db)

	/* ---------------------------------- meal ---------------------------------- */
	rmeal := mealrepo.NewMealRepository(db)
	smeal := mealservice.NewMealService(rmeal, rillness, rallergy, rpartner)

	/* ------------------------------ organization ------------------------------ */
	rorganization := organizationrepo.NewOrganizationRepository(db)

	/* --------------------------------- member --------------------------------- */
	rmember := memberrepo.NewMemberRepository(db)
	smember := memberservice.NewMemberService(rmember, ruser, rcaregiver, rallergy, rillness, *rorganization)

	return &DependencyInjection{
		UserService:    suser,
		AuthService:    sauth,
		MailService:    smail,
		MemberService:  smember,
		PartnerService: spartner,
		MealService:    smeal,
	}
}

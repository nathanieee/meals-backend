package di

import (
	"project-skbackend/configs"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/cartrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/cartservice"
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
	CartService    *cartservice.CartService
}

func NewDependencyInjection(db *gorm.DB, cfg *configs.Config, rdb *redis.Client) *DependencyInjection {
	/* -------------------------------- database -------------------------------- */
	db = db.Session(&gorm.Session{FullSaveAssociations: true})

	if cfg.DB.LogMode {
		db = db.Debug()
	}

	/* ------------------------------- repository ------------------------------- */
	ruser := userrepo.NewUserRepository(db)
	rpartner := partnerrepo.NewPartnerRepository(db)
	rcaregiver := caregiverrepo.NewCaregiverRepository(db)
	rallergy := allergyrepo.NewAllergyRepository(db)
	rmeal := mealrepo.NewMealRepository(db)
	rillness := illnessrepo.NewIllnessRepository(db)
	rmember := memberrepo.NewMemberRepository(db)
	rorganization := organizationrepo.NewOrganizationRepository(db)
	rcart := cartrepo.NewCartRepository(db)

	/* --------------------------------- service -------------------------------- */
	suser := userservice.NewUserService(ruser)
	spartner := partnerservice.NewPartnerService(rpartner)
	smail := mailservice.NewMailService(cfg)
	sauth := authservice.NewAuthService(cfg, ruser, smail, rdb)
	smeal := mealservice.NewMealService(rmeal, rillness, rallergy, rpartner)
	smember := memberservice.NewMemberService(rmember, ruser, rcaregiver, rallergy, rillness, *rorganization)
	scart := cartservice.NewCartService(rcart, rcaregiver, rmember)

	return &DependencyInjection{
		UserService:    suser,
		AuthService:    sauth,
		MailService:    smail,
		MemberService:  smember,
		PartnerService: spartner,
		MealService:    smeal,
		CartService:    scart,
	}
}

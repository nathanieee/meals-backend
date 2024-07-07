package di

import (
	"context"
	"project-skbackend/configs"
	"project-skbackend/internal/repositories/adminrepo"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/cartrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/orderrepo"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/internal/repositories/patronrepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/internal/services/consumerservice"
	"project-skbackend/internal/services/cronservice"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/internal/services/mealservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/orderservice"
	"project-skbackend/internal/services/organizationservice"
	"project-skbackend/internal/services/partnerservice"
	"project-skbackend/internal/services/patronservice"
	"project-skbackend/internal/services/producerservice"
	"project-skbackend/internal/services/userservice"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DependencyInjection struct {
	UserService         *userservice.UserService
	AuthService         *authservice.AuthService
	MailService         *mailservice.MailService
	MemberService       *memberservice.MemberService
	PartnerService      *partnerservice.PartnerService
	MealService         *mealservice.MealService
	CartService         *cartservice.CartService
	ConsumerService     *consumerservice.ConsumerService
	PatronService       *patronservice.PatronService
	OrganizationService *organizationservice.OrganizationService
	OrderService        *orderservice.OrderService
	CronService         *cronservice.CronService
}

func NewDependencyInjection(db *gorm.DB, ch *amqp.Channel, cfg *configs.Config, rdb *redis.Client, ctx context.Context) *DependencyInjection {
	/* -------------------------------- database -------------------------------- */
	if cfg.DB.LogMode {
		db = db.Debug()
	}

	/* ------------------------------- repository ------------------------------- */
	ruser := userrepo.NewUserRepository(db)
	rpart := partnerrepo.NewPartnerRepository(db)
	rcare := caregiverrepo.NewCaregiverRepository(db)
	rall := allergyrepo.NewAllergyRepository(db)
	rmeal := mealrepo.NewMealRepository(db)
	rill := illnessrepo.NewIllnessRepository(db)
	rmemb := memberrepo.NewMemberRepository(db)
	rorg := organizationrepo.NewOrganizationRepository(db)
	rcart := cartrepo.NewCartRepository(db)
	radmin := adminrepo.NewAdminRepository(db)
	rpatron := patronrepo.NewPatronRepository(db)
	rorder := orderrepo.NewOrderRepository(db, *cfg)

	/* --------------------------------- service -------------------------------- */
	sprod := producerservice.NewProducerService(ch, cfg, ctx)
	suser := userservice.NewUserService(ruser, radmin, rcare, rmemb, rorg, rpart)
	spart := partnerservice.NewPartnerService(rpart)
	smail := mailservice.NewMailService(cfg, ruser, sprod)
	sauth := authservice.NewAuthService(cfg, rdb, ruser, smail, suser)
	smeal := mealservice.NewMealService(rmeal, rill, rall, rpart)
	smemb := memberservice.NewMemberService(rmemb, ruser, rcare, rall, rill, *rorg)
	scart := cartservice.NewCartService(rcart, rcare, rmemb)
	scons := consumerservice.NewConsumerService(ch, cfg, smail)
	spatr := patronservice.NewPatronService(rpatron)
	sorga := organizationservice.NewOrganizationService(rorg)
	sordr := orderservice.NewOrderService(*cfg, rorder, rmeal, rmemb, ruser, rcare)
	scron := cronservice.NewCronService(cfg, rorder)

	return &DependencyInjection{
		UserService:         suser,
		AuthService:         sauth,
		MailService:         smail,
		MemberService:       smemb,
		PartnerService:      spart,
		MealService:         smeal,
		CartService:         scart,
		ConsumerService:     scons,
		PatronService:       spatr,
		OrganizationService: sorga,
		OrderService:        sordr,
		CronService:         scron,
	}
}

func (di *DependencyInjection) InitServices() {
	// * setup consumer service
	di.ConsumerService.ConsumeTask()

	// * init cron service
	di.CronService.Init()
}

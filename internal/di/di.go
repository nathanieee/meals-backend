package di

import (
	"context"
	"project-skbackend/configs"
	"project-skbackend/external/services/distancematrixservice"
	"project-skbackend/internal/repositories/adminrepo"
	"project-skbackend/internal/repositories/allergyrepo"
	"project-skbackend/internal/repositories/caregiverrepo"
	"project-skbackend/internal/repositories/cartrepo"
	"project-skbackend/internal/repositories/donationproofrepo"
	"project-skbackend/internal/repositories/donationrepo"
	"project-skbackend/internal/repositories/illnessrepo"
	"project-skbackend/internal/repositories/imagerepo"
	"project-skbackend/internal/repositories/mealrepo"
	"project-skbackend/internal/repositories/memberrepo"
	"project-skbackend/internal/repositories/orderrepo"
	"project-skbackend/internal/repositories/organizationrepo"
	"project-skbackend/internal/repositories/partnerrepo"
	"project-skbackend/internal/repositories/patronrepo"
	"project-skbackend/internal/repositories/userimagerepo"
	"project-skbackend/internal/repositories/userrepo"
	"project-skbackend/internal/services/allergyservice"
	"project-skbackend/internal/services/authservice"
	"project-skbackend/internal/services/cartservice"
	"project-skbackend/internal/services/consumerservice"
	"project-skbackend/internal/services/cronservice"
	"project-skbackend/internal/services/fileservice"
	"project-skbackend/internal/services/illnessservice"
	"project-skbackend/internal/services/mailservice"
	"project-skbackend/internal/services/mealservice"
	"project-skbackend/internal/services/memberservice"
	"project-skbackend/internal/services/orderservice"
	"project-skbackend/internal/services/organizationservice"
	"project-skbackend/internal/services/partnerservice"
	"project-skbackend/internal/services/patronservice"
	"project-skbackend/internal/services/producerservice"
	"project-skbackend/internal/services/userservice"
	"project-skbackend/packages/utils/utlogger"

	"github.com/minio/minio-go/v7"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DependencyInjection struct {
	// * internal services
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
	IllnessService      *illnessservice.IllnessService
	FileService         *fileservice.FileService
	AllergyService      *allergyservice.AllergyService

	// * external services
	DistanceMatrixService *distancematrixservice.DistanceMatrixService
}

func NewDependencyInjection(ctx context.Context, db *gorm.DB, ch *amqp.Channel, cfg *configs.Config, rdb *redis.Client, minio *minio.Client) *DependencyInjection {
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
	rimg := imagerepo.NewImageRepository(db)
	ruimg := userimagerepo.NewUserImageRepository(db)
	rordr := orderrepo.NewOrderRepository(db, *cfg)
	rdona := donationrepo.NewDonationRepository(db)
	rdnpr := donationproofrepo.NewDonationProofRepository(db)

	/* --------------------------------- service -------------------------------- */
	// * external services
	sdsmx := distancematrixservice.NewDistanceMatrixService(cfg)

	// * internal services
	sprod := producerservice.NewProducerService(ch, cfg, ctx)
	suser := userservice.NewUserService(ruser, radmin, rcare, rmemb, rorg, rpart, rpatron)
	spart := partnerservice.NewPartnerService(rpart, rordr)
	smail := mailservice.NewMailService(cfg, ruser, sprod)
	sauth := authservice.NewAuthService(cfg, rdb, ruser, smail, suser)
	smeal := mealservice.NewMealService(rmeal, rill, rall, rpart)
	smemb := memberservice.NewMemberService(rmemb, ruser, rcare, rall, rill, rorg)
	scart := cartservice.NewCartService(rcart, rcare, rmemb)
	scons := consumerservice.NewConsumerService(ch, cfg, smail)
	spatr := patronservice.NewPatronService(rpatron, rdona)
	sorga := organizationservice.NewOrganizationService(rorg)
	sordr := orderservice.NewOrderService(cfg, rorder, rmeal, rmemb, ruser, rcare, rcart)
	scron := cronservice.NewCronService(cfg, rorder)
	silln := illnessservice.NewIllnessService(rill)
	sfile := fileservice.NewFileService(cfg, ctx, *minio, ruser, rimg, ruimg, rdona, rdnpr)
	salle := allergyservice.NewAllergyService(rall)

	return &DependencyInjection{
		// * internal services
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
		IllnessService:      silln,
		FileService:         sfile,
		AllergyService:      salle,

		// * external services
		DistanceMatrixService: sdsmx,
	}
}

func (di *DependencyInjection) InitServices() {
	var (
		err error
	)

	// * setup consumer service
	di.ConsumerService.ConsumeTask()

	// * init cron service
	_, err = di.CronService.Init()
	utlogger.Fatal(err)
}

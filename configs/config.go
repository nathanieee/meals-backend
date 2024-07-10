package configs

import (
	"fmt"
	"project-skbackend/packages/utils/utlogger"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/redis/go-redis/v9"
)

type (
	Config struct {
		// * internal config
		API
		App
		Web
		File
		HTTP
		DB
		Mail
		JWT
		Credential
		Order

		// * external config
		Redis
		Xendit
		AWS
		Localstack
		Queue
	}

	API struct {
		VerifyTokenLength int    `env:"API_VERIFY_TOKEN_LENGTH" env-default:"8"`
		URL               string `env:"API_URL" env-default:"localhost"`
		Timezone          string `env:"API_TIMEZONE" env-default:"Asia/Makassar"`
		APIResetPassword
	}
	APIResetPassword struct {
		Cooldown int `env:"API_RESET_PASSWORD_COOLDOWN" env-default:"5"`
	}

	Order struct {
		OrderBuffer
		OrderMax
	}
	OrderBuffer struct {
		AutomaticallyCancelled      int `env:"ORDER_AUTOMATICALLY_CANCELLED_BUFFER" env-default:"10"`
		AutomaticallyBeingPickedUp  int `env:"ORDER_AUTOMATICALLY_BEING_PICKED_UP" env-default:"10"`
		AutomaticallyOutForDelivery int `env:"ORDER_AUTOMATICALLY_OUT_FOR_DELIVERY" env-default:"10"`
		AutomaticallyDelivered      int `env:"ORDER_AUTOMATICALLY_DELIVERED" env-default:"10"`
	}
	OrderMax struct {
		Member uint `env:"ORDER_MAX_MEMBER" env-default:"3"`
	}

	App struct {
		Name        string `env:"APP_NAME" env-default:"meals-app"`
		Version     string `env:"APP_VERSION" env-default:"1.0"`
		Url         string `env:"APP_URL"`
		Env         string `env:"APP_ENV" env-default:"development"`
		Timeout     int    `env:"APP_TIMEOUT" env-default:"30"`
		DeeplinkUrl string `env:"DEEPLINK_URL"`
	}

	Web struct {
		URL string `env:"WEB_URL"`
	}

	File struct {
		FileImage
	}
	FileImage struct {
		BaseDir    string `env:"IMAGE_BASE_DIR" env-default:"../assets/images"`
		ProfileDir string `env:"IMAGE_PROFILE_DIR" env-default:"/profile"`
		MealDir    string `env:"IMAGE_MEAL_DIR" env-default:"/meal"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" env-default:"8080"`
	}

	DB struct {
		PoolMax  int    `env:"DB_POOL_MAX" env-default:"10"`
		Name     string `env:"DB_NAME" env-default:"meals-pg"`
		User     string `env:"DB_USER" env-default:"root"`
		Password string `env:"DB_PASSWORD" env-default:"password"`
		Host     string `env:"DB_HOST" env-default:"localhost"`
		Port     string `env:"DB_PORT" env-default:"5432"`
		LogMode  bool   `env:"DB_LOG_MODE" env-default:"false"`
		SslMode  string `env:"DB_SSL_MODE" env-default:"disable"`
		Timezone string `env:"DB_TIMEZONE" env-default:"Asia/Makassar"`
	}

	Mail struct {
		Name        string `env:"MAIL_NAME"`
		From        string `env:"MAIL_FROM"`
		Password    string `env:"MAIL_PASSWORD"`
		TemplateDir string `env:"MAIL_TEMPLATE_DIR" env-default:"../web/templates"`
		SMTPHost    string `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
		SMTPPort    int    `env:"SMTP_PORT" env-default:"587"`
	}

	JWT struct {
		TimeUnit string `env:"JWT_TIME_UNIT" env-default:"hours"`
		JWTAccessToken
		JWTRefreshToken
	}
	JWTAccessToken struct {
		PublicKey  string `env:"ACCESS_TOKEN_PUBLIC_KEY"`
		PrivateKey string `env:"ACCESS_TOKEN_PRIVATE_KEY"`
		Life       int    `env:"ACCESS_TOKEN_LIFE" env-default:"3600"`
	}
	JWTRefreshToken struct {
		PublicKey  string `env:"REFRESH_TOKEN_PUBLIC_KEY"`
		PrivateKey string `env:"REFRESH_TOKEN_PRIVATE_KEY"`
		Life       int    `env:"REFRESH_TOKEN_LIFE" env-default:"86400"`
	}

	Credential struct {
		CredentialAdmin
	}
	CredentialAdmin struct {
		Email     string `env:"ADMIN_EMAIL"`
		Password  string `env:"ADMIN_PASSWORD"`
		FirstName string `env:"ADMIN_FIRST_NAME"`
		LastName  string `env:"ADMIN_LAST_NAME"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST"`
		Port     string `env:"REDIS_PORT"`
		Password string `env:"REDIS_PASSWORD"`
	}

	Xendit struct {
		SecretKey    string `env:"XEN_SECRET_KEY"`
		WebhookToken string `env:"XEN_WEBHOOK_TOKEN"`
	}

	AWS struct {
		AWSAccessKey
		Region string `env:"AWS_REGION" env-default:"ap-southeast-1"`
	}
	AWSAccessKey struct {
		PublicKey string `env:"AWS_PUBLIC_ACCESS_KEY"`
		SecretKey string `env:"AWS_SECRET_ACCESS_KEY"`
	}

	Localstack struct {
		Port  string `env:"LOCALSTACK_PORT" env-default:"4566"`
		Debug int    `env:"LOCALSTACK_DEBUG" env-default:"0"`
	}

	Queue struct {
		Host     string `env:"RABBIT_MQ_HOST"`
		Port     string `env:"RABBIT_MQ_PORT"`
		Username string `env:"RABBIT_MQ_USERNAME"`
		Password string `env:"RABBIT_MQ_PASSWORD"`
		QueueMail
	}
	QueueMail struct {
		QueueName    string `env:"MAIL_QUEUE_NAME"`
		ExchangeName string `env:"MAIL_EXCHANGE_NAME"`
		ExchangeType string `env:"MAIL_EXCHANGE_TYPE"`
		BindingKey   string `env:"MAIL_BINDING_KEY"`
	}
)

var (
	once     sync.Once
	instance *Config
	rdb      *redis.Client
)

func GetInstance() *Config {
	if instance == nil {
		once.Do(func() {
			cfg, err := newConfig()
			if err != nil {
				utlogger.Error(err)
			}

			instance = cfg
			rdb = instance.GetRedisClient()
		})
	}

	return instance
}

func newConfig() (*Config, error) {
	cfg := &Config{}
	err := cleanenv.ReadConfig(".env", cfg)
	if err != nil {
		fmt.Println("using environment variable")
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		utlogger.Error(err)
		return nil, err
	}

	return cfg, nil
}

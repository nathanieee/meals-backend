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
		File
		HTTP
		DB
		Mail
		JWT
		Credential

		// * external config
		Redis
		Xendit
		AWS
	}

	API struct {
		VerifyTokenLength int    `env:"API_VERIFY_TOKEN_LENGTH" env-default:"8"`
		Domain            string `env:"API_DOMAIN" env-default:"localhost"`
	}

	App struct {
		Name        string `env:"APP_NAME" env-default:"meals-app"`
		Version     string `env:"APP_VERSION" env-default:"1.0"`
		Url         string `env:"APP_URL"`
		Env         string `env:"APP_ENV" env-default:"development"`
		Timeout     int    `env:"APP_TIMEOUT" env-default:"30"`
		DeeplinkUrl string `env:"DEEPLINK_URL"`
	}

	File struct {
		Image struct {
			BaseDir    string `env:"IMAGE_BASE_DIR" env-default:"../assets/images"`
			ProfileDir string `env:"IMAGE_PROFILE_DIR" env-default:"/profile"`
			MealDir    string `env:"IMAGE_MEAL_DIR" env-default:"/meal"`
		}
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
	}

	Mail struct {
		From        string `env:"MAIL_FROM"`
		Password    string `env:"MAIL_PASSWORD"`
		TemplateDir string `env:"MAIL_TEMPLATE_DIR" env-default:"../web/templates"`
		SMTPHost    string `env:"SMTP_HOST" env-default:"smtp.gmail.com"`
		SMTPPort    string `env:"SMTP_PORT" env-default:"587"`
	}

	JWT struct {
		TimeUnit    string `env:"JWT_TIME_UNIT" env-default:"hours"`
		AccessToken struct {
			PublicKey  string `env:"ACCESS_TOKEN_PUBLIC_KEY"`
			PrivateKey string `env:"ACCESS_TOKEN_PRIVATE_KEY"`
			Life       int    `env:"ACCESS_TOKEN_LIFE" env-default:"3600"`
		}
		RefreshToken struct {
			PublicKey  string `env:"REFRESH_TOKEN_PUBLIC_KEY"`
			PrivateKey string `env:"REFRESH_TOKEN_PRIVATE_KEY"`
			Life       int    `env:"REFRESH_TOKEN_LIFE" env-default:"86400"`
		}
	}

	Credential struct {
		Admin struct {
			Email     string `env:"ADMIN_EMAIL"`
			Password  string `env:"ADMIN_PASSWORD"`
			FirstName string `env:"ADMIN_FIRST_NAME"`
			LastName  string `env:"ADMIN_LAST_NAME"`
		}
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
		AccessKey struct {
			PublicKey string `env:"AWS_PUBLIC_ACCESS_KEY"`
			SecretKey string `env:"AWS_SECRET_ACCESS_KEY"`
		}
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
				utlogger.LogError(err)
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
		utlogger.LogError(err)
		return nil, err
	}

	return cfg, nil
}

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
		API
		App
		HTTP
		DB
		Redis
		Mail
		JWT
		Credential
		Xendit
	}

	API struct {
		VerifyTokenLength int    `env:"API_VERIFY_TOKEN_LENGTH" default:"8"`
		Domain            string `env:"API_DOMAIN"`
	}

	App struct {
		Name        string `env:"APP_NAME" default:"MyApp"`
		Version     string `env:"APP_VERSION" default:"1.0"`
		Url         string `env:"APP_URL"`
		Env         string `env:"APP_ENV" default:"development"`
		Timeout     int    `env:"APP_TIMEOUT" default:"30"`
		DeeplinkUrl string `env:"DEEPLINK_URL"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT" default:"8080"`
	}

	DB struct {
		PoolMax  int    `env:"DB_POOL_MAX" default:"10"`
		Host     string `env:"DB_HOST"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_NAME"`
		Port     string `env:"DB_PORT" default:"5432"`
		SslMode  string `env:"DB_SSL_MODE" default:"disable"`
		LogMode  bool   `env:"DB_LOG_MODE" default:"false"`
	}

	Redis struct {
		Host     string `env:"REDIS_HOST"`
		Port     string `env:"REDIS_PORT"`
		Password string `env:"REDIS_PASSWORD"`
	}

	Mail struct {
		From              string `env:"MAIL_FROM"`
		Password          string `env:"MAIL_PASSWORD"`
		SMTPHost          string `env:"SMTP_HOST"`
		SMTPPort          string `env:"SMTP_PORT"`
		TemplateDirectory string `env:"TEMPLATE_DIRECTORY"`
	}

	JWT struct {
		TimeUnit    string `env:"JWT_TIME_UNIT" default:"hours"`
		AccessToken struct {
			PublicKey  string `env:"ACCESS_TOKEN_PUBLIC_KEY"`
			PrivateKey string `env:"ACCESS_TOKEN_PRIVATE_KEY"`
			Life       int    `env:"ACCESS_TOKEN_LIFE" default:"3600"`
		}
		RefreshToken struct {
			PublicKey  string `env:"REFRESH_TOKEN_PUBLIC_KEY"`
			PrivateKey string `env:"REFRESH_TOKEN_PRIVATE_KEY"`
			Life       int    `env:"REFRESH_TOKEN_LIFE" default:"86400"`
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

	Xendit struct {
		SecretKey    string `env:"XEN_SECRET_KEY"`
		WebhookToken string `env:"XEN_WEBHOOK_TOKEN"`
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

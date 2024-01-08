package configs

import (
	"fmt"
	"project-skbackend/packages/utils/utlogger"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App
		HTTP
		DB
		Mail
	}

	App struct {
		Name                  string `env:"APP_NAME"`
		Version               string `env:"APP_VERSION"`
		Url                   string `env:"APP_URL"`
		Secret                string `env:"APP_SECRET"`
		TokenLifespanDuration string `env:"TOKEN_DURATION"`
		TokenLifespan         int    `env:"TOKEN_LIFESPAN"`
		RefreshTokenLifespan  int    `env:"REFRESH_TOKEN_LIFESPAN"`
		DeeplinkUrl           string `env:"DEEPLINK_URL"`
		Timeout               int    `env:"APP_TIMEOUT"`
		Env                   string `env:"APP_ENV"`
	}

	HTTP struct {
		Port string `env:"HTTP_PORT"`
	}

	DB struct {
		PoolMax  int    `env:"DB_POOL_MAX"`
		Host     string `env:"DB_HOST"`
		User     string `env:"DB_USER"`
		Password string `env:"DB_PASSWORD"`
		Name     string `env:"DB_NAME"`
		Port     string `env:"DB_PORT"`
		SslMode  string `env:"DB_SSL_MODE"`
		LogMode  bool   `env:"DB_LOG_MODE"`
	}

	Mail struct {
		From              string `env:"MAIL_FROM"`
		Password          string `env:"MAIL_PASSWORD"`
		SMTPHost          string `env:"SMTP_HOST"`
		SMTPPort          string `env:"SMTP_PORT"`
		TemplateDirectory string `env:"TEMPLATE_DIRECTORY"`
	}
)

var (
	once     sync.Once
	instance *Config
)

func GetInstance() *Config {
	if instance == nil {
		once.Do(func() {
			cfg, err := newConfig()

			if err != nil {
				utlogger.LogError(err)
			}
			instance = cfg
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

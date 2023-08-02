package configs

import (
	"fmt"
	"log"
	"project-skbackend/internal/models"
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"gorm.io/gorm"
)

type (
	Config struct {
		App
		HTTP
		DB
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
		PoolMax      int    `env:"MASTER_DB_POOL_MAX"`
		Host         string `env:"MASTER_DB_HOST"`
		User         string `env:"MASTER_DB_USER"`
		Password     string `env:"MASTER_DB_PASSWORD"`
		DatabaseName string `env:"MASTER_DB_NAME"`
		Port         string `env:"MASTER_DB_PORT"`
		SslMode      string `env:"MASTER_SSL_MODE"`
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
				log.Fatal(err)
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
		fmt.Println("Using Environment Variable")
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func (db DB) GetDbConnectionUrl() string {
	connectionUrl := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		db.Host, db.User, db.Password, db.DatabaseName, db.Port, db.SslMode,
	)
	return connectionUrl
}

func (db DB) AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&models.User{},
		&models.Level{},
		&models.Role{},
	)
}

func (db DB) AutoSeed(gdb *gorm.DB) error {
	err := SeedAdminCredentials(gdb)

	return err
}

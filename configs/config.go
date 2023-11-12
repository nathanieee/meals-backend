package configs

import (
	"errors"
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
		PoolMax      int    `env:"DB_POOL_MAX"`
		Host         string `env:"DB_HOST"`
		User         string `env:"DB_USER"`
		Password     string `env:"DB_PASSWORD"`
		DatabaseName string `env:"DB_NAME"`
		Port         string `env:"DB_PORT"`
		SslMode      string `env:"SSL_MODE"`
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

func (db DB) DBSetup(gdb *gorm.DB) error {
	err := db.AutoSeedEnum(gdb)
	if err != nil {
		return err
	}

	err = db.AutoMigrate(gdb)
	if err != nil {
		return err
	}

	err = db.AutoSeedTable(gdb)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) AutoSeedEnum(gdb *gorm.DB) error {
	err := SeedEnum(gdb)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) AutoMigrate(gdb *gorm.DB) error {
	return gdb.AutoMigrate(
		&models.User{},
		&models.UserImage{},
		&models.Address{},
		&models.Admin{},
		&models.Allergy{},
		&models.Caregiver{},
		&models.Donation{},
		&models.FoodCategory{},
		&models.Illness{},
		&models.Image{},
		&models.Meal{},
		&models.MealAllergy{},
		&models.MealIllness{},
		&models.MealImage{},
		&models.Member{},
		&models.MemberAllergy{},
		&models.MemberIllness{},
		&models.Organization{},
		&models.Partner{},
		&models.Patron{},
		&models.Rating{},
	)
}

func (db DB) AutoSeedTable(gdb *gorm.DB) error {
	seedfuncs := []func(*gorm.DB) error{
		SeedAdminCredentials,
		SeedAllergyData,
	}

	var errs []error

	for _, seedfunc := range seedfuncs {
		if err := seedfunc(gdb); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.New("Error seeding table")
	}

	return nil
}

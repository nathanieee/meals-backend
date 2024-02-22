package configs

import (
	"errors"
	"fmt"
	"project-skbackend/internal/models"
	"project-skbackend/packages/utils/utlogger"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func (db DB) GetDbConnectionUrl() string {
	url := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		db.Host, db.User, db.Password, db.Name, db.Port, db.SslMode,
	)
	return url
}

func (db DB) DBSetup(gdb *gorm.DB) error {
	err := db.AutoSeedEnum(gdb)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	err = db.AutoMigrate(gdb)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	err = db.AutoSeedTable(gdb)
	if err != nil {
		utlogger.LogError(err)
		return err
	}

	return nil
}

func (db DB) AutoSeedEnum(gdb *gorm.DB) error {
	seedfuncs := []func(*gorm.DB) error{
		SeedAllergensEnum,
		SeedGenderEnum,
		SeedMealStatusEnum,
		SeedDonationStatusEnum,
		SeedImageTypeEnum,
		SeedPatronTypeEnum,
		SeedOrganizationTypeEnum,
	}

	var errs []error

	for _, seedfunc := range seedfuncs {
		if err := seedfunc(gdb); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		utlogger.LogError(errs...)
		return errors.New("error seeding enum")
	}

	return nil
}

func (db DB) GetLogLevel() logger.LogLevel {
	loglevel := logger.Warn
	if db.LogMode {
		loglevel = logger.Info
	}

	return loglevel
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
		SeedIllnessData,
		SeedOrganizationCredential,
		SeedPartnerCredential,
		SeedMealData,
	}

	var errs []error

	for _, seedfunc := range seedfuncs {
		if err := seedfunc(gdb); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		utlogger.LogError(errs...)
		return errors.New("error seeding table")
	}

	return nil
}

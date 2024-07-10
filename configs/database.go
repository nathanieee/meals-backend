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
	if err := db.AutoSeedEnum(gdb); err != nil {
		utlogger.Error(err)
		return err
	}

	if err := db.AutoMigrate(gdb); err != nil {
		utlogger.Error(err)
		return err
	}

	if err := db.AutoSeedData(gdb); err != nil {
		utlogger.Error(err)
		return err
	}

	return nil
}

func (db DB) AutoSeedEnum(gdb *gorm.DB) error {
	seedfuncs := []func(*gorm.DB) error{
		/* ---------------------------------- enum ---------------------------------- */
		SeedAllergensEnum,
		SeedGenderEnum,
		SeedMealStatusEnum,
		SeedDonationStatusEnum,
		SeedImageTypeEnum,
		SeedPatronTypeEnum,
		SeedOrganizationTypeEnum,
		SeedUserRoleEnum,
		SeedOrderStatusEnum,
	}

	var (
		errs []error
	)

	for _, seedfunc := range seedfuncs {
		if err := seedfunc(gdb); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		utlogger.Error(errs...)
		return errors.New("error seeding enum")
	}

	return nil
}

func (db DB) AutoSeedData(gdb *gorm.DB) error {
	seedfuncs := []func(*gorm.DB) error{
		// * dependent
		SeedMealCategoryData,

		/* ------------------------------- credentials ------------------------------ */
		SeedAdminCredentials,
		SeedOrganizationCredentials,
		SeedPartnerCredentials,
		SeedMemberCredentials,

		/* ---------------------------------- data ---------------------------------- */
		// * independent
		SeedAllergyData,
		SeedIllnessData,

		// * dependent
		SeedMealData,
		SeedCartData,
	}

	var (
		errs []error
	)

	for _, seedfunc := range seedfuncs {
		if err := seedfunc(gdb); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		utlogger.Error(errs...)
		return errors.New("error seeding data")
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
		&models.Patron{},
		&models.Donation{},
		&models.Illness{},
		&models.Image{},
		&models.Meal{},
		&models.MealAllergy{},
		&models.MealIllness{},
		&models.MealImage{},
		&models.MealCategory{},
		&models.Member{},
		&models.MemberAllergy{},
		&models.MemberIllness{},
		&models.Organization{},
		&models.Partner{},
		&models.Rating{},
		&models.Cart{},
		&models.Order{},
		&models.OrderHistory{},
		&models.OrderMeal{},
	)
}

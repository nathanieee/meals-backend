package configs

import (
	"errors"
	"fmt"
	"project-skbackend/internal/models"

	"gorm.io/gorm"
)

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
		SeedIllnessData,
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

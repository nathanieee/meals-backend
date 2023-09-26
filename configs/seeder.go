package configs

import (
	"errors"
	"os"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/admin"
	"project-skbackend/packages/consttypes"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func SeedAdminCredentials(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			gotime, _ := time.Parse(consttypes.DATEFORMAT, "2000-10-20")
			date := datatypes.Date(gotime)

			arepo := admin.NewAdminRepo(db)

			admins := []models.Admin{
				{
					User: models.User{
						Email:    os.Getenv("ADMIN_EMAIL"),
						Password: os.Getenv("ADMIN_PASSWORD"),
						Role:     consttypes.UR_ADMINISTRATOR,
					},
					FirstName:   "John",
					LastName:    "Doe",
					Gender:      consttypes.G_MALE,
					DateOfBirth: date,
				},
			}

			for _, admin := range admins {
				_, err := arepo.Create(&admin)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

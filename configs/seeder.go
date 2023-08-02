package configs

import (
	"errors"
	"os"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/user"
	"project-skbackend/packages/consttypes"

	"gorm.io/gorm"
)

func SeedUserRole(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.Level{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			// 	userRepository := user.NewUserRepo(db)

			// 	admins := []models.User{
			// 		{
			// 			EmailAddress: os.Getenv("ADMIN_EMAIL"),
			// 			Password:     os.Getenv("ADMIN_PASSWORD"),
			// 			FullName:     os.Getenv("ADMIN_FULLNAME"),
			// 			RoleID:       consttypes.ADMIN,
			// 		},
			// 	}

			// 	for _, admin := range admins {
			// 		userRepository.Store(&admin)
			// 	}
		}
	}

	return nil
}

func SeedAdminCredentials(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			userRepository := user.NewUserRepo(db)

			admins := []models.User{
				{
					EmailAddress: os.Getenv("ADMIN_EMAIL"),
					Password:     os.Getenv("ADMIN_PASSWORD"),
					FullName:     os.Getenv("ADMIN_FULLNAME"),
					RoleID:       consttypes.ADMIN,
				},
			}

			for _, admin := range admins {
				userRepository.Store(&admin)
			}
		}
	}

	return nil
}

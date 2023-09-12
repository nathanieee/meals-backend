package configs

import (
	"errors"
	"os"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/user"
	"project-skbackend/packages/consttypes"
	"project-skbackend/packages/utils"

	"gorm.io/gorm"
)

func SeedAdminCredentials(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			urepo := user.NewUserRepo(db)

			admins := []models.User{
				{
					Email:    os.Getenv("ADMIN_EMAIL"),
					Password: os.Getenv("ADMIN_PASSWORD"),
					Role:     consttypes.UR_ADMINISTRATOR,
				},
			}

			for _, admin := range admins {
				hashedPassword, err := utils.EncryptPassword(admin.Password)
				if err != nil {
					return err
				}
				admin.Password = hashedPassword

				urepo.Store(&admin)
			}
		}
	}

	return nil
}

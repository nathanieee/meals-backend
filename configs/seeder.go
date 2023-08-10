package configs

import (
	"errors"
	"os"
	"project-skbackend/internal/models"
	"project-skbackend/internal/repositories/level"
	"project-skbackend/internal/repositories/role"
	"project-skbackend/internal/repositories/user"
	"project-skbackend/packages/consttypes"

	"gorm.io/gorm"
)

func SeedUserLevel(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.Level{}) {
		if err := db.First(&models.Level{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			lrepo := level.NewLevelRepo(db)

			levels := []models.Level{
				{
					Model: gorm.Model{ID: uint(consttypes.USER)},
					Name:  "User",
				},
				{
					Model: gorm.Model{ID: uint(consttypes.ADMIN)},
					Name:  "Admin",
				},
			}

			for _, level := range levels {
				lrepo.Store(&level)
			}
		}
	}

	return nil
}

func SeedUserRole(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.Role{}) {
		if err := db.First(&models.Role{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			rrepo := role.NewRoleRepo(db)

			roles := []models.Role{
				{
					Name:    consttypes.ADMINISTRATOR,
					LevelID: consttypes.ADMIN,
				},
				{
					Name:    consttypes.CAREGIVER,
					LevelID: consttypes.USER,
				},
				{
					Name:    consttypes.MEMBER,
					LevelID: consttypes.USER,
				},
				{
					Name:    consttypes.PARTNER,
					LevelID: consttypes.USER,
				},
				{
					Name:    consttypes.PATRON,
					LevelID: consttypes.USER,
				},
			}

			for _, role := range roles {
				rrepo.Store(&role)
			}
		}
	}

	return nil
}

func SeedAdminCredentials(db *gorm.DB) error {
	if db.Migrator().HasTable(&models.User{}) {
		if err := db.First(&models.User{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			urepo := user.NewUserRepo(db)

			admins := []models.User{
				{
					Email:    os.Getenv("ADMIN_EMAIL"),
					Password: os.Getenv("ADMIN_PASSWORD"),
					FullName: os.Getenv("ADMIN_FULLNAME"),
					RoleID:   uint(consttypes.ADMIN),
				},
			}

			for _, admin := range admins {
				urepo.Store(&admin)
			}
		}
	}

	return nil
}

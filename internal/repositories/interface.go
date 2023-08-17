package repositories

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"
)

type (
	IUserRepo interface {
		FindAll(p models.Pagination) (*models.Pagination, error)
		Store(user *models.User) (*models.User, error)
		Update(user models.User, userID uint) (*models.User, error)
		FindByID(id uint) (*responses.UserResponse, error)
		FindByEmail(email string) (*responses.UserResponse, error)
		DeleteUser(user models.User) error
	}

	ILevelRepo interface {
		Store(l *models.Level) (*models.Level, error)
		Update(l models.Level, lid uint) (*models.Level, error)
		FindAll(p models.Pagination) (*models.Pagination, error)
		FindByID(lid uint) (*responses.LevelResponse, error)
		DeleteLevel(l models.Level) error
	}

	IRoleRepo interface {
		Store(r *models.Role) (*models.Role, error)
		Update(r models.Role, rid uint) (*models.Role, error)
		FindAll(p models.Pagination) (*models.Pagination, error)
		FindByID(lid uint) (*responses.RoleResponse, error)
		DeleteRole(r models.Role) error
	}
)

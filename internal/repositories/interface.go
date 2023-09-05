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
)

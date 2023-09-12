package repositories

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"

	"github.com/google/uuid"
)

type (
	IUserRepo interface {
		FindAll(p models.Pagination) (*models.Pagination, error)
		Store(user *models.User) (*models.User, error)
		Update(user models.User, uid uuid.UUID) (*models.User, error)
		FindByID(uid uuid.UUID) (*responses.UserResponse, error)
		FindByEmail(email string) (*responses.UserResponse, error)
		DeleteUser(user models.User) error
	}
)

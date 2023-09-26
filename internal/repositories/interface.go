package repositories

import (
	"project-skbackend/internal/controllers/responses"
	"project-skbackend/internal/models"

	"github.com/google/uuid"
)

type (
	IUserRepo interface {
		FindAll(p models.Pagination) (*models.Pagination, error)
		Create(user *models.User) (*models.User, error)
		Update(user models.User, uid uuid.UUID) (*models.User, error)
		FindByID(uid uuid.UUID) (*responses.UserResponse, error)
		FindByEmail(email string) (*responses.UserResponse, error)
		Delete(user models.User) error
	}

	IMemberRepo interface {
		Create(m *models.Member) (*models.Member, error)
		Update(m models.Member, mid uuid.UUID) (*models.Member, error)
		FindAll(p models.Pagination) (*models.Pagination, error)
		FindByID(mid uuid.UUID) (*responses.MemberResponse, error)
		Delete(m models.Member) error
	}
)
